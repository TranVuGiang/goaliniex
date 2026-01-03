package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/TranVuGiang/goaliniex"
)

func TestClient_GetKycInformation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := &goaliniex.KycInformationRequest{
		UserEmail: "tvugiang@gmail.com",
	}

	resp, err := client.GetKycInformation(ctx, req)
	if err != nil {
		t.Fatalf("GetKycInformation returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("response is nil")
	}

	if !resp.Success {
		t.Fatalf("request failed: %+v", resp)
	}

	if resp.Data.FirstName == "" && resp.Data.LastName == "" {
		t.Errorf("expected at least first or last name to be set")
	}

	if resp.Data.KycStatus == "" {
		t.Errorf("kycStatus should not be empty")
	}

	t.Logf("KYC Status: %s", resp.Data.KycStatus)
}
