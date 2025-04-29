package model

import "gorm.io/gorm"

type TemplateFields struct {
	ID         uint   `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	TemplateID uint   `gorm:"column:template_id;not null;" json:"template_id"`
	FieldName  string `gorm:"column:field_name;type:varchar(100);not null;" json:"field_name"`
	FieldType  string `gorm:"column:field_type;type:varchar(50);not null;" json:"field_type"`
	FieldOrder int    `gorm:"column:field_order;not null;" json:"field_order"`
}

func (this *TemplateFields) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "template_fields"
}

func NewTemplateFields() *gorm.DB {
	return NewDB().Model(&TemplateFields{})
}
