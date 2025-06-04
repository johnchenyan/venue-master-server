package mining

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/deatil/lakego-doak-admin/admin/pool"
	"golang.org/x/xerrors"
	"math"
	"strconv"
	"strings"
	"time"
)

func Run() {
	go fetchProfitLoop()
	go fetchStatusLoop()
	go miningPoolChannelLoop()
}

func miningPoolChannelLoop() {
	for request := range controller.MiningPoolChannel {
		err := processBTCProfit(request.MiningPool)
		if err != nil {
			fmt.Printf(err.Error())
		}

		err = processFBProfit(request.MiningPool)
		if err != nil {
			fmt.Printf(err.Error())
		}

		err = processStatus(request.MiningPool)
		if err != nil {
			fmt.Printf(err.Error())
		}

		// 将处理结果发送回结果通道
		if err != nil {
			request.ResultChan <- err
		} else {
			request.ResultChan <- nil // 处理成功
		}
	}
}

func fetchProfitLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	fmt.Println("start run fetch btc mining pool record....")

	for {

		miningPools, err := controller.ListBtcMiningPool()
		if err != nil {
			fmt.Printf("获取矿池失败: %v\n", err)
			time.Sleep(time.Second * 60)
			continue
		}

		for _, miningPool := range miningPools {
			if !miningPool.IsEnabled {
				continue
			}

			err := processBTCProfit(miningPool)
			if err != nil {
				fmt.Printf(err.Error())
			}

			err = processFBProfit(miningPool)
			if err != nil {
				fmt.Printf(err.Error())
			}

		}

		<-ticker.C
	}
}

func fetchStatusLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	fmt.Println("start run fetch btc mining pool status.......")

	for {

		miningPools, err := controller.ListBtcMiningPool()
		if err != nil {
			fmt.Printf("获取矿池失败: %v\n", err)
			time.Sleep(time.Second * 60)
			continue
		}

		for _, miningPool := range miningPools {
			if !miningPool.IsEnabled {
				continue
			}

			err := processStatus(miningPool)
			if err != nil {
				fmt.Printf(err.Error())
			}
		}

		<-ticker.C
	}
}

func processBTCProfit(miningPool model.MiningPool) error {
	btcProfits, err := pool.SpiderBTCProfit(miningPool.Link)
	if err != nil {
		return xerrors.Errorf("spider name: %v failed：%v\n", miningPool.PoolName, err)
	}

	for _, hs := range btcProfits {
		date, err := parseLastDayTime(miningPool.Link, hs.LastDayTime, false)
		if err != nil {
			fmt.Printf("parseLastDayTime name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		exist, err := controller.IsMiningRecordExist(miningPool.ID, date)
		if err != nil {
			fmt.Printf("IsMiningRecordExist name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		if exist {
			continue
		}

		hash, err := parseHashRate(hs.LastDayHashRate, hs.LastDayHashUnit)
		if err != nil {
			fmt.Printf("parseHashRate name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		btcRecv, err := parseProfitToFloat64(hs.LastDayRecv)
		if err != nil {
			fmt.Printf("parseProfitBtcToFloat64 name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		err = controller.CreateBtcMiningSettlementRecord(model.MiningSettlementRecord{
			PoolID:                        miningPool.ID,
			SettlementDate:                date,
			SettlementTheoreticalHashrate: miningPool.TheoreticalHashrate,
			SettlementHashrate:            hash,
			SettlementProfitBtc:           btcRecv,
			SettlementProfitFb:            0,
		})
		if err != nil {
			fmt.Printf("CreateBtcMiningSettlementRecord name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}
	}

	return nil
}

func processFBProfit(miningPool model.MiningPool) error {
	FBProFits, err := pool.SpiderFBProfit(miningPool.Link)
	if err != nil {
		return xerrors.Errorf("SpiderFBProfit error: %w", err)
	}
	for _, fp := range FBProFits {
		date, err := parseLastDayTime(miningPool.Link, fp.LastDayTime, true)
		if err != nil {
			fmt.Printf("parseLastDayTime name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		updated, err := controller.IsFBRecordUpdated(miningPool.ID, date)
		if err != nil {
			fmt.Printf("IsFBRecordUpdated name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		if updated {
			continue
		}

		fbRecv, err := parseProfitToFloat64(fp.LastDayRecv)
		if err != nil {
			fmt.Printf("parseProfitBtcToFloat64 name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}

		if fbRecv == 0 {
			continue
		}

		err = controller.UpdateBtcMiningFBProfit(miningPool.ID, date, fbRecv)
		if err != nil {
			fmt.Printf("UpdateBtcMiningFBProfit name: %v failed：%v\n", miningPool.PoolName, err)
			continue
		}
	}

	return nil
}

func processStatus(miningPool model.MiningPool) error {
	status, err := pool.SpiderRealTimeStatus(miningPool.Link)
	if err != nil {
		return xerrors.Errorf("spider name: %v failed：%v\n", miningPool.PoolName, err)
	}

	realTimeHash, err := getHashRate(status.RealTimeHash, status.HashUnit)
	if err != nil {
		return err
	}

	hsLast1D, err := getHashRate(status.HSLast1D, status.HashUnit)
	if err != nil {
		return err
	}

	err = controller.CreateBtcMiningPoolStatus(model.MiningPoolStatus{
		PoolID:          miningPool.ID,
		CurrentHashrate: realTimeHash,
		Last24hHashrate: hsLast1D,
		OnlineMachines:  status.Online,
		OfflineMachines: status.Offline,
	})
	if err != nil {
		return xerrors.Errorf("CreateBtcMiningPoolStatus name: %v failed：%v\n", miningPool.PoolName, err)
	}
	return nil
}

func parseHashRate(hashRate, hashUnit string) (float64, error) {
	var value float64
	hr, err := strconv.ParseFloat(hashRate, 64)
	if err != nil {
		return 0, err
	}
	switch hashUnit {
	case "TH/s":
		value = hr
	case "PH/s":
		value = hr * 1000
	case "EH/s":
		value = hr * 1000 * 1000
	}
	return value, nil
}

func parseLastDayTime(link, lastDayTime string, isFB bool) (string, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, lastDayTime)
	if err != nil {
		return "", fmt.Errorf("解析时间失败: %w", err)
	}

	if strings.Contains(link, "f2pool") || (strings.Contains(link, "antpool") && isFB) {
		t = t.AddDate(0, 0, -1)
	}

	return t.Format(layout), nil
}

func parseProfitToFloat64(profit string) (float64, error) {
	value, err := strconv.ParseFloat(profit, 64)
	if err != nil {
		return 0, fmt.Errorf("无法解析利润字符串: %s, 错误: %v", profit, err)
	}

	// 保留 8 位小数
	roundedValue := math.Round(value*1e8) / 1e8
	return roundedValue, nil
}

func getHashRate(hashStr, unit string) (float64, error) {
	hash, err := strconv.ParseFloat(hashStr, 64)
	if err != nil {
		return 0, xerrors.Errorf("ParseFloat err: %v", err)
	}
	switch unit {
	case "H/s":
		// 不需要转换
		return hash, nil
	case "GH/s":
		return hash * 1e9, nil
	case "TH/s":
		return hash * 1e12, nil
	case "PH/s":
		return hash * 1e15, nil
	case "EH/s":
		return hash * 1e18, nil
	default:
		return 0, fmt.Errorf("不支持的单位: %s", unit)
	}
}
