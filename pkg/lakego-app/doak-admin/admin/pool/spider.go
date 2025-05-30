package pool

import (
	"fmt"
	"strings"
)

func SpiderBTCProfit(url string) ([]*HashRateEntry, error) {
	if strings.Contains(url, "antpool") {
		rUrl, err := convertProfitURL(url)
		if err != nil {
			return nil, err
		}

		hss, err := fetchAntPoolBTCProfit(rUrl)
		if err != nil {
			return nil, err
		}

		return hss, nil
	} else if strings.Contains(url, "f2pool") {
		hs, err := fetchF2PoolRecv(url)
		if err != nil {
			return nil, err
		}
		return hs, nil
	}

	return nil, fmt.Errorf("no router path")
}

func SpiderFBProfit(url string) ([]FBProfitEntry, error) {
	var (
		profits    []FBProfitEntry
		err        error
		antPoolUrl string
	)

	if strings.Contains(url, "antpool") {
		antPoolUrl, err = convertURLForFB(url)
		if err != nil {
			return nil, err
		}
		profits, err = fetchAntPoolFBProfit(antPoolUrl)
	} else if strings.Contains(url, "f2pool") {
		profits, err = fetchF2poolFBProfit(url)
	}

	return profits, nil
}

func SpiderRealTimeStatus(url string) (*RealTimeStatus, error) {
	if strings.Contains(url, "antpool") {
		antPoolUrl, err := convertHashRateURL(url)
		if err != nil {
			return nil, err
		}
		status, err := fetchAntPoolHashRate(antPoolUrl)
		if err != nil {
			return nil, err
		}

		workerStatusUrl, err := convertStatusURL(url)
		if err != nil {
			return nil, err
		}
		online, offline, err := fetchAntPoolWorkerStatus(workerStatusUrl)
		if err != nil {
			return nil, err
		}
		status.Online = online
		status.Offline = offline

		return status, nil

	} else if strings.Contains(url, "f2pool") {
		status, err := fetchF2PoolRealTimeStatus(url)
		if err != nil {
			return nil, err
		}

		return status, nil
	}

	return nil, fmt.Errorf("no router path")
}
