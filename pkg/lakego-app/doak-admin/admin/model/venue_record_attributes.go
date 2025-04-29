package model

import "gorm.io/gorm"

type VenueRecordAttribute struct {
	ID         uint   `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	RecordID   uint   `gorm:"column:record_id;not null;" json:"record_id"`
	FieldID    uint   `gorm:"column:field_id;not null;" json:"field_id"`
	FieldName  string `gorm:"column:field_name;type:varchar(100);not null;" json:"field_name"`
	FieldValue string `gorm:"column:field_value;type:text;" json:"field_value"`
}

func (this *VenueRecordAttribute) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "venue_record_attributes"
}

func NewVenueRecordAttribute() *gorm.DB {
	return NewDB().Model(&VenueRecordAttribute{})
}
