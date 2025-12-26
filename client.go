package goaliniex

import (
	"context"
	"time"

	"github.com/TranVuGiang/goaliniex/config"
	"github.com/TranVuGiang/goaliniex/user"
	"resty.dev/v3"
)

type Client struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func NewAlixClient(cfg *config.Config) *Client {
	return &Client{
		cfg:         cfg,
		restyClient: newRestyClient(cfg),
	}
}

func (c *Client) SubmitKYC(ctx context.Context, req *user.SubmitKYCRequest) (*user.AlixResponse, error) {
	handle := user.NewSubmitKYCHandle(c.cfg, c.restyClient)
	return handle.SubmitKYC(ctx, req)
}

func (c *Client) GetUserInfo(ctx context.Context, userEmail string) (*user.UserAlixResponse, error) {
	handle := user.NewGetUserAlixHandle(c.cfg, c.restyClient)
	return handle.GetUserInfo(ctx, userEmail)
}

func newRestyClient(cfg *config.Config) *resty.Client {
	const defaultRetryCount int = 3

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(30 * time.Second).
		SetRetryCount(defaultRetryCount).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	return client
}
