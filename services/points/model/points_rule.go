package model
type PointsRule struct {
	ID        uint   `gorm:"primaryKey;column:id" json:"id"`
	Name      string `gorm:"size:50;column:name" json:"name"`
	Type      string `gorm:"size:30;column:type" json:"type"`
	ConfigJSON string `gorm:"type:text;column:config_json" json:"config_json"`
	Status    int16  `gorm:"default:1;column:status" json:"status"`
}
func (PointsRule) TableName() string { return "points_rules" }
