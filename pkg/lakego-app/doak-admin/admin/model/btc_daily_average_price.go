package model

import (
	"gorm.io/gorm"
	"time"
)

// DailyAveragePrice 表示每天的平均价格
type DailyAveragePrice struct {
	Date        string    `gorm:"column:date;type:varchar(20);not null;index:idx_date" json:"date"` // 添加索引
	CstAvgPrice string    `gorm:"column:cst_avg_price;type:decimal(15,2);not null" json:"cst_avg_price"`
	UtcAvgPrice string    `gorm:"column:utc_avg_price;type:decimal(15,2);not null" json:"utc_avg_price"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 自动填充
}

// TableName 设置表名，带前缀
func (d *DailyAveragePrice) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "daily_average_price"
}

// 新建模型实例
func NewDailyAveragePrice() *gorm.DB {
	return NewDB().Model(&DailyAveragePrice{})
}
