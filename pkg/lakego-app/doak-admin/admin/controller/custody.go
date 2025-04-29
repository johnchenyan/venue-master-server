package controller

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/gin-gonic/gin"
	"time"
)

var CustodyInfoChannel chan model.CustodyInfo

type Custody struct {
	Base
}

func init() {
	fmt.Println("init custody")
	CustodyInfoChannel = make(chan model.CustodyInfo, 20)
}

func (this *Custody) ListCustodyInfo(ctx *gin.Context) {
	custodyInfos, err := ListCustodyInfo()
	if err != nil {
		this.Error(ctx, fmt.Sprintf("获取托管信息失败: %s", err.Error()))
		return
	}

	this.SuccessWithData(ctx, "获取成功", custodyInfos)
}

// CreateCustodyInfo NewCustodyInfo 新增托管信息
func (this *Custody) CreateCustodyInfo(ctx *gin.Context) {
	var data CustodyInfo
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	ci, err := createCustodyInfo(data)
	if err != nil {
		this.Error(ctx, "新增托管信息数据失败")
	}

	// 更新此账号的托管统计
	// 使用 goroutine 异步发送数据到 CustodyInfoChannel
	go func(custodyInfo CustodyInfo) {
		csi := model.CustodyInfo{
			ID:              ci.ID,
			VenueName:       custodyInfo.VenueName,
			SubAccountName:  custodyInfo.SubAccountName,
			ObserverLink:    custodyInfo.ObserverLink,
			EnergyRatio:     custodyInfo.EnergyRatio,
			BasicHostingFee: custodyInfo.BasicHostingFee,
		}
		CustodyInfoChannel <- csi
	}(data)

	this.Success(ctx, "新增记托管信息成功！")
}

// DeleteCustodyInfo 删除托管信息
func (this *Custody) DeleteCustodyInfo(ctx *gin.Context) {
	custodyInfoId := ctx.Param("custodyInfoId")
	if custodyInfoId == "" {
		this.Error(ctx, "托管信息ID不能为空")
		return
	}

	// 根据记录ID删除对应的删除记录数据
	if err := deleteCustodyInfoById(custodyInfoId); err != nil {
		this.Error(ctx, "删除托管信息数据失败")
	}

	this.Success(ctx, "删除记录成功！")
}

// UpdateCustodyInfo 更新托管信息
func (this *Custody) UpdateCustodyInfo(ctx *gin.Context) {
	var data CustodyInfoUpdate
	if err := this.ShouldBindJSON(ctx, &data); err != nil {
		this.Error(ctx, "请求数据不正确")
		return
	}

	if data.ID == 0 {
		this.Error(ctx, "托管信息ID不能为空")
		return
	}

	if err := updateCustodyInfo(data); err != nil {
		this.Error(ctx, "更新记录失败")
	}

	this.Success(ctx, "更新记录成功！")
}

func (this *Custody) ListCustodyStatistics(ctx *gin.Context) {
	timeRange := ctx.Param("timeRange")
	if timeRange == "" {
		this.Error(ctx, "时间范围不能为空")
		return
	}
	// 计算时间范围的起点时间
	var startTime time.Time
	now := time.Now()

	switch timeRange {
	case "all", "":
		// 不做时间限制
		startTime = time.Time{} // 零值，代表不限制
	case "1days":
		startTime = now.AddDate(0, 0, -1)
	case "3days":
		println("3days")
		startTime = now.AddDate(0, 0, -3)
	case "7days":
		startTime = now.AddDate(0, 0, -7)
	case "1month":
		startTime = now.AddDate(0, -1, 0)
	case "3month":
		startTime = now.AddDate(0, -3, 0)
	case "6month":
		startTime = now.AddDate(0, -6, 0)
	// 你可以添加更多的范围
	default:
		// 默认值，比如全部
		startTime = time.Time{}
	}

	custodyStatistics, err := ListCustodyInfoWithTimeRange(startTime)
	if err != nil {
		this.Error(ctx, fmt.Sprintf("获取托管统计失败: %s", err.Error()))
		return
	}

	this.SuccessWithData(ctx, "获取成功", custodyStatistics)
}

func (this *Custody) ListDailyAveragePrice(ctx *gin.Context) {
	dailyAveragePrice, err := ListDailyAveragePrice()
	if err != nil {
		this.Error(ctx, fmt.Sprintf("获取价格信息失败: %s", err.Error()))
		return
	}

	this.SuccessWithData(ctx, "获取成功", dailyAveragePrice)
}

// 数据库操作

// 列出托管信息
func ListCustodyInfo() ([]model.CustodyInfo, error) {
	var custodyInfos []model.CustodyInfo
	err := model.NewCustodyInfoModel().Find(&custodyInfos).Error
	if err != nil {
		return nil, err
	}

	return custodyInfos, nil
}

// createCustodyInfo 新增托管信息
func createCustodyInfo(data CustodyInfo) (*model.CustodyInfo, error) {
	ci := model.CustodyInfo{
		VenueName:       data.VenueName,
		SubAccountName:  data.SubAccountName,
		ObserverLink:    data.ObserverLink,
		EnergyRatio:     data.EnergyRatio,
		BasicHostingFee: data.BasicHostingFee,
	}

	if err := model.NewCustodyInfoModel().Create(&ci).Error; err != nil {
		return nil, err // 返回错误
	}

	// 返回插入的 CustodyInfo 实例
	return &ci, nil
}

// deleteCustodyInfoById 删除托管信息
func deleteCustodyInfoById(custodyInfoID string) error {
	return model.NewCustodyInfoModel().Where("id = ?", custodyInfoID).Delete(&model.CustodyInfo{}).Error
}

// updateCustodyInfo 更新托管信息
func updateCustodyInfo(data CustodyInfoUpdate) error {
	updates := map[string]interface{}{
		"venue_name":        data.VenueName,
		"sub_account_name":  data.SubAccountName,
		"observer_link":     data.ObserverLink,
		"energy_ratio":      data.EnergyRatio,
		"basic_hosting_fee": data.BasicHostingFee,
	}
	return model.NewCustodyInfoModel().Where("id = ?", data.ID).
		Updates(updates).Error
}

func CreateCustodyStatistics(data model.CustodyStatistics) error {
	return model.NewCustodyStatisticsModel().Create(&data).Error
}

func ListCustodyInfoWithTimeRange(startTime time.Time) (custodyStatistics []model.CustodyStatistics, err error) {
	println("ListCustodyInfoWithTimeRange:;;;;")
	if !startTime.IsZero() {
		// 添加时间过滤条件
		reportDateStr := startTime.Format("2006-01-02")
		println("startTime:", reportDateStr)
		err = model.NewCustodyStatisticsModel().Where("report_date >= ?", reportDateStr).Preload("CustodyInfo").Find(&custodyStatistics).Error
	} else {
		err = model.NewCustodyStatisticsModel().Preload("CustodyInfo").Find(&custodyStatistics).Error
	}

	return custodyStatistics, err
}

func ListDailyAveragePrice() ([]model.DailyAveragePrice, error) {
	var dailyAveragePrice []model.DailyAveragePrice
	err := model.NewDailyAveragePrice().Find(&dailyAveragePrice).Error
	if err != nil {
		return nil, err
	}

	return dailyAveragePrice, nil
}
