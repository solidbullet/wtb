package model

// RechargePlan 充值档位定义
type RechargePlan struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`       // 门槛金额（分）
	FinalAmount int    `json:"final_amount"` // 到账金额（分）
	GiftAmount  int    `json:"gift_amount"`  // 赠送金额（分）
	MemberLevel int16  `json:"member_level"` // 达到的会员等级
	Desc        string `json:"desc"`
}

// RechargePlans 系统内置充值档位
var RechargePlans = []RechargePlan{
	{ID: 1, Name: "开通会员", Amount: 19900, FinalAmount: 0, GiftAmount: 0, MemberLevel: 1, Desc: "享受会员价折扣，余额不变"},
	{ID: 2, Name: "预充值 ¥1000", Amount: 100000, FinalAmount: 120000, GiftAmount: 20000, MemberLevel: 2, Desc: "到账 ¥1200（送 ¥200）"},
	{ID: 3, Name: "预充值 ¥2000", Amount: 200000, FinalAmount: 240000, GiftAmount: 40000, MemberLevel: 2, Desc: "到账 ¥2400（送 ¥400）"},
	{ID: 4, Name: "预充值 ¥5000", Amount: 500000, FinalAmount: 600000, GiftAmount: 100000, MemberLevel: 2, Desc: "到账 ¥6000（送 ¥1000）"},
}

// UpgradeInfo 用户当前可升级信息
type UpgradeInfo struct {
	CurrentPlan    *RechargePlan `json:"current_plan"`    // 当前已达到的档位（可能为nil）
	NextPlan       *RechargePlan `json:"next_plan"`       // 下一个可升级档位（可能为nil）
	TotalRecharged int           `json:"total_recharged"` // 累计充值金额（分）
	Balance        int           `json:"balance"`         // 当前余额（分）
}
