package model

import (
	"gorm.io/gorm"
	"time"
)

type VenueRecords struct {
	ID         uint      `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	TemplateID uint      `gorm:"column:template_id;not null;" json:"template_id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`

	Attributes []VenueRecordAttribute `gorm:"foreignKey:RecordID;references:ID;" json:"attributes"`
}

func (this *VenueRecords) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "venue_records"
}

func (this *VenueRecords) BeforeCreate(tx *gorm.DB) error {
	this.CreatedAt = time.Now()
	this.UpdatedAt = time.Now()
	return nil
}

func NewVenueRecord() *gorm.DB {
	return NewDB().Model(&VenueRecords{})
}
