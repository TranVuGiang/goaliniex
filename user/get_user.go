package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/TranVuGiang/goaliniex/config"
	"github.com/TranVuGiang/goaliniex/signature"
	"resty.dev/v3"
)

var ErrGetUserFromAlix = errors.New("failed to get user from alix")

type UserAlixResponse struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Data      *UserAlixData `json:"data"`
	ErrorCode int32         `json:"errorCode"`
}

type UserAlixData struct {
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	DateOfBirth      string `json:"dateOfBirth"`
	Gender           string `json:"gender"`
	Nationality      string `json:"nationality"`
	IDType           string `json:"idType"`
	NationalID       string `json:"nationalId"`
	IssueDate        string `json:"issueDate"`
	ExpiryDate       string `json:"expiryDate"`
	Address          string `json:"address"`
	FrontIDImage     string `json:"frontIdImage"`
	BackIDImage      string `json:"backIdImage"`
	HoldIDImage      string `json:"holdIdImage"`
	PhoneNumber      string `json:"phoneNumber"`
	PhoneCountryCode string `json:"phoneCountryCode"`
	KycStatus        string `json:"kycStatus"`
	RejectReason     string `json:"rejectReason"`
}

func (h *GetUserAlixHandle) GetUserInfo(ctx context.Context, userEmail string) (*UserAlixResponse, error) {
	targetURL := h.cfg.BaseURL + "/user/get-kyc-information"

	payload := fmt.Sprintf("%s|%s|%s", h.cfg.PartnerCode, userEmail, h.cfg.SecretKey)

	fmt.Println("PAYLOAD: ", payload)

	sig, err := signature.SignPayload(payload, h.cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}

	fmt.Println("SIGNATURE: ", sig)

	var respBody UserAlixResponse

	requestBody := map[string]any{
		"partnerCode": h.cfg.PartnerCode,
		"userEmail":   userEmail,
		"signature":   sig,
	}

	resp, err := h.restyClient.R().
		SetContext(ctx).
		SetBody(requestBody).
		SetResult(&respBody).
		SetError(&respBody).
		Post(targetURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetUserFromAlix, err)
	}

	if !resp.IsSuccess() {

		return nil, fmt.Errorf("%w: %+v", ErrGetUserFromAlix, respBody)
	}

	fmt.Println("Response: ", respBody)

	return &respBody, nil
}

type GetUserAlixHandle struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func NewGetUserAlixHandle(cfg *config.Config, restyClient *resty.Client) *GetUserAlixHandle {
	return &GetUserAlixHandle{
		cfg:         cfg,
		restyClient: restyClient,
	}
}
