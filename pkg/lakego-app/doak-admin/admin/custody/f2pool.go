package custody

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type F2PoolLastHashResponse struct {
	Status string             `json:"status"`
	Data   F2PoolLastHashData `json:"data"`
}

type F2PoolLastHashData struct {
	IncomeData       []IncomeData     `json:"income_data"`
	FilterCommentMap FilterCommentMap `json:"filter_comment_map"`
}

type IncomeData struct {
	HashRate      interface{} `json:"hash_rate"`
	CreatedAt     json.Number `json:"created_at"`
	Comment       string      `json:"comment"`
	FilterComment string      `json:"filter_comment"`
	Type          string      `json:"type"`
	Amount        float64     `json:"amount"`
	TxFee         float64     `json:"txfee"`
	CurrencyCode  string      `json:"currency_code"`
	Difficulty    string      `json:"difficulty"`
}

type FilterCommentMap struct {
	CurrencyCode string            `json:"currency_code"`
	FilterMap    map[string]string `json:"filter_map"`
}

func fetchF2PoolRecv(url string) ([]*HashRateEntry, error) {
	newUrl := fmt.Sprintf("%s%s", url, "&currency_code=btc&action=load_payout_history_income")

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, newUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	//req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "www.f2pool.kz")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp F2PoolLastHashResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}

	if len(apiResp.Data.IncomeData) == 0 {
		return nil, fmt.Errorf("empty item")
	}

	var result []*HashRateEntry

	for _, item := range apiResp.Data.IncomeData {
		timestampStr := item.CreatedAt.String()
		timestamp, err := strconv.ParseFloat(timestampStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		lastDayTime := time.Unix((int64(timestamp)), 0).Format("2006-01-02")

		btcRecv := item.Amount + item.TxFee

		hashRate, unit, err := getHashRateString(item.HashRate)
		if err != nil {
			return nil, fmt.Errorf("error getting hash rate: %v", err)
		}

		result = append(result, &HashRateEntry{
			LastDayHashRate: hashRate,
			LastDayRecv:     fmt.Sprintf("%.8f", btcRecv),
			LastDayHashUnit: unit,
			LastDayTime:     lastDayTime,
		})
	}

	//apiResp.OriginData.Summary.HashRate
	return result, nil
}

func getHashRateString(hashRate interface{}) (string, string, error) {
	switch v := hashRate.(type) {
	case string:
		// 尝试拆分字符串（空格分隔）
		parts := strings.Fields(v)
		if len(parts) == 2 {
			if parts[1] == "Thash/s" {
				return parts[0], "TH/s", nil
			}
			return parts[0], parts[1], nil
		}
		return "", "", fmt.Errorf("error getting hash rate string")
	case float64:
		return fmt.Sprintf("%g", v), "TH/s", nil // 如果是数字，格式化为字符串
	default:
		return "", "", fmt.Errorf("invalid hash rate format")
	}
}
