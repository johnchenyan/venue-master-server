package controller

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
	"math"
	"sort"
	"strings"
	"time"
)

const PriceRealTimeType = "realTime"
const PriceT1TimeType = "t1"

const PriceGreaterThan = "greaterThan7.5"
const PricelessThanEqual = "lessThanEqual7.5"

type Settlement struct {
	Base
}

func (this *Settlement) FindSettlementData(ctx *gin.Context) {
	var data SettlementQueryParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	// 转换 StartTime 和 EndTime 为 time.Time 类型
	startTime, err := time.Parse("2006-01-02", data.StartTime)
	if err != nil {
		this.Error(ctx, "EndTime 转换错误")
		return
	}

	endTime, err := time.Parse("2006-01-02", data.EndTime)
	if err != nil {
		this.Error(ctx, "EndTime 转换错误")
		return
	}

	// 根据价格类型查询数据
	switch data.Type {
	case PriceRealTimeType:
		results, err := fetchSettlementData(data, startTime, endTime, PriceGreaterThan)
		if err != nil {
			this.Error(ctx, "查询错误")
			return
		}
		this.SuccessWithData(ctx, "获取成功", results)

	case PriceT1TimeType:
		results, err := fetchSettlementDataT(data, startTime, endTime, PriceGreaterThan)
		if err != nil {
			this.Error(ctx, "查询错误")
			return
		}
		this.SuccessWithData(ctx, "获取成功", results)

	default:
		this.Error(ctx, "查找类型不正确")
	}
}

func (this *Settlement) FindSettlementDataWithPagination(ctx *gin.Context) {
	var data SettlementQueryWithPaginationParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	resultsInterface, total, err := querySettlementDataPage(data, true)
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	items, err := processSettlementDataPage(resultsInterface)
	if err != nil {
		this.Error(ctx, "处理数据失败："+err.Error())
	}

	this.SuccessWithData(ctx, "获取成功", SettlementQueryWithPaginationResult{
		Data:     items,
		Total:    total,
		Page:     data.Page,
		PageSize: data.PageSize,
	})

}

func (this *Settlement) SettlementPointList(ctx *gin.Context) {
	typ := ctx.Param("type")
	if typ == "" {
		this.Error(ctx, "类型不能为空")
		return
	}
	if typ == PriceRealTimeType { // 实时价格
		var sps []model.SettlementPoint
		err := model.NewSettlementPoint().Find(&sps).Error
		if err != nil {
			this.Error(ctx, "数据库查找失败")
		}

		pointsMap := make(map[string][]string)

		for _, sp := range sps {
			pointsMap[sp.SettlementPointName] = append(pointsMap[sp.SettlementPointName], sp.SettlementPointType)
		}
		this.SuccessWithData(ctx, "获取成功", pointsMap)
	} else if typ == PriceT1TimeType {
		var sps []model.SettlementPointT
		err := model.NewSettlementPointT().Find(&sps).Error
		if err != nil {
			this.Error(ctx, "数据库查找失败")
		}

		var results []string
		for _, sp := range sps {
			results = append(results, sp.SettlementPointName)
		}

		this.SuccessWithData(ctx, "获取成功", results)
	} else {
		this.Error(ctx, "查找类型不正确")
	}
}

func (this *Settlement) FindSettlementAverage(ctx *gin.Context) {
	var data SettlementQueryParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	// 转换 StartTime 和 EndTime 为 time.Time 类型
	startTime, err := time.Parse("2006-01-02", data.StartTime)
	if err != nil {
		this.Error(ctx, "EndTime 转换错误")
		return
	}

	endTime, err := time.Parse("2006-01-02", data.EndTime)
	if err != nil {
		this.Error(ctx, "EndTime 转换错误")
		return
	}

	timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

	// 根据价格类型查询数据
	switch data.Type {
	case PriceRealTimeType:
		var results []SettlementQueryAverageResult
		for name, typs := range data.NameMap {
			for _, typ := range typs {
				resultsInterface, err := querySettlementData(data.Type, PricelessThanEqual, name, typ, startTime, endTime)
				if err != nil {
					continue
				}

				averagePrice := calculateAverageElectricityPrice(resultsInterface)
				results = append(results, SettlementQueryAverageResult{
					name,
					typ,
					timeRange,
					averagePrice,
				})
			}
		}
		this.SuccessWithData(ctx, "获取成功", results)

	case PriceT1TimeType:
		var results []SettlementQueryAverageResultT
		for name, _ := range data.NameMap {
			resultsInterface, err := querySettlementData(data.Type, PricelessThanEqual, name, "", startTime, endTime)
			if err != nil {
				continue
			}
			averagePrice := calculateAverageElectricityPrice(resultsInterface)
			results = append(results, SettlementQueryAverageResultT{
				name,
				timeRange,
				averagePrice,
			})
		}

		this.SuccessWithData(ctx, "获取成功", results)

	default:
		this.Error(ctx, "查找类型不正确")
	}
}

