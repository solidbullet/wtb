package wechat

import "net/http"

type Client struct {
	config     Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{config: cfg, httpClient: &http.Client{}}
}
