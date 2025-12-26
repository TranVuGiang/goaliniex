# GoAliniex

Go SDK for Alix API integration - A simple and efficient library for KYC (Know Your Customer) operations.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Features

- 🔐 **RSA-SHA256 Signature** - Secure request signing with PKCS1/PKCS8 support
- 🚀 **Simple API** - Clean and intuitive interface
- ⚡ **HTTP Client Built-in** - Automatic retry and timeout configuration
- 🎯 **Context Support** - Full control over request lifecycle
- 📦 **Zero Dependencies** - Only uses `resty.dev/v3` for HTTP client

## Installation

```bash
go get github.com/TranVuGiang/goaliniex
```

## Quick Start

### 1. Initialize Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/TranVuGiang/goaliniex"
    "github.com/TranVuGiang/goaliniex/config"
    "github.com/TranVuGiang/goaliniex/user"
)

func main() {
    // Read private key from file
    privateKey, err := os.ReadFile("private_key.pem")
    if err != nil {
        log.Fatal(err)
    }

    // Create configuration
    cfg := &config.Config{
        BaseURL:     "https://api.alix.com",
        PartnerCode: "YOUR_PARTNER_CODE",
        SecretKey:   "YOUR_SECRET_KEY",
        PrivateKey:  privateKey,
    }

    // Initialize client
    client := goaliniex.NewAlixClient(cfg)
}
```

### 2. Get User Information

```go
ctx := context.Background()

