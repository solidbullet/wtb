package wechat

import "fmt"

// SubscribeMsgData 订阅消息数据项
type SubscribeMsgData struct {
	Value string `json:"value"`
}

// SendSubscribeMsg 发送订阅消息（Mock）
func (c *Client) SendSubscribeMsg(openid, templateID string, data map[string]SubscribeMsgData, page string) error {
	if c.config.AppID == "" {
		return fmt.Errorf("AppID 未配置")
	}
	// TODO: 接入微信小程序订阅消息 API
	return nil
}
