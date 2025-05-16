package controller

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"math"
	"sort"
	"time"
)

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

	var results []SettlementQueryResult

	for name, typs := range data.NameMap {
		for _, typ := range typs {
			settlementData, err := querySettlementLimitData(name, typ, startTime, endTime)
			if err != nil {
				fmt.Printf("查询错误：%v\n", err)
				continue
			}
			results = append(results, processSettlementData(settlementData)...)
		}
	}

	this.SuccessWithData(ctx, "获取成功", results)
}

func (this *Settlement) SettlementPointList(ctx *gin.Context) {
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

	var results []SettlementQueryAverageResult

	timeRange := fmt.Sprintf("%s 至 %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

	for name, typs := range data.NameMap {
		for _, typ := range typs {
			settlementData, err := querySettlementData(name, typ, startTime, endTime)
			if err != nil {
				fmt.Printf("查询错误：%v\n", err)
				continue
			}

			averagePrice := calculateAverageElectricityPrice(settlementData)
			results = append(results, SettlementQueryAverageResult{
				name,
				typ,
				timeRange,
				averagePrice,
			})
		}
	}

	this.SuccessWithData(ctx, "获取成功", results)
}

func CreateSettlementData(data model.SettlementData) error {
	return model.NewSettlementData().Create(&data).Error
}

func CreateSettlementPoint(data model.SettlementPoint) error {
	return model.NewSettlementPoint().Create(&data).Error
}

func querySettlementLimitData(name string, typ string, startDate time.Time, endDate time.Time) ([]model.SettlementData, error) {
	var results []model.SettlementData
	err := model.NewSettlementData().Where("settlement_point_name = ? AND settlement_point_type = ? AND delivery_date BETWEEN ? AND ? AND settlement_point_price > ?",
		name, typ, startDate, endDate, 7.5).Find(&results).Error
	return results, err
}

func querySettlementData(name string, typ string, startDate time.Time, endDate time.Time) ([]model.SettlementData, error) {
	var results []model.SettlementData
	err := model.NewSettlementData().Where("settlement_point_name = ? AND settlement_point_type = ? AND delivery_date BETWEEN ? AND ? AND settlement_point_price <= ?",
		name, typ, startDate, endDate, 7.5).Find(&results).Error
	return results, err
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

func calculateAverageElectricityPrice(data []model.SettlementData) float64 {
	if len(data) == 0 {
		return 0 // 防止除以零
	}

	total := 0.0
	for _, settlement := range data {
		total += settlement.SettlementPointPrice
	}

	average := total / float64(len(data))
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
