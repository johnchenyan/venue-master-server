package pool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var AntBaseURL string = "https://www.antpool.com/auth/v3/observer/api/hash/query"

type AntPoolRecvResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data RecvData `json:"data"`
}

type RecvData struct {
	TotalRecv   string        `json:"totalRecv"`
	Items       []RecvItem    `json:"items"`
	PageNum     float64       `json:"pageNum"`
	TotalPage   float64       `json:"totalPage"`
	PageSize    float64       `json:"pageSize"`
	TotalRecord float64       `json:"totalRecord"`
	AuxCoinList []interface{} `json:"auxCoinList"`
}

type RecvItem struct {
	CreateDate              json.Number   `json:"creatDate"`
	DayHashRate             string        `json:"dayHashRate"`
	DayHashRateUnit         string        `json:"dayHashRateUnit"`
	DayRecv                 string        `json:"dayRecv"`
	Type                    string        `json:"type"`
	PlusPercent             string        `json:"plusPercent"`
	PayStatus               string        `json:"payStatus"`
	PpaPpsAmount            string        `json:"ppaPpsAmount"`
	PpaPplnsAmount          string        `json:"ppaPplnsAmount"`
	MevAmount               string        `json:"mevAmount"`
	MevPercent              string        `json:"mevPercent"`
	UserHashrate            string        `json:"userHashrate"`
	UserHashrateUnit        string        `json:"userHashrateUnit"`
	IsContractUser          bool          `json:"isContractUser"`
	IsContractError         bool          `json:"isContractError"`
	OutContractHashRate     string        `json:"outContractHashRate"`
	OutContractHashRateUnit string        `json:"outContractHashRateUnit"`
	InContractHashRate      string        `json:"inContractHashRate"`
	InContractHashRateUnit  string        `json:"inContractHashRateUnit"`
	InContractModelList     []interface{} `json:"inContractModelList"`
	SilenceSwitcher         bool          `json:"silenceSwitcher"`
	OutContractModelList    []interface{} `json:"outContractModelList"`
}

// ****************** FB 结构体 ***************************** //

type PaymentItem struct {
	PayId           json.Number `json:"payId"`
	CreatDate       json.Number `json:"creatDate"`
	Amount          string      `json:"amount"`
	Status          string      `json:"status"`
	WalletAddress   string      `json:"walletAddress"`
	TransactionHash *string     `json:"transactionHash,omitempty"`
	TransactionLink *string     `json:"transactionLink,omitempty"`
	Type            string      `json:"type"`
	SubPayList      *string     `json:"subPayList,omitempty"`
}

type FBPayData struct {
	Items       []PaymentItem `json:"items"`
	AuxCoinList []string      `json:"auxCoinList"`
	PageNum     float64       `json:"pageNum"`
	TotalPage   float64       `json:"totalPage"`
	PageSize    float64       `json:"pageSize"`
	TotalRecord float64       `json:"totalRecord"`
}

type AntPoolFBPayApiResponse struct {
	Code string    `json:"code"`
	Msg  string    `json:"msg"`
	Data FBPayData `json:"data"`
}

