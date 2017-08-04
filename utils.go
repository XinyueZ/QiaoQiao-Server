package qiaoqiao

func generateBarcodeUrl(code string) (s string) {
	s = "http://www.searchupc.com/drawupc.aspx?q=" + code
	return
}
