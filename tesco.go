package qiaoqiao

type TescoResult struct {
	TescoProduct []TescoProduct  `json:"products"`
}

type TescoProduct struct {
	UPC              string `json:"gtin"`
	TescoId          string `json:"catId"`
	ShortDescription string `json:"description"`
	LongDescription  string `json:"marketingText"`
	Brand            string `json:"brand"`
}
