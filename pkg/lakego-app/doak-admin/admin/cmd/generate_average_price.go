package cmd

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/deatil/lakego-doak/lakego/command"
	"strconv"
)

var GenerateAveragePriceCmd = &command.Command{
	Use:          "lakego-admin:generate-average-price",
	Short:        "lakego-admin generate-average-price.",
	Example:      "{execfile} lakego-admin:generate-average-price",
	SilenceUsage: true,
	PreRun: func(cmd *command.Command, args []string) {

	},
	Run: func(cmd *command.Command, args []string) {
		err := GenerateAveragePrice()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func GenerateAveragePrice() error {
	candles, err := controller.CandleList()
	if err != nil {
		return err
	}

	priceMap := make(map[string][]float64)

	for _, candle := range candles {
		utcTime := candle.Timestamp.UTC()
		dataKey := utcTime.Format("2006-01-02")

		// 调整 UTC 0 点的数据到前一天的 key
		if utcTime.Hour() == 0 {
			// 将 key 设为前一天
			dataKey = utcTime.AddDate(0, 0, -1).Format("2006-01-02")
		}

		price, err := strconv.ParseFloat(candle.PriceClose, 64)
		if err != nil {
			fmt.Println("收盘价转换失败：", err.Error())
			continue
		}
		priceMap[dataKey] = append(priceMap[dataKey], price)
		fmt.Println("日期：", dataKey, "时间：", candle.Timestamp, "UTC 时间", utcTime)
	}

	for dataKey, prices := range priceMap {
		if len(prices) < 24 {
			fmt.Println(dataKey, "数据不完整，跳过，长度：", len(prices))
			continue
		}
		var total float64
		for _, price := range prices {
			total += price
		}
		avgPrice := total / float64(len(prices))

		exist, err := controller.IsDailyAveragePriceExist(dataKey)
		if err != nil {
			continue
		}
		if !exist {
			// 创建平均价格记录
			err = controller.CreateDailyAveragePrice(model.DailyAveragePrice{
				Date:        dataKey,
				CstAvgPrice: "0.0",                         // 可以根据需要存储 CST 平均价格
				UtcAvgPrice: fmt.Sprintf("%.2f", avgPrice), // 或者存储 UTC 平均价格
			})
			if err != nil {
				fmt.Println("记录平均价格失败:", err)
				return err
			}
		}
	}

	return nil
}
