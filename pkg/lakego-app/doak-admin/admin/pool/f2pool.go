package pool

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

// **************** 上次结算收益 ********************* //

type OutCome struct {
	CreatedAt    json.Number `json:"created_at"`
	Address      string      `json:"address"`
	TxID         string      `json:"txid"`
	Chain        string      `json:"chain"`
	CurrencyCode string      `json:"currency_code"`
	Amount       string      `json:"amount"`
	Tax          float64     `json:"tax"`
	Type         string      `json:"type"`
}

type F2PoolOutComeResponse struct {
	Status string    `json:"status"`
	Data   []OutCome `json:"data"`
}

// worker 状态结构体
type Worker struct {
	LocalHash           string      `json:"local_hash"`
	LastShare           int64       `json:"last_share"`
	Hashrate            string      `json:"hashrate"`
	HashesAccepted      json.Number `json:"hashes_accepted"`
	SharesAccepted      int         `json:"shares_accepted"`
	HashesLastDay       json.Number `json:"hashes_last_day"`
	StaleHashesRejected json.Number `json:"stale_hashes_rejected"`
	StaleSharesRejected int         `json:"stale_shares_rejected"`
	Currency            string      `json:"currency"`
	StaleHashesLastDay  json.Number `json:"stale_hashes_last_day"`
	LocalHashesLastDay  json.Number `json:"local_hashes_last_day"`
	DelayHashesLastDay  json.Number `json:"delay_hashes_last_day"`
	IP                  string      `json:"ip"`
	GroupID             *int        `json:"group_id"`
	GroupName           *string     `json:"group_name"`
	Status              int         `json:"status"`
	DelayrateLastDay    string      `json:"delayrate_last_day"`
	Delayrate           float64     `json:"delayrate"`
	HashrateLastDay     string      `json:"hashrate_last_day"`
	Stalerate           float64     `json:"stalerate"`
	StalerateLastDay    string      `json:"stalerate_last_day"`
	LocalrateLastDay    string      `json:"localrate_last_day"`
	Name                string      `json:"name"`
	OriginName          string      `json:"origin_name"`
	Tag                 *string     `json:"tag"`
	TagName             *string     `json:"tag_name"`
	OnlineTime7d        float64     `json:"online_time_7d"`
	OnlineTime30d       float64     `json:"online_time_30d"`
	OnlineTime60d       float64     `json:"online_time_60d"`
}

// OriginData 结构体定义了originData中的所有字段
type OriginData struct {
	Workers             []Worker               `json:"workers"`
	Tags                []interface{}          `json:"tags"`
	TagsDict            map[string]interface{} `json:"tags_dict"`
	Limit               int                    `json:"limit"`
	Offset              int                    `json:"offset"`
	Tab                 string                 `json:"tab"`
	Summary             Summary                `json:"summary"`
	TagsOverview        []TagsOverview         `json:"tagsOverview"`
	Count               int                    `json:"count"`
	WorkerLengthAll     int                    `json:"worker_length_all"`
	WorkerLengthOnline  int                    `json:"worker_length_online"`
	WorkerLengthOffline int                    `json:"worker_length_offline"`
	WorkerLengthDead    int                    `json:"worker_length_dead"`
}

// Summary 结构体定义了summary中的所有字段
type Summary struct {
	HashRate      json.Number `json:"hash_rate"`
	HashRateDaily json.Number `json:"hash_rate_daily"`
	RejectRate    string      `json:"reject_rate"`
	LocalHash     json.Number `json:"local_hash"`
	DelayRate     json.Number `json:"delay_rate"`
}

// TagsOverview 结构体定义了tagsOverview中的所有字段
type TagsOverview struct {
	HashRate json.Number `json:"hash_rate"`
	Total    int         `json:"total"`
	Online   int         `json:"online"`
	Offline  int         `json:"offline"`
	Expire   int         `json:"expire"`
	TagName  *string     `json:"tag_name"`
}

