package controller

import (
	"github.com/deatil/lakego-doak-admin/admin/model"
	"time"
)

type Candle struct {
	Base
}

//func (this *Candle) CreateCandle(ctx *gin.Context) {
//	//custodyInfos, err := ListCustodyInfo()
//	//if err != nil {
//	//	this.Error(ctx, fmt.Sprintf("获取托管信息失败: %s", err.Error()))
//	//	return
//	//}
//	//
//	//this.SuccessWithData(ctx, "获取成功", custodyInfos)
//}

// 列出托管信息
func CreateCandle(candle model.BtcUsdCandle) error {
	return model.NewBtcUsdCandle().Create(&candle).Error
}

func GetMaxTimestamp() (time.Time, error) {
	var maxTimestamp time.Time
	err := model.NewBtcUsdCandle().Select("MAX(`timestamp`)").Scan(&maxTimestamp).Error
	if err != nil {
		return maxTimestamp, err
	}

	return maxTimestamp, nil
}

func GetCandleRange(start, end time.Time) ([]model.BtcUsdCandle, error) {
	var candles []model.BtcUsdCandle
	err := model.NewBtcUsdCandle().
		Where("timestamp > ? AND timestamp <= ?", start, end).
		Find(&candles).Error
	if err != nil {
		return nil, err
	}
	return candles, nil
}
