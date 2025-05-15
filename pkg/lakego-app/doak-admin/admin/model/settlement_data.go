package model

import (
	"gorm.io/gorm"
	"time"
)

type SettlementData struct {
	SettlementPointName  string    `gorm:"column:settlement_point_name;type:varchar(20);not null" json:"settlement_point_name"`
	SettlementPointType  string    `gorm:"column:settlement_point_type;type:varchar(20);not null" json:"settlement_point_type"`
	DeliveryDate         time.Time `gorm:"column:delivery_date;type:date;not null" json:"delivery_date"`
	DeliveryHour         uint8     `gorm:"column:delivery_hour;type:tinyint;not null" json:"delivery_hour"`
	DeliveryInterval     uint8     `gorm:"column:delivery_interval;type:tinyint;not null" json:"delivery_interval"`
	SettlementPointPrice float64   `gorm:"column:settlement_point_price;type:float;not null" json:"settlement_point_price"`
}

// TableName 设置表名，带前缀
func (s *SettlementData) TableName() string {
	prefix := GetConfig("prefix").(string) // 获取配置中的前缀
	return prefix + "settlement_data"      // 更改为实际表名
}

// 新建模型实例（可选）
func NewSettlementData() *gorm.DB {
	return NewDB().Model(&SettlementData{})
}