func (this *Settlement) DownLoadSettlementData(ctx *gin.Context) {
	var data SettlementQueryWithPaginationParam
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	resultsInterface, total, err := querySettlementDataPage(data, false)
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	items, err := processSettlementDataPage(resultsInterface)
	if err != nil {
		this.Error(ctx, "处理数据失败："+err.Error())
	}

	this.SuccessWithData(ctx, "获取成功", SettlementQueryWithPaginationResult{
		Data:     items,
		Total:    total,
		Page:     data.Page,
		PageSize: data.PageSize,
	})

}

func CreateSettlementData(data model.SettlementData) error {
	return model.NewSettlementData().Create(&data).Error
}

func CreateSettlementDataT(data model.SettlementDataT) error {
	return model.NewSettlementDataT().Create(&data).Error
}

func CreateSettlementPoint(data model.SettlementPoint) error {
	return model.NewSettlementPoint().Create(&data).Error
}

func CreateSettlementPointT(data model.SettlementPointT) error {
	return model.NewSettlementPointT().Create(&data).Error
}

func querySettlementData(timeType, priceType, name, pointType string, startDate time.Time, endDate time.Time) (interface{}, error) {
	var args []interface{}
	var whereClause string

	if timeType == PriceRealTimeType {
		whereClause += "settlement_point_name = ? AND settlement_point_type = ?"
		args = append(args, name, pointType)
	} else if timeType == PriceT1TimeType {
		whereClause += "(settlement_point_name = ?)"
		args = append(args, name)
	} else {
		return nil, xerrors.Errorf("无效的价格类型")
	}

	whereClause += " AND delivery_date BETWEEN ? AND ? AND "
	args = append(args, startDate, endDate)

	if priceType == PriceGreaterThan {
		whereClause += "settlement_point_price > ?"
	} else if priceType == PricelessThanEqual {
		whereClause += "settlement_point_price <= ?"
	} else {
		return nil, xerrors.Errorf("价格类型错误")
	}
	args = append(args, 7.5)

	return query(timeType, whereClause, args)
}

// fetchSettlementData 查询SettlementData类型数据
func fetchSettlementData(data SettlementQueryParam, startTime, endTime time.Time, priceCondition string) ([]SettlementQueryResult, error) {
	var results []SettlementQueryResult

	for name, typs := range data.NameMap {
		for _, typ := range typs {
			resultsInterface, err := querySettlementData(data.Type, priceCondition, name, typ, startTime, endTime)
			if err != nil {
				return nil, err
			}

			if settlementData, ok := resultsInterface.([]model.SettlementData); ok {
				results = append(results, processSettlementData(settlementData)...)
			}
		}
	}

	return results, nil
}

// fetchSettlementDataT 查询SettlementDataT类型数据
func fetchSettlementDataT(data SettlementQueryParam, startTime, endTime time.Time, priceCondition string) ([]SettlementQueryResultT, error) {
	var results []SettlementQueryResultT

	for name := range data.NameMap {
		resultsInterface, err := querySettlementData(data.Type, priceCondition, name, "", startTime, endTime)
		if err != nil {
			return nil, err
		}

		if settlementDataT, ok := resultsInterface.([]model.SettlementDataT); ok {
			results = append(results, processSettlementDataT(settlementDataT)...)
		}
	}

	return results, nil
}

// 以下是查询数据
func query(typ string, whereClause string, args []interface{}) (interface{}, error) {
	// 根据 typ 选择不同的数据模型
	if typ == PriceRealTimeType {
		var results []model.SettlementData
		err := model.NewSettlementData().Model(&model.SettlementData{}).Where(whereClause, args...).Find(&results).Error
		if err != nil {
			return nil, err
		}
		// 将结果转换为 interface{}
		return results, nil
	} else if typ == PriceT1TimeType {
		var results []model.SettlementDataT // 假设 SettlementDataT 是另一种模型
		err := model.NewSettlementDataT().Model(&model.SettlementDataT{}).Where(whereClause, args...).Find(&results).Error
		if err != nil {
			return nil, err
		}

		return results, nil

	} else {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}
}

