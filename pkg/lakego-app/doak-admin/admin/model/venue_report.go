package model

import (
	"gorm.io/gorm"
	"time"
)

type VenueReport struct {
	ID                 uint      `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	SiteID             uint      `gorm:"column:site_id;type:int;not null;" json:"site_id"`
	SubAccount         string    `gorm:"column:sub_account;type:varchar(20);not null;" json:"sub_account"`
	RecordDate         time.Time `gorm:"column:record_date;type:date;not null;" json:"record_date"`
	AntpoolHashRate    string    `gorm:"column:antpool_hashrate;type:varchar(20);not null;" json:"antpool_hashrate"`
	F2poolHashRate     string    `gorm:"column:f2pool_hashrate;type:varchar(20);not null;" json:"f2pool_hashrate"`
	AntpoolDailyIncome string    `gorm:"column:antpool_daily_income;type:varchar(20);not null;" json:"antpool_daily_income"`
	F2poolDailyIncome  string    `gorm:"column:f2pool_daily_income;type:varchar(20);not null;" json:"f2pool_daily_income"`
	FBIncome           string    `gorm:"column:fb_income;type:varchar(20);not null;" json:"fb_income"`
	CreatedAt          time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

func (report *VenueReport) BeforeCreate(tx *gorm.DB) error {
	report.CreatedAt = time.Now()
	return nil
}

func (report *VenueReport) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "venue_report"
}

func NewVenueReportModel() *gorm.DB {
	return NewDB().Model(&VenueReport{})
}
