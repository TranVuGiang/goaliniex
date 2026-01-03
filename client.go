package goaliniex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/TranVuGiang/goaliniex/signer"
)

const UserAgent = "aliniex-go-sdk"

var (
	// Request lifecycle errors.
	ErrNilRequest    = errors.New("request is nil")
	ErrRequestBuild  = errors.New("failed to build request")
	ErrRequestSign   = errors.New("failed to sign request")
	ErrRequestEncode = errors.New("failed to encode request body")
	ErrInvalidParams = errors.New("invalid request params")

	// HTTP / transport errors.
	ErrHTTPFailure      = errors.New("http request failed")
	ErrUnexpectedStatus = errors.New("unexpected http status code")
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Option func(*Client)

type Client struct {
	BaseURL     string
	PartnerCode string
	SecretKey   string
	PrivateKey  []byte
	Logger      Logger
	Debug       bool
	HTTPClient  HTTPClient
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.BaseURL = url
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.Logger = logger
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.Debug = debug
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

func NewClient(
	baseURL string,
	partnerCode string,
	secretKey string,
	privateKey []byte,
	opts ...Option,
) *Client {
	client := &Client{
		BaseURL:     baseURL,
		PartnerCode: partnerCode,
		SecretKey:   secretKey,
		PrivateKey:  privateKey,
		HTTPClient:  http.DefaultClient,
		Logger:      slog.Default(),
		Debug:       false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) debug(msg string, attrs ...any) {
	if c.Debug {
		c.Logger.Debug(msg, attrs...)
	}
}

func paramsToMap(params any) (map[string]any, error) {
	if params == nil {
		return map[string]any{}, nil
	}

	if m, ok := params.(map[string]any); ok {
		return m, nil
	}

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) buildRequest(req *request) error {
	if req == nil {
		return ErrNilRequest
	}

	fullURL := c.BaseURL + req.Endpoint

	headers := http.Header{}
	if req.Header != nil {
		headers = req.Header.Clone()
	}

	headers.Set("Content-Type", "application/json")
	headers.Set("User-Agent", UserAgent)

	signature, err := signer.Sign(c.PrivateKey, req.SigningData)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestSign, err)
	}

	bodyMap, err := paramsToMap(req.Params)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidParams, err)
	}

	bodyMap["partnerCode"] = c.PartnerCode
	bodyMap["signature"] = signature

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestEncode, err)
	}

	c.debug("http request", "url", fullURL)
	c.debug("http request body", "body", string(bodyBytes))

	req.FullURL = fullURL
	req.Header = headers
	req.Body = bytes.NewReader(bodyBytes)

	return nil
}

func (c *Client) execute(ctx context.Context, req *request) ([]byte, error) {
	if err := c.buildRequest(req); err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		req.FullURL,
		req.Body,
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header = req.Header

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTPFailure, err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.debug("http response", "status", resp.StatusCode)
	c.debug("http response body", "body", string(responseBody))

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"%w: status=%d body=%s",
			ErrUnexpectedStatus,
			resp.StatusCode,
			string(responseBody),
		)
	}

	return responseBody, nil
}