// 以下是查询数据，分页查询
func queryPage(typ string, whereClause string, args []interface{}, page, pageSize int) (interface{}, int64, error) {
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 根据 typ 选择不同的数据模型
	if typ == PriceRealTimeType {
		// 查询总条数
		err := model.NewSettlementData().Where(whereClause, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}

		var results []model.SettlementData
		err = model.NewSettlementData().Model(&model.SettlementData{}).Where(whereClause, args...).Limit(pageSize).Offset(offset).Find(&results).Error
		if err != nil {
			fmt.Println("查询错误", err.Error())
			return nil, total, err
		}
		return results, total, nil
	} else if typ == PriceT1TimeType {
		// 查询总条数
		err := model.NewSettlementDataT().Where(whereClause, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		var results []model.SettlementDataT // 假设 SettlementDataT 是另一种模型
		err = model.NewSettlementDataT().Model(&model.SettlementDataT{}).Where(whereClause, args...).Limit(pageSize).Offset(offset).Find(&results).Error
		if err != nil {
			return nil, total, err
		}
		return results, total, nil
	} else {
		return nil, 0, fmt.Errorf("unsupported type: %s", typ)
	}
}

// 分页查询
func querySettlementDataPage(data SettlementQueryWithPaginationParam, isPage bool) (interface{}, int64, error) {
	var (
		err error

		queryAllTime bool
		queryAllName bool
		startTime    time.Time
		endTime      time.Time
	)

	if len(data.NameMap) == 0 { // 查询所有的接入点
		queryAllName = true
	}

	if data.StartTime == "" || data.EndTime == "" { // 查询所有时间
		queryAllTime = true
	}

	if data.StartTime != "" {
		startTime, err = time.Parse("2006-01-02", data.StartTime)
		if err != nil {
			return nil, 0, xerrors.Errorf("StartTime 转换错误")
		}
	}

	if data.EndTime != "" {
		endTime, err = time.Parse("2006-01-02", data.EndTime)
		if err != nil {
			return nil, 0, xerrors.Errorf("EndTime 转换错误")
		}
	}

	// 构建 whereClause 和 args
	var whereClause string
	var conditions []string
	var args []interface{}

	if !queryAllName {
		for name, types := range data.NameMap {
			for _, typ := range types {
				if data.Type == PriceRealTimeType {
					conditions = append(conditions, "(settlement_point_name = ? AND settlement_point_type = ?)")
					args = append(args, name, typ)
				} else if data.Type == PriceT1TimeType {
					conditions = append(conditions, "(settlement_point_name = ?)")
					args = append(args, name)
				}
			}
		}
		whereClause = strings.Join(conditions, " OR ")
	}

	if data.Price == PriceGreaterThan {
		if whereClause != "" { // 如果已有条件，添加 AND
			whereClause += " AND "
		}
		whereClause += "settlement_point_price > ?"
		args = append(args, 7.5)
	} else if data.Price == PricelessThanEqual {
		if whereClause != "" { // 如果已有条件，添加 AND
			whereClause += " AND "
		}
		whereClause += "settlement_point_price <= ?"
		args = append(args, 7.5)
	}

	if !queryAllTime {
		if whereClause != "" { // 如果已有条件，添加 AND
			whereClause += " AND "
		}
		whereClause += "delivery_date BETWEEN ? AND ?"
		args = append(args, startTime, endTime)
	}

	if isPage {
		return queryPage(data.Type, whereClause, args, data.Page, data.PageSize)
	} else {
		result, err := query(data.Type, whereClause, args)
		if err != nil {
			return nil, 0, err
		}
		return result, 0, nil
	}
}

type tmp struct {
	SettlementPointName string
	SettlementPointType string
	time                time.Time
}

func processSettlementData(data []model.SettlementData) []SettlementQueryResult {
	var (
		tmps    []tmp
		results []SettlementQueryResult
	)

	for _, data := range data {
		tmps = append(tmps, tmp{
			SettlementPointName: data.SettlementPointName,
			SettlementPointType: data.SettlementPointType,
			time:                calculateDeliveryTime(data),
		})
	}

	sort.SliceStable(tmps, func(i, j int) bool {
		return tmps[i].time.Before(tmps[j].time)
	})

	for i := 0; i < len(tmps); {
		current := tmps[i]
		startTime := current.time
		endTime := current.time
		name := current.SettlementPointName
		typ := current.SettlementPointType

		for j := i + 1; j < len(tmps); j++ {
			next := tmps[j]
			if next.time.Sub(endTime).Minutes() <= 15 {
				endTime = next.time
			} else {
				break
			}
		}

		// 记录合并后的结果
		timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02 15:04"), endTime.Add(15*time.Minute).Format("2006-01-02 15:04"))
		timeLength := int(endTime.Sub(startTime).Minutes()) + 15

		results = append(results, SettlementQueryResult{
			Name:       name,
			Type:       typ,
			TimeRange:  timeRange,
			TimeLength: fmt.Sprintf("%d", timeLength),
		})

		// 更新 i 的位置，跳到下一个不同的时间段
		i += int(endTime.Sub(startTime).Minutes()/15) + 1

	}

	return results
}

type tmpT struct {
	SettlementPointName string
	time                time.Time
}

func processSettlementDataT(data []model.SettlementDataT) []SettlementQueryResultT {
	var (
		tmps    []tmpT
		results []SettlementQueryResultT
	)

	for _, data := range data {
		tmps = append(tmps, tmpT{
			SettlementPointName: data.SettlementPointName,
			time:                data.DeliveryDate.Add(time.Hour * time.Duration(data.DeliveryHour)),
		})
	}

	sort.SliceStable(tmps, func(i, j int) bool {
		return tmps[i].time.Before(tmps[j].time)
	})

	for i := 0; i < len(tmps); {
		current := tmps[i]
		startTime := current.time
		endTime := current.time
		name := current.SettlementPointName

		for j := i + 1; j < len(tmps); j++ {
			next := tmps[j]
			if next.time.Sub(endTime).Hours() <= 1 {
				endTime = next.time
			} else {
				break
			}
		}

		// 记录合并后的结果
		timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02 15:04"), endTime.Add(1*time.Hour).Format("2006-01-02 15:04"))
		timeLength := int(endTime.Sub(startTime).Hours()) + 1

		results = append(results, SettlementQueryResultT{
			Name:       name,
			TimeRange:  timeRange,
			TimeLength: fmt.Sprintf("%d", timeLength),
		})

		// 更新 i 的位置，跳到下一个不同的时间段
		i += int(endTime.Sub(startTime).Hours()) + 1

	}

	return results
}

func calculateAverageElectricityPrice(data interface{}) float64 {
	total := 0.00
	count := 0

	// 处理 []model.SettlementData 类型
	if settlementData, ok := data.([]model.SettlementData); ok {
		count = len(settlementData)
		for _, settlement := range settlementData {
			total += settlement.SettlementPointPrice
		}
	} else if settlementDataT, ok := data.([]model.SettlementDataT); ok { // 处理 []model.SettlementDataT 类型
		count = len(settlementDataT)
		for _, settlement := range settlementDataT {
			total += settlement.SettlementPointPrice
		}
	}

	if count == 0 {
		return 0 // 防止除以零
	}

	average := total / float64(count)
	average = math.Round(average*100) / 100 // 保留两位小数
	return average
}

// 计算实际的时间
func calculateDeliveryTime(data model.SettlementData) time.Time {
	// 使用 DeliveryDate、DeliveryHour 和 DeliveryInterval 计算出实际的时间
	hour := int(data.DeliveryHour) - 1
	minute := int(data.DeliveryInterval-1) * 15 // 将 DeliveryInterval 转换为分钟
	return data.DeliveryDate.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
}

func processSettlementDataPage(data interface{}) ([]SettlementItem, error) {
	var result []SettlementItem
	if settlementData, ok := data.([]model.SettlementData); ok {
		for _, settlement := range settlementData {
			startTime := calculateDeliveryTime(settlement)
			timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02 15:04"), startTime.Add(15*time.Minute).Format("2006-01-02 15:04"))
			result = append(result, SettlementItem{
				Name:  settlement.SettlementPointName,
				Type:  settlement.SettlementPointType,
				Time:  timeRange,
				Price: settlement.SettlementPointPrice,
			})
		}
	} else if settlementDataT, ok := data.([]model.SettlementDataT); ok {
		for _, settlementT := range settlementDataT {
			startTime := settlementT.DeliveryDate.Add(time.Hour * time.Duration(settlementT.DeliveryHour))
			timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02 15:04"), startTime.Add(1*time.Hour).Format("2006-01-02 15:04"))
			result = append(result, SettlementItem{
				Name:  settlementT.SettlementPointName,
				Type:  "-",
				Time:  timeRange,
				Price: settlementT.SettlementPointPrice,
			})
		}
	}

	return result, nil
}
