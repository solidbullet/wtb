package wechat

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type SessionResult struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (c *Client) Code2Session(code string) (*SessionResult, error) {
	u := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		url.QueryEscape(c.config.AppID),
		url.QueryEscape(c.config.AppSecret),
		url.QueryEscape(code))

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("wx code2session: %w", err)
	}
	defer resp.Body.Close()

	var result SessionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wx error %d: %s", result.ErrCode, result.ErrMsg)
	}
	return &result, nil
}
