package controller

import (
	"errors"
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var MiningPoolChannel chan MiningPoolRequest

type MiningPoolRequest struct {
	MiningPool model.MiningPool
	ResultChan chan<- error
}

type BtcMiningPool struct {
	Base
}

func init() {
	fmt.Println("init mining pool channel")
	MiningPoolChannel = make(chan MiningPoolRequest, 20)
}

// btc mining pool

func (this *BtcMiningPool) ListBtcMiningPool(ctx *gin.Context) {
	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}

	category := ctx.Param("poolCategory")
	if category == "" {
		this.Error(ctx, "类型不能为空")
		return
	}

	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ? AND pool_category = ?", typ, category).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
	}

	this.SuccessWithData(ctx, "获取成功", miningPools)
}

func (this *BtcMiningPool) CreateBtcMiningPool(ctx *gin.Context) {
	var data BtcMiningPoolParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	hr, err := strconv.ParseFloat(data.TheoreticalHashrate, 64)
	if err != nil {
		this.Error(ctx, "无效的理论算力")
		return
	}

	bp, err := createBtcMiningPool(model.MiningPool{
		PoolName:            data.PoolName,
		PoolType:            data.PoolType,
		PoolCategory:        data.PoolCategory,
		TheoreticalHashrate: hr,
		Link:                data.Link,
	})
	if err != nil {
		this.Error(ctx, "新增矿池失败")
		return
	}

	// 创建结果通道
	resultChan := make(chan error)
	// 创建 MiningPoolRequest 结构体
	request := MiningPoolRequest{
		MiningPool: bp,
		ResultChan: resultChan,
	}

	// 使用 goroutine 异步发送数据到 CustodyInfoChannel
	go func(req MiningPoolRequest) {
		MiningPoolChannel <- req
	}(request)

	// 等待处理结果
	select {
	case err := <-resultChan:
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		this.Success(ctx, "新增矿池成功！")
		return
	case <-time.After(30 * time.Second): // 设置超时
		this.Error(ctx, "请求超时")
		return
	}

	this.Success(ctx, "新增矿池成功！")
}

func (this *BtcMiningPool) UpdateBtcMiningPool(ctx *gin.Context) {
	var data BtcMiningPoolParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	hr, err := strconv.ParseFloat(data.TheoreticalHashrate, 64)
	if err != nil {
		this.Error(ctx, "无效的理论算力")
		return
	}

	err = updateBtcMiningPool(model.MiningPool{
		ID:                  data.ID,
		PoolName:            data.PoolName,
		PoolType:            data.PoolType,
		PoolCategory:        data.PoolCategory,
		TheoreticalHashrate: hr,
		Link:                data.Link,
	})
	if err != nil {
		this.Error(ctx, "新增矿池失败")
		return
	}

	this.Success(ctx, "新增矿池成功！")
}

func ListBtcMiningPool() ([]model.MiningPool, error) {
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Find(&miningPools).Error
	if err != nil {
		return nil, err
	}

	return miningPools, nil
}

func createBtcMiningPool(data model.MiningPool) (model.MiningPool, error) {
	err := model.NewMiningPool().Create(&data).Error
	return data, err
}

func updateBtcMiningPool(data model.MiningPool) error {
	updates := map[string]interface{}{
		"pool_name":            data.PoolName,
		"pool_type":            data.PoolType,
		"pool_category":        data.PoolCategory,
		"theoretical_hashrate": data.TheoreticalHashrate,
		"link":                 data.Link,
	}

	return model.NewMiningPool().Where("id = ?", data.ID).
		Updates(updates).Error
}

// btc mining pool settlement records

func CreateBtcMiningSettlementRecord(data model.MiningSettlementRecord) error {
	return model.NewMiningSettlementRecord().Create(&data).Error
}

func UpdateBtcMiningFBProfit(pool_id uint, date string, data float64) error {
	return model.NewMiningSettlementRecord().Where("pool_id = ? AND settlement_date = ?", pool_id, date).
		Update("settlement_profit_fb", data).Error
}

