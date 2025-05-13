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
