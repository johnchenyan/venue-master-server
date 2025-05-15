package model

import (
	"gorm.io/gorm"
)

type SettlementPoint struct {
	SettlementPointID   uint   `gorm:"column:settlement_point_id;primaryKey;autoIncrement" json:"settlement_point_id"`      // 结算点唯一标识
	SettlementPointName string `gorm:"column:settlement_point_name;type:varchar(20);not null" json:"settlement_point_name"` // 结算点名称
	SettlementPointType string `gorm:"column:settlement_point_type;type:varchar(20);not null" json:"settlement_point_type"` // 结算点类型
}

// TableName 设置表名，带前缀
func (s *SettlementPoint) TableName() string {
	prefix := GetConfig("prefix").(string) // 获取配置中的前缀
	return prefix + "settlement_points"    // 更改为实际表名
}

// 新建模型实例（可选）
func NewSettlementPoint() *gorm.DB {
	return NewDB().Model(&SettlementPoint{})
}
