package goaliniex_test

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/TranVuGiang/goaliniex"
)

func newTestClient(t *testing.T) *goaliniex.Client {
	t.Helper()

	privateKey, err := os.ReadFile("./alix-private-key.pem")
	if err != nil {
		t.Skipf("skipping test: unable to read private key: %v", err)
	}

	partnerCode := strings.TrimSpace(os.Getenv("ALIX_PARTNER_CODE"))
	if partnerCode == "" {
		t.Skip("skipping test: ALIX_PARTNER_CODE not set")
	}

	secretKey := strings.TrimSpace(os.Getenv("ALIX_SECRET_KEY"))
	if secretKey == "" {
		t.Skip("skipping test: ALIX_SECRET_KEY not set")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	client, err := goaliniex.NewClient(
		"https://sandbox.alixpay.com",
		partnerCode,
		secretKey,
		privateKey,
		goaliniex.WithDebug(true),
		goaliniex.WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return client
}
