package controller

import (
	"errors"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"gorm.io/gorm"
)

type BtcDailyAveragePrice struct {
	Base
}

// 创建数据
func CreateDailyAveragePrice(averagePrice model.DailyAveragePrice) error {
	return model.NewDailyAveragePrice().Create(&averagePrice).Error
}

// 判断某个日期是否存在
func IsDailyAveragePriceExist(date string) (bool, error) {
	var record model.DailyAveragePrice
	err := model.NewDailyAveragePrice().Where("date = ?", date).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到，说明不存在
			return false, nil
		}
		// 其他错误
		return false, err
	}
	// 找到，存在
	if record.Date != date {
		return false, nil
	}
	return true, nil
}

func GetDailyAveragePrice(date string) (string, bool, error) {
	var record model.DailyAveragePrice

	err := model.NewDailyAveragePrice().Where("date = ?", date).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return "", false, nil
		}
		// 其他错误
		return "", false, err
	}

	if record.Date != date {
		return "", false, nil
	}

	// 记录存在，返回价格
	return record.UtcAvgPrice, true, nil
}
