package model

import (
	"gorm.io/gorm"
	"time"
)

type BtcUsdCandle struct {
	Timestamp  time.Time `gorm:"column:timestamp;type:timestamp;not null;index:idx_timestamp" json:"timestamp"`
	PriceLow   string    `gorm:"column:price_low;type:decimal(15,4);default:null" json:"price_low"`
	PriceHigh  string    `gorm:"column:price_high;type:decimal(15,4);default:null" json:"price_high"`
	PriceOpen  string    `gorm:"column:price_open;type:decimal(15,4);default:null" json:"price_open"`
	PriceClose string    `gorm:"column:price_close;type:decimal(15,4);default:null" json:"price_close"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 设置表名，带前缀
func (b *BtcUsdCandle) TableName() string {
	prefix := GetConfig("prefix").(string) // 获取配置中的前缀
	return prefix + "btc_usd_candle"
}

// 新建模型实例（可选）
func NewBtcUsdCandle() *gorm.DB {
	return NewDB().Model(&BtcUsdCandle{})
}
