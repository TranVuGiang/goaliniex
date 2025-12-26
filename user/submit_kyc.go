package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/TranVuGiang/goaliniex/config"
	"github.com/TranVuGiang/goaliniex/signature"
	"resty.dev/v3"
)

var ErrSubmitKYCToAlix = errors.New("failed to submit kyc to alix")

type SubmitKYCRequest struct {
	UserEmail        string  `json:"userEmail"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	DateOfBirth      string  `json:"dateOfBirth"`
	Gender           string  `json:"gender"` // "male" | "female"
	Nationality      string  `json:"nationality"`
	Type             string  `json:"type"` // "ID_CARD" | "PASSPORT"
	NationalID       string  `json:"nationalId"`
	IssueDate        string  `json:"issueDate"`
	ExpiryDate       string  `json:"expiryDate"`
	AddressLine1     string  `json:"addressLine1"`
	AddressLine2     string  `json:"addressLine2"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	ZipCode          string  `json:"zipCode"`
	FrontIDImage     string  `json:"frontIdImage"`
	BackIDImage      string  `json:"backIdImage"`
	HoldIDImage      string  `json:"holdIdImage"`
	PhoneNumber      *string `json:"phoneNumber,omitempty"`
	PhoneCountryCode *string `json:"phoneCountryCode,omitempty"`
}

type AlixResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      *AlixData `json:"data"`
	ErrorCode int32     `json:"errorCode"`
}

type AlixData struct {
	NationalID string `json:"nationalId"`
	KycStatus  string `json:"kycStatus"`
	Signature  string `json:"signature"`
}

type KYCRequestBody struct {
	SubmitKYCRequest
	Signature string `json:"signature"`
}

func (h *SubmitKYCHandle) SubmitKYC(ctx context.Context, req *SubmitKYCRequest) (*AlixResponse, error) {
	targetURL := h.cfg.BaseURL + "/user/submit-kyc"

	payload := fmt.Sprintf("%s|%s|%s|%s", h.cfg.PartnerCode, req.UserEmail, req.Nationality, h.cfg.SecretKey)

	sig, err := signature.SignPayload(payload, h.cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}

	var alixResponse AlixResponse

	requestBody := KYCRequestBody{
		SubmitKYCRequest: *req,
		Signature:        sig,
	}

	resp, err := h.restyClient.R().
		SetContext(ctx).
		SetBody(requestBody).
		SetResult(&alixResponse).
		SetError(&alixResponse).
		Post(targetURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSubmitKYCToAlix, err)
	}

	if !resp.IsSuccess() {

		return nil, fmt.Errorf("%w: %+v", ErrSubmitKYCToAlix, alixResponse)
	}

	return &alixResponse, nil
}

type SubmitKYCHandle struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func NewSubmitKYCHandle(cfg *config.Config, restyClient *resty.Client) *SubmitKYCHandle {
	return &SubmitKYCHandle{
		cfg:         cfg,
		restyClient: restyClient,
	}
}
