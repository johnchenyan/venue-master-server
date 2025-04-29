package model

import (
	"time"

	"gorm.io/gorm"
)

type CustodyInfo struct {
	ID              uint      `gorm:"column:id;type:int;not null;primaryKey;autoIncrement" json:"id"`
	VenueName       string    `gorm:"column:venue_name;type:varchar(50);not null" json:"venue_name"`
	SubAccountName  string    `gorm:"column:sub_account_name;type:varchar(50);not null" json:"sub_account_name"`
	ObserverLink    string    `gorm:"column:observer_link;type:varchar(100)" json:"observer_link,omitempty"`        // 可为空
	EnergyRatio     string    `gorm:"column:energy_ratio;type:varchar(20)" json:"energy_ratio,omitempty"`           // 可为空
	BasicHostingFee string    `gorm:"column:basic_hosting_fee;type:varchar(20)" json:"basic_hosting_fee,omitempty"` // 可为空
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"`
}

func (p *CustodyInfo) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	return nil
}

// 表名设置，带前缀
func (p *CustodyInfo) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "custody_info"
}

// 创建模型实例
func NewCustodyInfoModel() *gorm.DB {
	return NewDB().Model(&CustodyInfo{})
}