// 判断某个日期是否存在
func IsMiningRecordExist(pool_id uint, date string) (bool, error) {
	var record model.MiningSettlementRecord
	err := model.NewMiningSettlementRecord().Where("pool_id = ? AND settlement_date = ?", pool_id, date).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到，说明不存在
			return false, nil
		}
		// 其他错误
		return false, err
	}
	if record.SettlementDate != date || record.PoolID != pool_id {
		return false, nil
	}
	return true, nil
}

// 判断某个日期是否存在
func IsFBRecordUpdated(pool_id uint, date string) (bool, error) {
	var record model.MiningSettlementRecord
	err := model.NewMiningSettlementRecord().Where("pool_id = ? AND settlement_date = ?", pool_id, date).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到，说明不存在
			return false, nil
		}
		// 其他错误
		return false, err
	}
	if record.PoolID == pool_id && record.SettlementProfitFb == 0 {
		return false, nil
	}
	return true, nil
}

// btc mining pool settlement status

func CreateBtcMiningPoolStatus(data model.MiningPoolStatus) error {
	return model.NewMiningPoolStatus().Create(&data).Error
}

func (this *BtcMiningPool) ListBtcMiningPoolHashRate(ctx *gin.Context) {
	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}

	category := ctx.Param("poolCategory")
	if category == "" {
		this.Error(ctx, "类型不能为空")
		return
	}

	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ? AND pool_category = ?", typ, category).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
	}

	var result []BtcMiningPoolHashResult

	for _, miningPool := range miningPools {
		var settlementRecord model.MiningSettlementRecord
		var poolStatus model.MiningPoolStatus
		// 查询矿池的最新结算记录
		err = model.NewMiningSettlementRecord().
			Where("pool_id = ?", miningPool.ID).
			Order("settlement_date DESC").
			First(&settlementRecord).Error // 使用 First 获取最新记录
		if err != nil {
			continue
		}

		// 查询矿池的状态
		err = model.NewMiningPoolStatus().
			Where("pool_id = ?", miningPool.ID).
			Order("last_update DESC").
			First(&poolStatus).Error // 使用 First 获取最新状态
		if err != nil {
			continue
		}

		currentHash := formatHashRate(poolStatus.CurrentHashrate)
		lastHash := formatHashRate(poolStatus.Last24hHashrate)
		effective, err := conversionHashRateEffective(settlementRecord.SettlementHashrate, miningPool.TheoreticalHashrate)
		if err != nil {
			continue
		}

		result = append(result, BtcMiningPoolHashResult{
			PoolName:                miningPool.PoolName,
			CurrentHashRate:         currentHash,
			Online:                  poolStatus.OnlineMachines,
			Offline:                 poolStatus.OfflineMachines,
			LastHashRate:            lastHash,
			TheoreticalHashRate:     fmt.Sprintf("%.2f PH/s", miningPool.TheoreticalHashrate),
			LastHashRateEffective:   effective,
			LastSettlementProfitBtc: settlementRecord.SettlementProfitBtc,
			LastSettlementProfitFB:  settlementRecord.SettlementProfitFb,
			Link:                    miningPool.Link,
			LastSettlementDate:      settlementRecord.SettlementDate,
			UpdateTime:              poolStatus.LastUpdate.Format("2006-01-02 15:04:05"),
		})
	}

	this.SuccessWithData(ctx, "获取成功", result)
}

func formatHashRate(hash float64) (result string) {
	if hash >= 1e18 {
		result = fmt.Sprintf("%.2f EH/s", hash/1e18)
	} else if hash >= 1e15 { // PH/s
		result = fmt.Sprintf("%.2f PH/s", hash/1e15)
	} else if hash >= 1e12 { // TH/s
		result = fmt.Sprintf("%.2f TH/s", hash/1e12)
	} else if hash >= 1e9 { // TH/s
		result = fmt.Sprintf("%.2f GH/s", hash/1e9)
	} else {
		result = fmt.Sprintf("%.2f H/s", hash) // 保留原始单位
	}
	return result
}

// lastSettlementHash 单位为TH/s, theoreticalHash 单位为PH/s
func conversionHashRateEffective(lastSettlementHash, theoreticalHash float64) (string, error) {
	conversionRate := lastSettlementHash / (theoreticalHash * 1000) * 100

	// 返回转化率，格式化为字符串
	return fmt.Sprintf("%.2f%%", conversionRate), nil
}
