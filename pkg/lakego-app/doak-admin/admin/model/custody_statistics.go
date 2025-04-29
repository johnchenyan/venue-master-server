package model

import (
	"time"

	"gorm.io/gorm"
)

type CustodyStatistics struct {
	ID                   uint      `gorm:"column:id;type:int;not null;primaryKey;autoIncrement" json:"id"`
	CustodyID            uint      `gorm:"column:custody_id;type:int;not null" json:"custody_id"`
	EnergyRatio          string    `gorm:"column:energy_ratio;type:varchar(20)" json:"energy_ratio,omitempty"`           // 可为空
	BasicHostingFee      string    `gorm:"column:basic_hosting_fee;type:varchar(20)" json:"basic_hosting_fee,omitempty"` // 可为空
	HourlyComputingPower string    `gorm:"column:hourly_computing_power;type:varchar(20)" json:"hourly_computing_power,omitempty"`
	TotalHostingFee      string    `gorm:"column:total_hosting_fee;type:varchar(20)" json:"total_hosting_fee,omitempty"`
	TotalIncomeBTC       string    `gorm:"column:total_income_btc;type:varchar(20)" json:"total_income_btc,omitempty"`
	TotalIncomeUSD       string    `gorm:"column:total_income_usd;type:varchar(20)" json:"total_income_usd,omitempty"`
	NetIncome            string    `gorm:"column:net_income;type:varchar(20)" json:"net_income,omitempty"`
	HostingFeeRatio      string    `gorm:"column:hosting_fee_ratio;type:varchar(20)" json:"hosting_fee_ratio,omitempty"`
	ReportDate           string    `gorm:"column:report_date;type:varchar(20);autoCreateTime" json:"report_date"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"`

	// 关联
	CustodyInfo CustodyInfo `gorm:"foreignKey:CustodyID" json:"custody_info"`
}

// 表名设置，带前缀
func (p *CustodyStatistics) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "custody_statistics"
}

// 创建模型实例
func NewCustodyStatisticsModel() *gorm.DB {
	return NewDB().Model(&CustodyStatistics{})
}