// 24小时哈希和实时哈希结构体
type ApiResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	UserID                     string  `json:"userId"`
	NickName                   string  `json:"nickName"`
	HSNow                      string  `json:"hsNow"`
	HSNowUnit                  string  `json:"hsNowUnit"`
	HSLast1D                   string  `json:"hsLast1D"`
	HSLast1DUnit               string  `json:"hsLast1DUnit"`
	UserHSNow                  string  `json:"userHsNow"`
	UserHSNowUnit              string  `json:"userHsNowUnit"`
	UserHSLast1D               string  `json:"userHsLast1D"`
	UserHSLast1DUnit           string  `json:"userHsLast1DUnit"`
	IsContractUser             bool    `json:"isContractUser"`
	IsContractError            bool    `json:"isContractError"`
	IsContractPercent          bool    `json:"isContractPercent"`
	ContractHashrateUnit       string  `json:"contractHashrateUnit"`
	ContractHashrate           string  `json:"contractHashrate"`
	ContractPercent            string  `json:"contractPercent"`
	RealContractHashrateUnit   string  `json:"realContractHashrateUnit"`
	RealContractHashrate       string  `json:"realContractHashrate"`
	RealContractHashrate1DUnit string  `json:"realContractHashrate1DUnit"`
	RealContractHashrate1D     string  `json:"realContractHashrate1D"`
	Time                       float64 `json:"time"`
	TotalWorkerNum             float64 `json:"totalWorkerNum"`
	OnlineWorkerNum            float64 `json:"onlineWorkerNum"`
	OfflineWorkerNum           float64 `json:"offlineWorkerNum"`
	DisableWorkerNum           float64 `json:"disableWorkerNum"`
	Hash1DBaseUnit             float64 `json:"hash1DBaseUnit"`
	Hash5mBaseUnit             float64 `json:"hash5mBaseUnit"`
	LastDayIntegrateHS         string  `json:"lastDayIntegrateHs"`
	LastDayIntegrateHSUnit     string  `json:"lastDayIntegrateHsUnit"`
	LastDayIntegrateHSBaseUnit float64 `json:"lastDayIntegrateHsBaseUnit"`
	SilenceSwitcher            bool    `json:"silenceSwitcher"`
}

// worker状态结构体

// Worker 结构体定义了每个worker的字段
type AntPoolWorker struct {
	WorkerId          string  `json:"workerId"`
	HsLast10Min       string  `json:"hsLast10Min"`
	HsLast1Hour       string  `json:"hsLast1Hour"`
	HsLast1H          string  `json:"hsLast1H"`
	HsLast1D          string  `json:"hsLast1D"`
	RejectRatio       string  `json:"rejectRatio"`
	ShareLastTime     float64 `json:"shareLastTime"`
	WorkerStatus      float64 `json:"workerStatus"`
	UserWorkerId      string  `json:"userWorkerId"`
	GroupId           float64 `json:"groupId"`
	OnlineTimeLast24h float64 `json:"onlineTimeLast24h"`
	ReconnectLast24h  float64 `json:"reconnectLast24h"`
	BurstBlock72h     bool    `json:"burstBlock72h"`
	Id                float64 `json:"id"`
	CreateTime        float64 `json:"createTime"`
}

// WorkerStatus 结构体定义了worker状态的字段
type AntPoolWorkerStatus struct {
	TotalWorkerNum   float64 `json:"totalWorkerNum"`
	OnlineWorkerNum  float64 `json:"onlineWorkerNum"`
	OfflineWorkerNum float64 `json:"offlineWorkerNum"`
	DisableWorkerNum float64 `json:"disableWorkerNum"`
}

// Data 结构体定义了返回JSON中的data字段
type AntPoolWorkerStatusData struct {
	Items        []AntPoolWorker     `json:"items"`
	PageNum      float64             `json:"pageNum"`
	TotalPage    float64             `json:"totalPage"`
	PageSize     float64             `json:"pageSize"`
	TotalRecord  float64             `json:"totalRecord"`
	WorkerStatus AntPoolWorkerStatus `json:"workerStatus"`
}

// Response 结构体定义了整个返回的JSON数据结构
type AntPoolWorkerStatusResponse struct {
	Code string                  `json:"code"`
	Msg  string                  `json:"msg"`
	Data AntPoolWorkerStatusData `json:"data"`
}

// convertURL 将url1的查询部分附加到新的基础URL上
func convertProfitURL(sUrl string) (string, error) {
	parsedURL, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	accessKey := queryParams.Get("accessKey")
	observerUserId := queryParams.Get("observerUserId")

	dUrl := fmt.Sprintf("https://www.antpool.com/auth/v3/observer/api/recv?accessKey=%s&coinType=BTC&observerUserId=%s&pageNum=1&pageSize=1000",
		accessKey, observerUserId)

	return dUrl, nil
}

