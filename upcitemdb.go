package qiaoqiao

type UpcItemDbResult struct {
	Code   string `json:"code"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Items  []UpcItemDbItem `json:"items"`
}

type UpcItemDbItem struct {
	Ean         string `json:"ean"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Isbn        string `json:"isbn"`
	Publisher   string `json:"publisher"`
	Images      []string `json:"images"`
}

func (p *UpcItemDbResult) getStatus() (status int) {
	if p.Code == "INVALID_UPC" || p.Total == 0 || len(p.Items) < 1 {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *UpcItemDbResult) getProduct() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Items[0].Title
}
func (p *UpcItemDbResult) getDescription() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return p.Items[0].Description
}
func (p *UpcItemDbResult) getPeople() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return ""
}
func (p *UpcItemDbResult) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return ""
	}
	return generateBarcodeUrl(p.Items[0].Ean)
}
func (p *UpcItemDbResult) getCompany() Company {
	if p.getStatus() == StatusRequestUnsuccessfully {
		return Company{"", ""}
	}
	return Company{p.Items[0].Publisher, ""}
}
