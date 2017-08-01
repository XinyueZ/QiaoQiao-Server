package qiaoqiao

type EANdataAttributes struct {
	Product         string `json:"product"`
	Description     string `json:"description"`
	LongDescription string `json:"long_desc"`
	Language        string `json:"language_text"`
	LanguageText    string `json:"language_text_long"`
	Price           string `json:"price_new"`
	PriceUnit       string `json:"price_new_extra"`
}

type Barcode struct {
	Url string `json:"EAN13"`
}

type Company struct {
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type ENAdataStatus struct {
	Code string `json:"code"`
}

type Product struct {
	Attributes EANdataAttributes `json:"attributes"`
	EAN13      string `json:"EAN13"`
	ISBN10     string `json:"ISBN10"`
	Barcode    Barcode `json:"ISBN10"`
	Image      string `json:"image"`
}

type EANdataResult struct {
	Status  ENAdataStatus `json:"code"`
	Product Product  `json:"product"`
	Company Company `json:"company"`
}
