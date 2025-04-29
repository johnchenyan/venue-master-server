package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrNotFound = errors.New("not found")

type VenueTemplates struct {
	ID           uint      `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	TemplateName string    `gorm:"column:template_name;type:varchar(100);not null;" json:"template_name"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`

	Fields []TemplateFields `gorm:"foreignKey:TemplateID;references:ID;" json:"fields"`
}

func (this *VenueTemplates) BeforeCreate(tx *gorm.DB) error {
	this.CreatedAt = time.Now()
	this.UpdatedAt = time.Now()
	return nil
}

func (this *VenueTemplates) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "venue_templates"
}

func NewVenueTemplate() *gorm.DB {
	return NewDB().Model(&VenueTemplates{})
}
