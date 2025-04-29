package model

import (
	"gorm.io/gorm"
)

type LinkInfo struct {
	ID          uint   `gorm:"column:id;type:int;not null;primaryKey;autoIncrement;" json:"id"`
	SiteName    string `gorm:"column:site_name;type:varchar(50);not null;" json:"site_name"`
	SubAccount  string `gorm:"column:sub_account;type:varchar(50);not null;" json:"sub_account"`
	AntpoolLink string `gorm:"column:antpool_link;type:varchar(100);not null;" json:"antpool_link"`
	F2poolLink  string `gorm:"column:f2pool_link;type:varchar(100);not null;" json:"f2pool_link"`
	SortOrder   uint   `gorm:"column:sort_order;type:int;not null;default:0;" json:"sort_order"` // 设置默认值为0
	//CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

//func (link *LinkInfo) BeforeCreate(tx *gorm.DB) error {
//	link.CreatedAt = time.Now()
//	return nil
//}

func (link *LinkInfo) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "link_info"
}

func NewLinkInfoModel() *gorm.DB {
	return NewDB().Model(&LinkInfo{})
}
