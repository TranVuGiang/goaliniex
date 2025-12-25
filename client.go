package goaliniex

import (
	"time"

	"github.com/TranVuGiang/goaliniex/config"
	"github.com/TranVuGiang/goaliniex/user"
	"resty.dev/v3"
)

type Client struct {
	cfg         *config.Config
	restyClient *resty.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:         cfg,
		restyClient: newRestyClient(cfg),
	}
}

func (c *Client) NewSubmitKYCHandle() *user.SubmitKYCHandle {
	return user.NewSubmitKYCHandle(c.cfg, c.restyClient)
}

func (c *Client) NewUserAlixHandle() *user.GetUserAlixHandle {
	return user.NewGetUserAlixHandle(c.cfg, c.restyClient)
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
