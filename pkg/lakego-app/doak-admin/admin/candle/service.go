package candle

import (
	"encoding/json"
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const CandleURL = "https://api.exchange.coinbase.com/products/BTC-USD/candles?granularity=3600"

func Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	fmt.Println("start run candle....")

	for {
		data, err := fetchCandles()
		if err != nil {
			fmt.Println(err.Error())
		}
		updateCandlesToDB(data)
		updateYesterdaysAveragePrice()
		<-ticker.C
	}
}

func updateYesterdaysAveragePrice() {
	// 定义 CST 时区
	//loc, err := time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	fmt.Println("加载时区失败:", err)
	//	return
	//}

	now := time.Now().UTC()
	if now.Hour() < 1 {
		// 当前时间小于1点，不执行
		fmt.Println("当前时间未到1点，不执行任务")
		return
	}

	yesterday := now.AddDate(0, 0, -1)
	start := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	candles, err := controller.GetCandleRange(start, end)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 计算昨天的平均收盘价
	var totalPrice float64
	var count int
	for _, candle := range candles {
		price, err := strconv.ParseFloat(candle.PriceClose, 64)
		if err != nil {
			// 转换失败，可选择跳过或记录日志
			fmt.Println("转换价格失败:", err, "价格:", candle.PriceClose)
			return
		}
		totalPrice += price
		count++
	}

	avgPrice := totalPrice / float64(count)
	date := start.Format("2006-01-02")
	exist, err := controller.IsDailyAveragePriceExist(date)
	if err != nil {
		fmt.Println("查找日期错误：", err.Error())
		return
	}
	if !exist {
		err := controller.CreateDailyAveragePrice(model.DailyAveragePrice{
			Date:        date,
			CstAvgPrice: "0.00",
			UtcAvgPrice: fmt.Sprintf("%.2f", avgPrice),
		})
		if err != nil {
			println(err.Error())
			return
		}
	}
}

type CandleData [][]float64

func fetchCandles() ([][]float64, error) {
	resp, err := http.Get(CandleURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码：%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var candles [][]float64
	err = json.Unmarshal(body, &candles)
	if err != nil {
		return nil, err
	}

	return candles, nil
}

func updateCandlesToDB(candles [][]float64) {
	maxTimestamp, err := controller.GetMaxTimestamp()
	if err != nil {
		maxTimestamp = time.Now().AddDate(0, 0, -30)
	}

	reverseSlice(candles)

	for _, candle := range candles {
		timestamp := time.Unix(int64(candle[0]), 0)

		// 只插入比已存在最大时间戳更早的记录
		if !maxTimestamp.IsZero() && !timestamp.After(maxTimestamp) {
			continue
		}

		priceLow := fmt.Sprintf("%.2f", candle[1])
		priceHigh := fmt.Sprintf("%.2f", candle[2])
		priceOpen := fmt.Sprintf("%.2f", candle[3])
		priceClose := fmt.Sprintf("%.2f", candle[4])

		err := controller.CreateCandle(model.BtcUsdCandle{
			Timestamp:  timestamp,
			PriceLow:   priceLow,
			PriceHigh:  priceHigh,
			PriceOpen:  priceOpen,
			PriceClose: priceClose,
		})
		if err != nil {
			println(err.Error())
		}
	}
}

// reverse2DFloat64Slice 将二维切片倒序
func reverseSlice(slice [][]float64) {
	left, right := 0, len(slice)-1
	for left < right {
		slice[left], slice[right] = slice[right], slice[left]
		left++
		right--
	}
}
