package controller

import (
	"errors"
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"math"
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

	// 如果是备用矿池，找到对应的主矿池则添加，没有则报错
	if data.PoolCategory == "备用矿池" {
		exist, err := IsMasterPoolExist(data.PoolType, data.PoolName)
		if err != nil {
			this.Error(ctx, "查找对应的主矿池失败")
			return
		}
		if !exist {
			this.Error(ctx, "对应的主矿池不存在，请添加！")
			return
		}
	}

	bp, err := createBtcMiningPool(model.MiningPool{
		PoolName:            data.PoolName,
		PoolType:            data.PoolType,
		Country:             data.Country,
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

	// 如果是备用矿池，找到对应的主矿池则添加，没有则报错
	if data.PoolCategory == "备用矿池" {
		exist, err := IsMasterPoolExist(data.PoolType, data.PoolName)
		if err != nil {
			this.Error(ctx, "查找对应的主矿池失败")
			return
		}
		if !exist {
			this.Error(ctx, "对应的主矿池不存在，请添加！")
			return
		}
	}

	err = updateBtcMiningPool(model.MiningPool{
		ID:                  data.ID,
		PoolName:            data.PoolName,
		PoolType:            data.PoolType,
		Country:             data.Country,
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

// 判断备用矿池的主矿池是否存在
func IsMasterPoolExist(pool_type, pool_name string) (bool, error) {
	var pool model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ? AND pool_category = ? AND pool_name = ?", pool_type, "主矿池", pool_name).First(&pool).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到，说明不存在
			return false, nil
		}
		// 其他错误
		return false, err
	}
	if pool_type == pool.PoolType && pool_name == pool.PoolName {
		return true, nil
	}
	return false, nil
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
			LastSettlementHashRate:  fmt.Sprintf("%.2f TH/s", settlementRecord.SettlementHashrate),
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

// overview

func (this *BtcMiningPool) TotalRealTimeStatus(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TotalRealTimeStatus 程序发生了未捕获的错误: %v\n", err)
			// 记录错误或采取其他行动
		}
	}()

	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ?", typ).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
	}

	var poolIDs []uint
	for _, miningPool := range miningPools {
		poolIDs = append(poolIDs, miningPool.ID)
	}

	// 计算半小时前的时间
	halfHourAgo := time.Now().Add(-30 * time.Minute)

	// 查询最新状态
	var latestStatuses []model.MiningPoolStatus
	// 使用子查询获取每个矿池最新状态的时间
	subQuery := model.NewMiningPoolStatus().
		Select("pool_id, MAX(last_update) AS latest_update").
		Where("last_update >= ? AND pool_id IN (?)", halfHourAgo, poolIDs).
		Group("pool_id")

	// 根据子查询获取最新状态的详细信息
	status := model.MiningPoolStatus{}
	tableName := status.TableName()
	err = model.NewMiningPoolStatus().
		Table(tableName).
		Joins("JOIN (?) AS latest ON "+tableName+".pool_id = latest.pool_id AND "+tableName+".last_update = latest.latest_update", subQuery).
		Find(&latestStatuses).Error

	if err != nil {
		this.Error(ctx, "状态获取失败")
		return
	}

	var (
		totalCurrentHashrate       float64
		totalMasterCurrentHashrate float64
		totalBackUpCurrentHashrate float64
		totalOnline                int
		totalOffline               int
	)

	for _, status := range latestStatuses {
		totalCurrentHashrate += status.CurrentHashrate
		totalOnline += status.OnlineMachines
		totalOffline += status.OfflineMachines

		for _, miningPool := range miningPools {
			if status.PoolID == miningPool.ID {
				if miningPool.PoolCategory == "主矿池" {
					totalMasterCurrentHashrate += status.CurrentHashrate
				} else {
					totalBackUpCurrentHashrate += status.CurrentHashrate
				}
			}
		}
	}

	// 将结果打包成一个 map
	result := map[string]any{
		"totalCurrentHashRate":       fmt.Sprintf("%.2f", totalCurrentHashrate/1e15),
		"totalMasterCurrentHashrate": fmt.Sprintf("%.2f", totalMasterCurrentHashrate/1e15),
		"totalBackUpCurrentHashrate": fmt.Sprintf("%.2f", totalBackUpCurrentHashrate/1e15),
		"totalOnline":                totalOnline,
		"totalOffline":               totalOffline,
	}

	this.SuccessWithData(ctx, "获取成功", result)
}

func (this *BtcMiningPool) TotalLastDayStatus(ctx *gin.Context) {
	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ?", typ).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
	}

	var poolIDs []uint
	for _, miningPool := range miningPools {
		poolIDs = append(poolIDs, miningPool.ID)
	}

	// 查询 SettlementDate 最大的日期
	var maxSettlementDate string
	err = model.NewMiningSettlementRecord().
		Select("MAX(settlement_date)").
		Where("pool_id IN ?", poolIDs).
		Scan(&maxSettlementDate).Error

	if err != nil || maxSettlementDate == "" {
		this.Error(ctx, "获取最大结算日期失败")
		return
	}

	// 查询所有 SettlementDate 等于最大日期的记录
	var latestStatuses []model.MiningSettlementRecord
	err = model.NewMiningSettlementRecord().
		Where("settlement_date = ? AND pool_id IN ?", maxSettlementDate, poolIDs).
		Find(&latestStatuses).Error

	if err != nil {
		this.Error(ctx, "状态获取失败")
		return
	}

	var (
		totalHash            float64
		totalTheoreticalHash float64
		totalProfitBtc       float64
		totalMasterProfitBtc float64
		totalBackUpProfitBtc float64
		totalHashEfficiency  string
	)

	for _, status := range latestStatuses {
		totalHash += status.SettlementHashrate // TH/s
		totalProfitBtc += status.SettlementProfitBtc

		for _, miningPool := range miningPools {
			if miningPool.ID == status.PoolID && miningPool.PoolCategory == "主矿池" { // PH/s
				totalTheoreticalHash += status.SettlementTheoreticalHashrate
				totalMasterProfitBtc += status.SettlementProfitBtc
			} else if miningPool.ID == status.PoolID && miningPool.PoolCategory == "备用矿池" {
				totalBackUpProfitBtc += status.SettlementProfitBtc
			}
		}
	}

	if totalTheoreticalHash != 0 {
		totalHashEfficiency = fmt.Sprintf("%.2f", 100*totalHash/(totalTheoreticalHash*1e3))
	}

	totalProfitBtc = math.Round(totalProfitBtc*1e8) / 1e8 // 保留8位小数

	// 将结果打包成一个 map
	result := map[string]any{
		"lastSettlementDate":   maxSettlementDate,
		"totalProfitBtc":       totalProfitBtc,
		"totalMasterProfitBtc": totalMasterProfitBtc,
		"totalBackUpProfitBtc": totalBackUpProfitBtc,
		"totalHashEfficiency":  totalHashEfficiency,
	}

	this.SuccessWithData(ctx, "获取成功", result)
}

