package custody

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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

// convertURL 将url1的查询部分附加到新的基础URL上
func convertURL(sUrl string) (string, error) {
	parsedURL, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	accessKey := queryParams.Get("accessKey")
	observerUserId := queryParams.Get("observerUserId")

	dUrl := fmt.Sprintf("https://www.antpool.com/auth/v3/observer/api/recv?accessKey=%s&coinType=BTC&observerUserId=%s&pageNum=1&pageSize=40",
		accessKey, observerUserId)

	return dUrl, nil
}

func fetchAntPoolRecv(url string) ([]*HashRateEntry, error) {
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

	for _, item := range apiResp.Data.Items {
		timestampStr := item.CreateDate.String()
		timestamp, err := strconv.ParseFloat(timestampStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		lastDayTime := time.Unix((int64(timestamp))/1000, 0).Format("2006-01-02")

		result = append(result, &HashRateEntry{
			LastDayHashRate: item.DayHashRate,
			LastDayHashUnit: item.DayHashRateUnit,
			LastDayRecv:     item.DayRecv,
			LastDayTime:     lastDayTime,
		})
	}

	return result, nil
}