// Get user KYC information by email
resp, err := client.GetUserInfo(ctx, "user@example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("KYC Status: %s\n", resp.Data.KycStatus)
fmt.Printf("Full Name: %s %s\n", resp.Data.FirstName, resp.Data.LastName)
fmt.Printf("National ID: %s\n", resp.Data.NationalID)
```

### 3. Submit KYC

```go
ctx := context.Background()

// Prepare KYC request
kycRequest := &user.SubmitKYCRequest{
    UserEmail:    "newuser@example.com",
    FirstName:    "John",
    LastName:     "Doe",
    DateOfBirth:  "1990-01-01",
    Gender:       "male",
    Nationality:  "VN",
    Type:         "ID_CARD",
    NationalID:   "123456789",
    IssueDate:    "2015-01-01",
    ExpiryDate:   "2025-01-01",
    AddressLine1: "123 Main Street",
    AddressLine2: "Apt 4B",
    City:         "Hanoi",
    State:        "Hanoi",
    ZipCode:      "100000",
    FrontIDImage: "base64_encoded_image",
    BackIDImage:  "base64_encoded_image",
    HoldIDImage:  "base64_encoded_image",
}

// Submit KYC
resp, err := client.SubmitKYC(ctx, kycRequest)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Success: %v\n", resp.Success)
fmt.Printf("Message: %s\n", resp.Message)
```

## Advanced Usage

### With Timeout

```go
import "time"

// Create context with 10 second timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.GetUserInfo(ctx, "user@example.com")
if err != nil {
    // Handle timeout or other errors
    log.Fatal(err)
}
```

### With Cancellation

```go
ctx, cancel := context.WithCancel(context.Background())

// Cancel after 5 seconds in goroutine
go func() {
    time.Sleep(5 * time.Second)
    cancel()
}()

resp, err := client.GetUserInfo(ctx, "user@example.com")
if err != nil {
    if ctx.Err() == context.Canceled {
        fmt.Println("Request was cancelled")
    }
}
```

### Error Handling

```go
import (
    "errors"
    "github.com/TranVuGiang/goaliniex/user"
)

resp, err := client.GetUserInfo(ctx, "user@example.com")
if err != nil {
    if errors.Is(err, user.ErrGetUserFromAlix) {
        fmt.Println("Failed to get user from Alix API")
    }
    log.Fatal(err)
}

// Check response status
if !resp.Success {
    fmt.Printf("API Error: %s (Code: %d)\n", resp.Message, resp.ErrorCode)
}
```

## Configuration

### Config Structure

```go
type Config struct {
    BaseURL     string // Alix API base URL
    PartnerCode string // Your partner code
    SecretKey   string // Your secret key
    PrivateKey  []byte // RSA private key in PEM format
}
```

### Private Key Format

The library supports both PKCS1 and PKCS8 private key formats:

**PKCS1:**
```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----
```

**PKCS8:**
```
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEF...
-----END PRIVATE KEY-----
```

### HTTP Client Configuration

The library automatically configures the HTTP client with:
- **Timeout:** 30 seconds
- **Retry Count:** 3 attempts
- **Retry Wait Time:** 1 second
- **Max Retry Wait Time:** 5 seconds

## API Reference

### Client Methods

#### `NewAlixClient(cfg *config.Config) *Client`
Creates a new Alix client instance.

#### `GetUserInfo(ctx context.Context, userEmail string) (*user.UserAlixResponse, error)`
Retrieves KYC information for a user by email.

**Parameters:**
- `ctx` - Context for request control
- `userEmail` - User's email address

**Returns:**
- `UserAlixResponse` - User KYC information
- `error` - Error if request fails

#### `SubmitKYC(ctx context.Context, req *user.SubmitKYCRequest) (*user.AlixResponse, error)`
Submits KYC information for a user.

**Parameters:**
- `ctx` - Context for request control
- `req` - KYC request data

**Returns:**
- `AlixResponse` - Submission response
- `error` - Error if request fails

### Response Types

#### `UserAlixResponse`

```go
type UserAlixResponse struct {
    Success   bool
    Message   string
    Data      *UserAlixData
    ErrorCode int32
}

type UserAlixData struct {
    FirstName        string
    LastName         string
    DateOfBirth      string
    Gender           string
    Nationality      string
    IDType           string
    NationalID       string
    IssueDate        string
    ExpiryDate       string
    Address          string
    FrontIDImage     string
    BackIDImage      string
    HoldIDImage      string
    PhoneNumber      string
    PhoneCountryCode string
    KycStatus        string
    RejectReason     string
}
```

#### `AlixResponse`

```go
type AlixResponse struct {
    Success   bool
    Message   string
    Data      *AlixData
    ErrorCode int32
}

type AlixData struct {
    NationalID string
    KycStatus  string
    Signature  string
}
```

## Local Development

If you want to use this library locally before pushing to GitHub:

### In your `go.mod`:

```go
module myproject

go 1.21

require github.com/TranVuGiang/goaliniex v0.0.0

replace github.com/TranVuGiang/goaliniex => /Users/mbp/Template/go_lib/goaliniex
```

### Then run:

```bash
go mod tidy
```

This allows you to develop and test the library locally. When ready to publish:

1. Remove the `replace` directive
2. Push to GitHub
3. Tag a version:
```bash
git tag v1.0.0
git push origin v1.0.0
```
4. Update your project:
```bash
go get github.com/TranVuGiang/goaliniex@v1.0.0
```

## Best Practices

### 1. Reuse Client Instance

```go
// ✅ Good - Create client once
client := goaliniex.NewAlixClient(cfg)

for _, email := range emails {
    resp, err := client.GetUserInfo(ctx, email)
    // Process response
}
```

```go
// ❌ Bad - Creating client in loop
for _, email := range emails {
    client := goaliniex.NewAlixClient(cfg)  // Don't do this!
    resp, err := client.GetUserInfo(ctx, email)
}
```

### 2. Always Use Context

```go
// ✅ Good - With timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.GetUserInfo(ctx, email)
```

```go
// ⚠️  Acceptable - Background context
ctx := context.Background()
resp, err := client.GetUserInfo(ctx, email)
```

### 3. Handle Errors Properly

```go
resp, err := client.GetUserInfo(ctx, email)
if err != nil {
    // Log error with context
    log.Printf("Failed to get user info for %s: %v", email, err)
    return err
}

// Check API response status
if !resp.Success {
    log.Printf("API returned error: %s (code: %d)", resp.Message, resp.ErrorCode)
    return fmt.Errorf("API error: %s", resp.Message)
}
```

### 4. Secure Private Key Storage

```go
// ✅ Good - Read from secure file
privateKey, err := os.ReadFile("/secure/path/private_key.pem")
if err != nil {
    log.Fatal(err)
}

// ✅ Good - From environment variable
privateKeyPEM := os.Getenv("ALIX_PRIVATE_KEY")
privateKey := []byte(privateKeyPEM)

// ❌ Bad - Hardcoded in source code
privateKey := []byte("-----BEGIN RSA PRIVATE KEY-----\n...")  // Don't do this!
```

## Example Integration with Echo Framework

```go
package handler

import (
    "net/http"

    "github.com/TranVuGiang/goaliniex"
    "github.com/TranVuGiang/goaliniex/user"
    "github.com/labstack/echo/v4"
)

type GetUserHandler struct {
    alixClient *goaliniex.Client
}

func NewGetUserHandler(client *goaliniex.Client) *GetUserHandler {
    return &GetUserHandler{alixClient: client}
}

func (h *GetUserHandler) Handle(c echo.Context) error {
    email := c.QueryParam("email")

    resp, err := h.alixClient.GetUserInfo(c.Request().Context(), email)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    if !resp.Success {
        return echo.NewHTTPError(http.StatusBadRequest, resp.Message)
    }

    return c.JSON(http.StatusOK, resp.Data)
}
```

## Troubleshooting

### Invalid Signature Error

**Problem:** API returns "Invalid signature" error

**Solutions:**
1. Verify your private key matches the public key registered with Alix
2. Check the payload format matches exactly: `partnerCode|email|secretKey`
3. Ensure private key is in correct PEM format
4. Verify partner code and secret key are correct

### Timeout Errors

**Problem:** Requests timing out

**Solutions:**
1. Increase context timeout:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
```

2. Check network connectivity to Alix API
3. Verify BaseURL is correct

### Import Errors

**Problem:** Cannot find package

**Solution:**
```bash
go get github.com/TranVuGiang/goaliniex@latest
go mod tidy
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- Open an issue on [GitHub](https://github.com/TranVuGiang/goaliniex/issues)
- Contact: tranvugiang@example.com

## Changelog

### v1.0.0 (Current)
- Initial release
- Support for GetUserInfo API
- Support for SubmitKYC API
- RSA-SHA256 signature implementation
- Context support for all requests
- Automatic retry and timeout configuration
