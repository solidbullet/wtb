package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

type tokenCache struct {
	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

var cache = &tokenCache{}

func (c *Client) GetAccessToken() (string, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.accessToken != "" && time.Now().Before(cache.expiresAt.Add(-5*time.Minute)) {
		return cache.accessToken, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", c.config.AppID, c.config.AppSecret)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat error %d: %s", result.ErrCode, result.ErrMsg)
	}

	cache.accessToken = result.AccessToken
	cache.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, nil
}

func (c *Client) GetWXACodeUnlimited(accessToken, scene, page string, checkPath bool, envVersion string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimited?access_token=%s", accessToken)
	body := map[string]interface{}{
		"scene":       scene,
		"check_path":  checkPath,
	}
	if page != "" {
		body["page"] = page
	}
	if envVersion != "" {
		body["env_version"] = envVersion
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 && data[0] == '{' {
		var result struct {
			ErrCode int    `json:"errcode"`
			ErrMsg  string `json:"errmsg"`
		}
		if err := json.Unmarshal(data, &result); err == nil && result.ErrCode != 0 {
			return nil, fmt.Errorf("wxacode error %d: %s", result.ErrCode, result.ErrMsg)
		}
	}
	return data, nil
}
