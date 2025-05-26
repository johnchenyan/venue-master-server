package controller

// 场地相关数据结构

type Field struct {
	ID     int    `json:"ID"`
	Value  string `json:"value"`
	Status string `json:"status"` // "default" | "new" | "modified" | "deleted"
}

type VenueTemplateChange struct {
	ID                 uint    `json:"ID"`
	TemplateNameBefore string  `json:"templateNameBefore"`
	TemplateNameAfter  string  `json:"templateNameAfter"`
	Fields             []Field `json:"fields"`
}

// "default" | "new" | "modified" | "deleted";
const (
	VenueTemplateStatusDefault  string = "default"
	VenueTemplateStatusNew      string = "new"
	VenueTemplateStatusModified string = "modified"
	VenueTemplateStatusDeleted  string = "deleted"
)

type VenueTemplateNew struct {
	ID           uint    `json:"ID"`
	TemplateName string  `json:"templateName"`
	Fields       []Field `json:"fields"`
}

// VenueRecordNew templateID
type VenueRecordNew struct {
	TemplateID uint           `json:"templateID"`
	Fields     []FieldsRecord `json:"fields"`
}

type VenueRecordUpdate struct {
	RecordID uint           `json:"RecordID"`
	Fields   []FieldsRecord `json:"Fields"`
}

type FieldsRecord struct {
	ID        uint   `json:"id"`
	FieldName string `json:"fieldName"`
	Value     string `json:"value"`
}

// 托管相关数据结构
type CustodyInfo struct {
	ID              uint   `json:"id"`
	VenueName       string `json:"venue_name"`
	SubAccountName  string `json:"sub_account_name"`
	ObserverLink    string `json:"observer_link"`
	EnergyRatio     string `json:"energy_ratio"`
	BasicHostingFee string `json:"basic_hosting_fee"`
}

type CustodyInfoUpdate = CustodyInfo

type CustodyHostingFeeCurve struct {
	Year     string  `json:"year"`
	Value    float64 `json:"value"`
	Category string  `json:"category"`
}

type SettlementQueryParam struct {
	Type      string              `json:"type"`
	NameMap   map[string][]string `json:"name"`
	StartTime string              `json:"start"`
	EndTime   string              `json:"end"`
}

type SettlementQueryWithPaginationParam struct {
	Type      string              `json:"type"`
	NameMap   map[string][]string `json:"name"`
	StartTime string              `json:"start"`
	EndTime   string              `json:"end"`
	Price     string              `json:"price"`
	Page      int                 `json:"page"`      // 当前页码
	PageSize  int                 `json:"page_size"` // 每页条目数
}

type SettlementQueryWithPaginationResult struct {
	Data     []SettlementItem `json:"data"`      // 查询结果数据列表
	Total    int64            `json:"total"`     // 数据总条目数
	Page     int              `json:"page"`      // 当前页码
	PageSize int              `json:"page_size"` // 每页条目数
}

// SettlementItem 表示单个条目的结构体
type SettlementItem struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

type SettlementQueryResult struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	TimeRange  string `json:"time_range"`
	TimeLength string `json:"time_length"` // 分钟
}

type SettlementQueryResultT struct {
	Name       string `json:"name"`
	TimeRange  string `json:"time_range"`
	TimeLength string `json:"time_length"` // 时
}

type SettlementQueryAverageResult struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	TimeRange string  `json:"time_range"`
	Average   float64 `json:"average"`
}

type SettlementQueryAverageResultT struct {
	Name      string  `json:"name"`
	TimeRange string  `json:"time_range"`
	Average   float64 `json:"average"`
}

//
//type SettlementQueryEntry struct {
//	Name       string `json:"name"`
//	Type       string `json:"start"`
//	TimeRange  string `json:"time_range"`
//	TimeLength string `json:"time_length"` // 分钟
//}
