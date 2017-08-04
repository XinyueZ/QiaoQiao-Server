package qiaoqiao

type BarcodableResult struct {
	Status  int         `json:"status"`
	Message string         `json:"message"`
	Item    BarcodableItem `json:"item"`
}

type BarcodableItem struct {
	Ean   string           `json:"ean"`
	Isbn  string           `json:"isbn"`
	Asins []BarcodableAsin `json:"asins"`
}

type BarcodableAsin struct {
	Asin         string   `json:"asin"`
	Title        string   `json:"title"`
	Mpn          string   `json:"mpn"`
	PartNumber   string   `json:"part_number"`
	Brand        string   `json:"brand"`
	Manufacturer string   `json:"manufacturer"`
	Url          string   `json:"url"`
	Images       []string `json:"images"`
	Categories   []string `json:"categories"`
}

func (p *BarcodableResult) getStatus() (status int) {
	switch p.Status {
	case 200:
		status = StatusRequestSuccessfully
	default:
		status = StatusRequestUnsuccessfully
	}
	return
}

func (p *BarcodableResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Item.Asins[0].Title
}
func (p *BarcodableResult) getDescription() (desc string) {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	desc = p.getProduct()
	for _, s := range p.Item.Asins[0].Categories {
		desc = "\n" + desc + s + "\n"
	}
	return desc
}
func (p *BarcodableResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return ""
}
func (p *BarcodableResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return generateBarcodeUrl(p.Item.Ean)
}
func (p *BarcodableResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return Company{
		p.Item.Asins[0].Manufacturer + "\n" + p.Item.Asins[0].Brand,
		"",
	}
}
