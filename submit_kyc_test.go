package goaliniex_test

import (
	"context"
	"testing"
	"time"

	"github.com/TranVuGiang/goaliniex"
)

func TestClient_SubmitKyc(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)

	req := &goaliniex.SubmitKycRequest{
		UserEmail:        "chinontest0000001@gmail.com",
		FirstName:        "Chí",
		LastName:         "Nguyen Thanh",
		DateOfBirth:      "1994-06-06",
		Gender:           goaliniex.GenderMale,
		Nationality:      "VN",
		DocumentType:     goaliniex.IDTypeIDCard,
		NationalID:       "12130163",
		IssueDate:        "1994-06-06",
		ExpiryDate:       "1994-06-06",
		AddressLine1:     "Ho Chi Minh",
		AddressLine2:     "Go Vap",
		City:             "Ho Chi Minh",
		State:            "Ho Chi Minh",
		ZipCode:          "70000",
		FrontIDImage:     "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/.........",
		BackIDImage:      "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/.........",
		HoldIDImage:      "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/.........",
		PhoneNumber:      "0968861116",
		PhoneCountryCode: "84",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.SubmitKyc(ctx, req)
	if err != nil {
		t.Fatalf("SubmitKyc failed: %v", err)
	}

	if !resp.Success {
		t.Fatalf("API returned unsuccessful response: %+v", resp)
	}

	if resp.Data == nil {
		t.Fatal("response data is nil")
	}

	if resp.Data.NationalID == "" {
		t.Error("expected nationalId in response")
	}

	if resp.Data.KycStatus == "" {
		t.Error("expected kycStatus in response")
	}
}