type DailyEfficiency struct {
	Date       string  `json:"date"`
	Efficiency float64 `json:"efficiency"`
}

type EfficiencyResponse struct {
	Efficiencies      []DailyEfficiency `json:"efficiencys"`
	AverageEfficiency float64           `json:"averageEfficiency"`
}

func (this *BtcMiningPool) TotalLastWeekStatus(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("TotalLastWeekStatus 程序发生了未捕获的错误: %v\n", err)
			// 记录错误或采取其他行动
		}
	}()

	fmt.Printf("TotalLastWeekStatus: ctx address: %p\n", ctx)
	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ?", typ).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
		return
	}

	var poolIDs []uint
	for _, miningPool := range miningPools {
		poolIDs = append(poolIDs, miningPool.ID)
	}

	// 查询 SettlementDate 最大的日期
	var maxSettlementDate string
	err = model.NewMiningSettlementRecord().
		Select("MAX(settlement_date)").
		Where("pool_id IN ?", poolIDs).
		Scan(&maxSettlementDate).Error

	if err != nil || maxSettlementDate == "" {
		this.Error(ctx, "获取最大结算日期失败")
		return
	}

	// 计算从 maxSettlementDate 往前的七天的有效率
	startDate, err := time.Parse("2006-01-02", maxSettlementDate)
	if err != nil {
		this.Error(ctx, "日期格式不正确")
		return
	}

	var efficiencies []DailyEfficiency
	var totalEfficiency float64
	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, -i).Format("2006-01-02")
		efficiency, err := getOneDayHashEfficiency(miningPools, poolIDs, date)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}

		totalEfficiency += efficiency
		efficiencies = append(efficiencies, DailyEfficiency{
			Date:       date,
			Efficiency: efficiency,
		})
	}

	averageEfficiency := totalEfficiency / float64(len(efficiencies))
	averageEfficiency = math.Round(averageEfficiency*100) / 100 // 保留两位小数

	response := EfficiencyResponse{
		Efficiencies:      efficiencies,
		AverageEfficiency: averageEfficiency,
	}

	this.SuccessWithData(ctx, "获取成功", response)
}