// Response 结构体定义了JSON数据的顶层结构
type F2PoolResponse struct {
	Status          string     `json:"status"`
	Data            []Worker   `json:"data"`
	OriginData      OriginData `json:"originData"`
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
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

func fetchF2poolFBProfit(url string) ([]FBProfitEntry, error) {
	// https://www.f2pool.kz/mining-user/b4a9423f365f8204bd8b9f3cb9c0dd87?user_name=amkaz04&params=user_name=amkaz04&currency_code=fb-mm&account=amkaz04&action=load_payout_history_outcome
	// https://www.f2pool.kz/mining-user/b4a9423f365f8204bd8b9f3cb9c0dd87?user_name=amkaz04&params=user_name=amkaz04&currency_code=btc&account=amkaz04&action=load_payout_history_outcome
	newUrl := fmt.Sprintf("%s%s", url, "&currency_code=fb-mm&action=load_payout_history_outcome")
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
	var apiResp F2PoolOutComeResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}

	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var result []FBProfitEntry

	for _, item := range apiResp.Data {
		timestampStr := item.CreatedAt.String()
		timestamp, err := strconv.ParseFloat(timestampStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		lastDayTime := time.Unix((int64(timestamp))/1000, 0).Format("2006-01-02")

		result = append(result, FBProfitEntry{
			LastDayTime: lastDayTime,
			LastDayRecv: item.Amount,
		})
	}

	return result, nil

}

func fetchF2PoolRealTimeStatus(url string) (*RealTimeStatus, error) {
	// https://www.f2pool.kz/mining-user/3b9959eda2f92927e18951b032a55ca1?user_name=amkaz01
	// url := "https://www.f2pool.kz/mining-user/3b9959eda2f92927e18951b032a55ca1?user_name=amkaz01&currency=BTC&action=get_pagination_workers"
	// Send a GET request to the URL

	newUrl := fmt.Sprintf("%s%s", url, "&currency=BTC&action=get_pagination_workers")

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, newUrl, nil)

	if err != nil {
		fmt.Println(err)
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
	var apiResp F2PoolResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}

	realTimeHash, unit, err := convertF2PoolHashRate(apiResp.OriginData.Summary.HashRate.String())
	if err != nil {
		return nil, fmt.Errorf("error convertF2PoolHashRateL: %v", err)
	}

	hsLast1D, _, err := convertF2PoolHashRate(apiResp.OriginData.Summary.HashRateDaily.String())
	if err != nil {
		return nil, fmt.Errorf("error convertF2PoolHashRateL: %v", err)
	}

	//apiResp.OriginData.Summary.HashRate
	return &RealTimeStatus{
		RealTimeHash: realTimeHash,
		HSLast1D:     hsLast1D,
		HashUnit:     unit,
		Online:       apiResp.OriginData.WorkerLengthOnline,
		Offline:      apiResp.OriginData.WorkerLengthOffline,
	}, nil
}

// convertHashRate 函数将字符串形式的哈希率转换为PH/s或TH/s，并保留两位小数
func convertF2PoolHashRate(hashRateStr string) (string, string, error) {
	hashRate, err := strconv.ParseFloat(hashRateStr, 64)
	if err != nil {
		fmt.Println("Error converting hash rate:", err)
		return "", "", fmt.Errorf("error converting hash rate")
	}

	// 转换为PH/s或TH/s
	if hashRate >= 1e18 {
		return fmt.Sprintf("%.2f", hashRate/1e1), "EH/s", nil
	} else if hashRate >= 1e15 { // PH/s
		return fmt.Sprintf("%.2f", hashRate/1e15), "PH/s", nil
	} else if hashRate >= 1e12 { // TH/s
		return fmt.Sprintf("%.2f", hashRate/1e12), "TH/s", nil
	} else if hashRate >= 1e9 { // TH/s
		return fmt.Sprintf("%.2f", hashRate/1e9), "GH/s", nil
	} else {
		return fmt.Sprintf("%.2f", hashRate), "H/s", nil
	}
}