func fetchAntPoolBTCProfit(url string) ([]*HashRateEntry, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp AntPoolRecvResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(apiResp.Data.Items) == 0 {
		return nil, fmt.Errorf("empty item")
	}

	var result []*HashRateEntry
	recvMap := make(map[string]*HashRateEntry) // 用于存储每个日期对应的信息

	for _, item := range apiResp.Data.Items {
		timestampStr := item.CreateDate.String()
		timestamp, err := strconv.ParseFloat(timestampStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		lastDayTime := time.Unix((int64(timestamp))/1000, 0).Format("2006-01-02")

		// 将 DayRecv 从字符串转换为 float64
		dayRecvValue, err := strconv.ParseFloat(item.DayRecv, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting DayRecv to float64: %v", err)
		}

		// 检查该日期是否已经存在于 map 中
		if entry, exists := recvMap[lastDayTime]; exists {
			// 累加 DayRecv 值并转换为字符串
			currentRecv, _ := strconv.ParseFloat(entry.LastDayRecv, 64)
			newTotalRecv := currentRecv + dayRecvValue
			entry.LastDayRecv = fmt.Sprintf("%.8f", newTotalRecv) // 转换为字符串并保留两位小数
		} else {
			// 如果不存在，创建新的 HashRateEntry 记录
			recvMap[lastDayTime] = &HashRateEntry{
				LastDayHashRate: item.DayHashRate,
				LastDayHashUnit: item.DayHashRateUnit,
				LastDayRecv:     item.DayRecv, // 初始值直接使用 DayRecv 字符串
				LastDayTime:     lastDayTime,
			}
		}

		// 构建结果列表
		for _, entry := range recvMap {
			result = append(result, entry)
		}
	}

	return result, nil
}

// https://www.antpool.com/auth/v3/observer/api/pay?accessKey=SelaP5Guh9SZ4orfNKHK&coinType=BTC&observerUserId=AMWV02&pageNum=1&pageSize=40&mergeCoinType=FB
func convertURLForFB(sUrl string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	accessKey := queryParams.Get("accessKey")
	observerUserId := queryParams.Get("observerUserId")

	dUrl := fmt.Sprintf("https://www.antpool.com/auth/v3/observer/api/pay?accessKey=%s&coinType=BTC&observerUserId=%s&pageNum=1&pageSize=1000&mergeCoinType=FB",
		accessKey, observerUserId)

	return dUrl, nil
}

func fetchAntPoolFBProfit(url string) ([]FBProfitEntry, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp AntPoolFBPayApiResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	if len(apiResp.Data.Items) == 0 {
		return nil, fmt.Errorf("empty item")
	}

	var result []FBProfitEntry

	for _, item := range apiResp.Data.Items {
		timestampStr := item.CreatDate.String()
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

// convertURL 将url1的查询部分附加到新的基础URL上
func convertHashRateURL(sUrl string) (string, error) {

	// 分割url1以获取查询部分
	parts := strings.Split(sUrl, "?")
	if len(parts) != 2 {
		return "", fmt.Errorf("the provided URL does not contain a query part")
	}

	// 解析基础URL
	u, err := url.Parse(AntBaseURL)
	if err != nil {
		return "", err
	}

	// 设置查询部分
	u.RawQuery = parts[1]

	// 返回新的URL字符串
	return u.String(), nil
}

func fetchAntPoolHashRate(url string) (*RealTimeStatus, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp ApiResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &RealTimeStatus{
		RealTimeHash: apiResp.Data.HSNow,
		HSLast1D:     apiResp.Data.HSLast1D,
		HashUnit:     apiResp.Data.HSNowUnit,
	}, nil
}

// convertURL 将url1的查询部分附加到新的基础URL上
func convertStatusURL(sUrl string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	accessKey := queryParams.Get("accessKey")
	observerUserId := queryParams.Get("observerUserId")

	dUrl := fmt.Sprintf("https://www.antpool.com/auth/v3/observer/api/worker/list?search=&workerStatus=0&accessKey=%s&coinType=BTC&observerUserId=%s&pageNum=1&pageSize=10",
		accessKey, observerUserId)

	return dUrl, nil
}

func fetchAntPoolWorkerStatus(url string) (int, int, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResp AntPoolWorkerStatusResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return 0, 0, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return int(apiResp.Data.WorkerStatus.OnlineWorkerNum), int(apiResp.Data.WorkerStatus.OfflineWorkerNum), nil
}