func (this *BtcMiningPool) LastestHashRateEfficiency(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("LastestHashRateEfficiency 程序发生了未捕获的错误: %v\n", err)
			// 记录错误或采取其他行动
		}
	}()

	fmt.Printf("LastestHashRateEfficiency: ctx address: %p\n", ctx)

	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	dayStr := ctx.Param("day")
	if dayStr == "" {
		this.Error(ctx, "时间不能为空")
		return
	}
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ?", typ).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
		return
	}

	var poolIDs []uint
	for _, miningPool := range miningPools {
		poolIDs = append(poolIDs, miningPool.ID)
	}

	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		this.Error(ctx, "时间转换失败")
		return
	}

	// 查询 SettlementDate 最大的日期
	var maxSettlementDate string
	err = model.NewMiningSettlementRecord().
		Select("MAX(settlement_date)").
		Where("pool_id IN ?", poolIDs).
		Scan(&maxSettlementDate).Error

	if err != nil || maxSettlementDate == "" {
		this.Error(ctx, "获取最大结算日期失败")
		return
	}

	// 计算从 maxSettlementDate 往前的七天的有效率
	startDate, err := time.Parse("2006-01-02", maxSettlementDate)
	if err != nil {
		this.Error(ctx, "日期格式不正确")
		return
	}

	var efficiencies []struct {
		Date       string  `json:"date"`
		Efficiency float64 `json:"efficiency"`
	}
	for i := 0; i < int(day); i++ {
		date := startDate.AddDate(0, 0, -i).Format("2006-01-02")
		efficiency, err := getOneDayHashEfficiency(miningPools, poolIDs, date)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}

		efficiencies = append(efficiencies, struct {
			Date       string  `json:"date"`
			Efficiency float64 `json:"efficiency"`
		}{
			Date:       date,
			Efficiency: efficiency,
		})
	}

	this.SuccessWithData(ctx, "获取成功", efficiencies)
}

func getOneDayHashEfficiency(miningPools []model.MiningPool, poolIDs []uint, day string) (float64, error) {
	// 查询所有 SettlementDate 等于最大日期的记录
	var latestStatuses []model.MiningSettlementRecord
	err := model.NewMiningSettlementRecord().
		Where("settlement_date = ? AND pool_id IN ?", day, poolIDs).
		Find(&latestStatuses).Error

	if err != nil {
		return 0, xerrors.Errorf("状态获取失败")
	}

	var (
		totalHash            float64
		totalTheoreticalHash float64
		totalHashEfficiency  float64
	)

	for _, status := range latestStatuses {
		totalHash += status.SettlementHashrate // TH/s

		for _, miningPool := range miningPools {
			if miningPool.ID == status.PoolID && miningPool.PoolCategory == "主矿池" { // PH/s
				totalTheoreticalHash += status.SettlementTheoreticalHashrate
			}
		}
	}

	if totalTheoreticalHash != 0 {
		totalHashEfficiency = 100 * totalHash / (totalTheoreticalHash * 1e3)
		totalHashEfficiency = math.Round(totalHashEfficiency*100) / 100 // 保留两位小数
	}
	return totalHashEfficiency, nil
}

func (this *BtcMiningPool) LastestHashRate(ctx *gin.Context) {
	typ := ctx.Param("poolType")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	dayStr := ctx.Param("day")
	if dayStr == "" {
		this.Error(ctx, "时间不能为空")
		return
	}
	var miningPools []model.MiningPool
	err := model.NewMiningPool().Where("pool_type = ?", typ).Find(&miningPools).Error
	if err != nil {
		this.Error(ctx, "数据库获取失败")
		return
	}

	var poolIDs []uint
	for _, miningPool := range miningPools {
		poolIDs = append(poolIDs, miningPool.ID)
	}

	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		this.Error(ctx, "时间转换失败")
		return
	}

	// 查询 SettlementDate 最大的日期
	var maxSettlementDate string
	err = model.NewMiningSettlementRecord().
		Select("MAX(settlement_date)").
		Where("pool_id IN ?", poolIDs).
		Scan(&maxSettlementDate).Error

	if err != nil || maxSettlementDate == "" {
		this.Error(ctx, "获取最大结算日期失败")
		return
	}

	// 计算从 maxSettlementDate 往前的七天的有效率
	startDate, err := time.Parse("2006-01-02", maxSettlementDate)
	if err != nil {
		this.Error(ctx, "日期格式不正确")
		return
	}

	var DayHashRates []struct {
		Date        string  `json:"date"`
		DayHashRate float64 `json:"day_hash_rate"`
	}
	for i := 0; i < int(day); i++ {
		date := startDate.AddDate(0, 0, -i).Format("2006-01-02")
		hash, err := getOneDayHash(miningPools, poolIDs, date)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}

		DayHashRates = append(DayHashRates, struct {
			Date        string  `json:"date"`
			DayHashRate float64 `json:"day_hash_rate"`
		}{
			Date:        date,
			DayHashRate: hash,
		})
	}

	this.SuccessWithData(ctx, "获取成功", DayHashRates)
}

func getOneDayHash(miningPools []model.MiningPool, poolIDs []uint, day string) (float64, error) {
	// 查询所有 SettlementDate 等于最大日期的记录
	var latestStatuses []model.MiningSettlementRecord
	err := model.NewMiningSettlementRecord().
		Where("settlement_date = ? AND pool_id IN ?", day, poolIDs).
		Find(&latestStatuses).Error

	if err != nil {
		return 0, xerrors.Errorf("状态获取失败")
	}

	var totalHash float64

	for _, status := range latestStatuses {
		totalHash += status.SettlementHashrate // TH/s
	}

	totalHash = totalHash / 1e3
	totalHash = math.Round(totalHash*100) / 100

	return totalHash, nil
}
