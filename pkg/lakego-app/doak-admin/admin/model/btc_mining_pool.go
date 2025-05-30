package model

import (
	"gorm.io/gorm"
	"time"
)

// MiningPool 表示矿池信息
type MiningPool struct {
	ID                  uint      `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`                       // 主键，自增ID
	PoolName            string    `gorm:"column:pool_name;type:varchar(255);not null" json:"pool_name"`                // 矿池名称
	PoolType            string    `gorm:"column:pool_type;type:enum('自营', '矿机贷款');not null" json:"pool_type"`          // 矿池类型
	PoolCategory        string    `gorm:"column:pool_category;type:enum('主矿池', '备用矿池');not null" json:"pool_category"` // 矿池类别
	TheoreticalHashrate float64   `gorm:"column:theoretical_hashrate;type:decimal(10,2)" json:"theoretical_hashrate"`  // 理论算力
	Link                string    `gorm:"column:link;type:varchar(255)" json:"link"`                                   // 链接
	SortOrder           int       `gorm:"column:sort_order;type:int;default:0" json:"sort_order"`                      // 排序字段
	IsEnabled           bool      `gorm:"column:is_enabled;type:tinyint(1);default:1" json:"is_enabled"`               // 启用状态
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"`           // 更新时间
}

// TableName 设置表名，带前缀
func (m *MiningPool) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "mining_pools"
}

// 新建模型实例
func NewMiningPool() *gorm.DB {
	return NewDB().Model(&MiningPool{})
}

// MiningSettlementRecord 表示矿池结算记录
type MiningSettlementRecord struct {
	ID                  uint      `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`                                 // 主键，自增ID
	PoolID              uint      `gorm:"column:pool_id;type:int;not null" json:"pool_id"`                                       // 外键，关联到矿池表的 ID
	SettlementDate      string    `gorm:"column:settlement_date;type:varchar(30);not null" json:"settlement_date"`               // 结算日期
	SettlementHashrate  float64   `gorm:"column:settlement_hashrate;type:decimal(15,2);not null" json:"settlement_hashrate"`     // 结算算力
	SettlementProfitBtc float64   `gorm:"column:settlement_profit_btc;type:decimal(15,8);not null" json:"settlement_profit_btc"` // 结算收益 BTC
	SettlementProfitFb  float64   `gorm:"column:settlement_profit_fb;type:decimal(15,2);not null" json:"settlement_profit_fb"`   // 结算算力 FB
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"`                     // 创建时间
}

// TableName 设置表名，带前缀
func (m *MiningSettlementRecord) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "mining_settlement_records"
}

// 新建模型实例
func NewMiningSettlementRecord() *gorm.DB {
	return NewDB().Model(&MiningSettlementRecord{})
}

// MiningPoolStatus 表示矿池状态信息
type MiningPoolStatus struct {
	ID              uint      `gorm:"column:id;type:int;primaryKey;autoIncrement" json:"id"`                         // 主键，自增ID
	PoolID          uint      `gorm:"column:pool_id;type:int;not null" json:"pool_id"`                               // 外键，关联到矿池表的 ID
	CurrentHashrate float64   `gorm:"column:current_hashrate;type:decimal(15,2);not null" json:"current_hashrate"`   // 实时算力
	OnlineMachines  int       `gorm:"column:online_machines;type:int;not null" json:"online_machines"`               // 实时在线机器数
	OfflineMachines int       `gorm:"column:offline_machines;type:int;not null" json:"offline_machines"`             // 离线机器数
	Last24hHashrate float64   `gorm:"column:last_24h_hashrate;type:decimal(15,2);not null" json:"last_24h_hashrate"` // 24小时算力
	LastUpdate      time.Time `gorm:"column:last_update;type:timestamp;autoUpdateTime" json:"last_update"`           // 最后更新时间
}

// TableName 设置表名，带前缀
func (m *MiningPoolStatus) TableName() string {
	prefix := GetConfig("prefix").(string)
	return prefix + "mining_pool_status"
}

// 新建模型实例
func NewMiningPoolStatus() *gorm.DB {
	return NewDB().Model(&MiningPoolStatus{})
}
