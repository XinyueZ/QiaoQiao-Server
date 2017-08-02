package qiaoqiao

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"google.golang.org/appengine"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"net/url"
	"sort"
)

type ProductUpc struct {
	r         *http.Request
	targetUrl string
}

func newProductUpc(r *http.Request, targetUrl string) (p *ProductUpc) {
	p = new(ProductUpc)
	p.r = r
	p.targetUrl = targetUrl
	return
}

func (p *ProductUpc) get(language string, code string, response chan []byte, service string) {
	switch service {
	case "eandata":
		get(p.r, fmt.Sprintf(p.targetUrl, code, EANDATE_KEY), response)
	case "aws":
		for _, assoc := range AWS_ASSOCIATE_LIST {
			get(p.r, getAWSapi(p, &assoc, code), response)
		}
	}
}

type ProductUpcResponse struct {
	r           *http.Request
	Status      int  `json:"status"`
	Product     string  `json:"product"`
	Description string `json:"description"`
	Barcode     string `json:"description"`
	Company     Company `json:"company"`
	People      string `json:"people"`
	Source      string `json:"source"`
}

func newProductUpcResponse(r *http.Request, eandata *EANdataResult) (p *ProductUpcResponse) {
	p = new(ProductUpcResponse)

	p.Source = "eandata"
	p.r = r
	p.Status, _ = strconv.Atoi(eandata.Status.Code)
	if p.Status == 404 {
		p.Status = StatusRequestUnsuccessfully
		return
	}

	p.Product = eandata.Product.Attributes.Product
	if eandata.Product.Attributes.LongDescription != "" {
		p.Description = eandata.Product.Attributes.LongDescription
	} else {
		p.Description = eandata.Product.Attributes.Description
	}
	p.People = eandata.Product.Attributes.Author
	p.Barcode = eandata.Product.Barcode.Url
	p.Company = eandata.Company
	return
}

func (p *ProductUpcResponse) show(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(p)
	if err == nil {
		fmt.Fprintf(w, "%s", bytes)
	} else {
		NewStatus(w, "noid", StatusRequestUnsuccessfully, "Can't give you UPC information.").show(appengine.NewContext(p.r))
	}
}

const timestampFormat = "2006-01-02T15:04:05Z"

var timeNowFunc = time.Now

type AmazonAssociate struct {
	Tag  string
	Host string
}

func escape(s string) (r string) {
	r = strings.Replace(url.QueryEscape(s), "+", "%20", -1)
	return
}

func getAWSapi(p *ProductUpc, associate *AmazonAssociate, searchedId string) (ep string) {
	tm := timeNowFunc().UTC().Format(timestampFormat)
	keyValue := make(map[string]string)
	keyValue["Service"] = "AWSECommerceService"
	keyValue["Operation"] = "ItemLookup"
	keyValue["AWSAccessKeyId"] = AWS_ACCESS_ID
	keyValue["AssociateTag"] = associate.Tag
	keyValue["ItemId"] = searchedId
	keyValue["IdType"] = "EAN"
	keyValue["SearchIndex"] = "All"
	keyValue["Timestamp"] = tm
	queryKeys := make([]string, 0, len(keyValue))
	for key := range keyValue {
		queryKeys = append(queryKeys, key) // "key1", "key2" , "key3" .....
	}
	sort.Strings(queryKeys)
	queryKeysAndValues := make([]string, len(queryKeys))
	for i, key := range queryKeys {
		queryKeysAndValues[i] = escape(key) + "=" + escape(keyValue[key])
	}
	query := strings.Join(queryKeysAndValues, "&")
	msg := "GET\n" +
		associate.Host + "\n" +
		AWS_PATH + "\n" +
		query
	mac := hmac.New(sha256.New, []byte(AWS_SECURITY_KEY))
	mac.Write([]byte(msg))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	ep = fmt.Sprintf(
		p.targetUrl,
		associate.Host,
		AWS_PATH,
		AWS_ACCESS_ID,
		associate.Tag,
		searchedId,
		tm,
		sig)
	return
}
