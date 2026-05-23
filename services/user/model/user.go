package model

import "time"

type User struct {
	ID               uint      `gorm:"primaryKey;column:id" json:"id"`
	OpenID           string    `gorm:"uniqueIndex;size:64;column:openid" json:"openid"`
	UnionID          string    `gorm:"size:64;column:unionid" json:"unionid"`
	Nickname         string    `gorm:"size:100;column:nickname" json:"nickname"`
	AvatarURL        string    `gorm:"size:500;column:avatar_url" json:"avatar_url"`
	Phone            string    `gorm:"size:20;column:phone" json:"phone"`
	MemberLevel      int16     `gorm:"default:0;column:member_level" json:"member_level"`
	Balance          int       `gorm:"default:0;column:balance" json:"balance"`
	TotalConsumption int       `gorm:"default:0;column:total_consumption" json:"total_consumption"`
	TotalOrders      int       `gorm:"default:0;column:total_orders" json:"total_orders"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string { return "users" }
