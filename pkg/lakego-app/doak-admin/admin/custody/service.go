package custody

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"golang.org/x/xerrors"
	"strconv"
	"strings"
	"time"
)

func Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	defer close(controller.CustodyInfoChannel)

	fmt.Println("start run custody....")

	// 启动一个 goroutine 来处理 collect 调用
	go func() {
		for custodyInfo := range controller.CustodyInfoChannel {
			if err := collect(custodyInfo); err != nil {
				fmt.Println("场地：", custodyInfo.VenueName, "子账户：", custodyInfo.SubAccountName, err.Error())
			}
		}
	}()

	for {

		custodyInfos, err := controller.ListCustodyInfo()
		if err != nil {
			fmt.Printf("获取托管信息失败: %v\n", err)
			continue
		}

		fmt.Println("custody info:", len(custodyInfos))

		// 处理函数
		for _, custodyInfo := range custodyInfos {
			// 获取观察链接算力、btc收益
			if err := collect(custodyInfo); err != nil {
				fmt.Println("场地： ", custodyInfo.VenueName, "子账户：", custodyInfo.SubAccountName, err.Error())
			}
		}

		<-ticker.C
	}

}

func collect(custodyInfo model.CustodyInfo) error {
	hss, err := spider(custodyInfo.ObserverLink)
	if err != nil {
		return xerrors.Errorf("获取数据失败：%s", err.Error())
	}

	for _, hs := range hss {
		err := processHashRate(hs, custodyInfo)
		if err != nil {
			fmt.Println(time.Now(), "场地:", custodyInfo.VenueName, "子账户:", custodyInfo.SubAccountName, err.Error())
			continue
		}
	}

	return nil
}

func processHashRate(hs *HashRateEntry, custodyInfo model.CustodyInfo) error {
	// 计算总能耗
	energy, err := totalEnergy(*hs, custodyInfo)
	if err != nil {
		return xerrors.Errorf("计算总能耗失败：%s", err.Error())
	}

	// 计算总的托管费
	totalHostingFee, err := totalHostingFee(energy, custodyInfo)
	if err != nil {
		return xerrors.Errorf("计算总托管费失败：%s", err.Error())
	}

	// 转换时间，主要是针对鱼池
	lastDayTime, err := convertUTC(custodyInfo.ObserverLink, hs.LastDayTime)
	if err != nil {
		return xerrors.Errorf("转换时间失败：%s", err.Error())
	}

	// 判断日平均价格是否存在，不存在返回
	averagePrice, exist, err := controller.GetDailyAveragePrice(lastDayTime)
	if err != nil {
		return xerrors.Errorf("获取日平均价格失败：%s", err.Error())
	}
	if !exist {
		return xerrors.Errorf("日均价格记录不存在")
	}

	totalIncomeUSD, err := totalIncomeUSD(hs.LastDayRecv, averagePrice)
	if err != nil {
		return xerrors.Errorf("计算总USD收益失败：%s", err.Error())
	}

	netIncome := totalIncomeUSD - totalHostingFee

	totalHostingFeeRatio := totalHostingFee / totalIncomeUSD * 100

	err = controller.CreateCustodyStatistics(model.CustodyStatistics{
		CustodyID:            custodyInfo.ID,
		EnergyRatio:          custodyInfo.EnergyRatio,
		BasicHostingFee:      custodyInfo.BasicHostingFee,
		HourlyComputingPower: hs.LastDayHashRate,
		TotalHostingFee:      fmt.Sprintf("%.2f", totalHostingFee),
		TotalIncomeBTC:       hs.LastDayRecv,
		TotalIncomeUSD:       fmt.Sprintf("%.2f", totalIncomeUSD),
		NetIncome:            fmt.Sprintf("%.2f", netIncome),
		HostingFeeRatio:      fmt.Sprintf("%.2f%%", totalHostingFeeRatio),
		ReportDate:           hs.LastDayTime,
	})
	if err != nil {
		return xerrors.Errorf("插入统计数据失败：%s", err.Error())
	}

	return nil
}

func totalEnergy(hs HashRateEntry, custodyInfo model.CustodyInfo) (float64, error) {
	var value float64
	hashRate, err := strconv.ParseFloat(hs.LastDayHashRate, 64)
	if err != nil {
		return 0, err
	}
	switch hs.LastDayHashUnit {
	case "TH/s":
		value = hashRate
	case "PH/s":
		value = hashRate * 1000
	case "EH/s":
		value = hashRate * 1000 * 1000
	}

	energyRatio, err := strconv.ParseFloat(custodyInfo.EnergyRatio, 64)
	if err != nil {
		return 0, err
	}

	return value * energyRatio * 3600 * 24, nil
}

func totalHostingFee(energy float64, custodyInfo model.CustodyInfo) (float64, error) {
	energyKwh := joulesToKWh(energy)
	basicHostingFee, err := strconv.ParseFloat(custodyInfo.BasicHostingFee, 64)
	if err != nil {
		return 0, err
	}
	return basicHostingFee * energyKwh, nil
}

func joulesToKWh(joules float64) float64 {
	return joules / 3600000
}

func totalIncomeUSD(incomeBtc, averagePrice string) (float64, error) {
	incomeBtcFloat, err := strconv.ParseFloat(incomeBtc, 64)
	if err != nil {
		return 0, err
	}

	averagePriceFloat, err := strconv.ParseFloat(averagePrice, 64)
	if err != nil {
		return 0, err
	}
	return averagePriceFloat * incomeBtcFloat, nil
}

func convertUTC(link, lastDayTime string) (string, error) {
	if strings.Contains(link, "antpool") {
		return lastDayTime, nil
	}

	layout := "2006-01-02"
	t, err := time.Parse(layout, lastDayTime)
	if err != nil {
		return "", fmt.Errorf("解析时间失败: %w", err)
	}

	if strings.Contains(link, "f2pool") {
		t = t.AddDate(0, 0, -1)
	}

	return t.Format(layout), nil
}
