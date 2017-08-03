package qiaoqiao

import (
	"encoding/xml"
	"time"
	"fmt"
	"sort"
	"strings"
	"crypto/hmac"
	"crypto/sha256"
	"net/url"

	"encoding/base64"
)

const timestampFormat = "2006-01-02T15:04:05Z"

var timeNowFunc = time.Now

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

func (p *ItemLookupResponse) getStatus() (status int) {
	if len(p.Items.Item) > 0 {
		status = StatusRequestSuccessfully
	} else {
		status = StatusRequestUnsuccessfully
	}
	return
}

func (p *ItemLookupResponse) getProduct() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item[0].ItemAttributes.Title
	}
	return ""
}
func (p *ItemLookupResponse) getDescription() (desc string) {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item[0].ItemAttributes.Title
	}
	return ""
}
func (p *ItemLookupResponse) getPeople() (people string) {
	if p.getStatus() == StatusRequestSuccessfully {
		people = ""
		for _, s := range p.Items.Item[0].ItemAttributes.Author {
			people = people + s + " "
		}
		return
	}
	return ""
}
func (p *ItemLookupResponse) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item[0].DetailPageURL
	}
	return ""
}
func (p *ItemLookupResponse) getCompany() Company {
	if p.getStatus() == StatusRequestSuccessfully {
		return Company{
			p.Items.Item[0].ItemAttributes.Manufacturer,
			p.Items.Item[0].DetailPageURL,
		}
	}
	return Company{"", ""}
}

type ItemLink struct {
	Description string
	URL         string
}

type ItemLinks struct {
	ItemLink []ItemLink
}

type Image struct {
	URL    string
	Height Size
	Width  Size
}

type Size struct {
	Value int    `xml:",chardata"`
	Units string `xml:",attr"`
}

type ImageSets struct {
	ImageSet []ImageSet
}

type ImageSet struct {
	Category       string `xml:",attr"`
	SwatchImage    Image
	SmallImage     Image
	ThumbnailImage Image
	TinyImage      Image
	MediumImage    Image
	LargeImage     Image
}

type ItemAttributes struct {
	Author            []string
	Artist            string
	Actor             string
	AspectRatio       string
	AudienceRating    string
	Binding           string
	Creator           Creator
	EAN               string
	EANList           EANList
	CatalogNumberList CatalogNumberList
	Format            []string
	IsAdultProduct    bool
	ISBN              string
	Label             string
	Languages         Languages
	ListPrice         Price
	Manufacturer      string
	NumberOfPages     int
	PackageDimensions PackageDimensions
	ProductGroup      string
	ProductTypeName   string
	PublicationDate   *Date
	PackageQuantity   int
	PartNumber        string
	UPC               string
	UPCList           UPCList
	Publisher         string
	Studio            string
	Title             string
	NumberOfDiscs     []int
}

type Date struct {
	time.Time
}

type Creator struct {
	Role string `xml:",attr"`
	Name string `xml:",chardata"`
}

type EANList struct {
	Element []string `xml:"EANListElement"`
}

type Languages struct {
	Language []Language
}

type Language struct {
	Name        string
	Type        string
	AudioFormat string
}

type Price struct {
	Amount         string
	CurrencyCode   string
	FormattedPrice string
}

type PackageDimensions struct {
	Height Size
	Length Size
	Weight Size
	Width  Size
}

type OfferSummary struct {
	LowestNewPrice   Price
	LowestUsedPrice  Price
	TotalNew         int
	TotalUsed        int
	TotalCollectible int
	TotalRefurbished int
}

type Offers struct {
	TotalOffers     int
	TotalOfferPages int
	MoreOffersURL   string `xml:"MoreOffersUrl"`
	Offer           []Offer
}

type Offer struct {
	OfferAttributes OfferAttributes
	OfferListing    OfferListing
	LoyaltyPoints   LoyaltyPoints
	Merchant        Merchant
}

type Merchant struct {
	Name string
}

type OfferAttributes struct {
	Condition string
}

type OfferListing struct {
	ID                              string `xml:"OfferListingId"`
	Price                           Price
	Availability                    string
	AvailabilityAttributes          AvailabilityAttributes
	IsEligibleForSuperSaverShipping bool
	IsEligibleForPrime              bool
}

type AvailabilityAttributes struct {
	AvailabilityType string
	MinimumHours     int
	MaximumHours     int
}

type LoyaltyPoints struct {
	Points                 int
	TypicalRedemptionValue Price
}

type CustomerReviews struct {
	IFrameURL  string
	HasReviews bool
}

type SimilarProducts struct {
	SimilarProduct []SimilarProduct
}

type asinTitle struct {
	ASIN  string
	Title string
}

type SimilarProduct struct {
	asinTitle
}

type TopSellers struct {
	TopSeller []TopSeller
}

type TopSeller struct {
	asinTitle
}

type NewReleases struct {
	NewRelease []NewRelease
}

type NewRelease struct {
	asinTitle
}

type SimilarViewedProducts struct {
	SimilarViewedProduct []SimilarViewedProduct
}

type SimilarViewedProduct struct {
	asinTitle
}

type TopItemSet struct {
	Type    string
	TopItem []TopItem
}

type TopItem struct {
	asinTitle
	DetailPageURL string
	ProductGroup  string
	Author        string
}

type CatalogNumberList struct {
	Element []string `xml:"CatalogNumberListElement"`
}

type UPCList struct {
	Element []string `xml:"UPCListElement"`
}

type Item struct {
	XMLName         xml.Name `xml:"Item"`
	ASIN            string
	DetailPageURL   string
	SalesRank       int
	ItemLinks       ItemLinks
	SmallImage      Image
	MediumImage     Image
	LargeImage      Image
	ImageSets       ImageSets
	ItemAttributes  ItemAttributes
	OfferSummary    OfferSummary
	Offers          Offers
	CustomerReviews CustomerReviews
	SimilarProducts SimilarProducts
	BrowseNodes     BrowseNodes
}
type BrowseNodes struct {
	BrowseNode []BrowseNode
	Request    Request
}

type Request struct {
	IsValid bool
	Errors  *Errors
}

type Errors struct {
	XMLName   xml.Name `xml:"Errors"`
	ErrorNode []Error  `xml:"Error"`
}

type ErrorCode string

type Error struct {
	Code    ErrorCode
	Message string
}

func (e Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("Error %v: %v", e.Code, e.Message)
	}
	return ""
}

func (e Errors) Error() string {
	if len(e.ErrorNode) > 0 {
		return e.ErrorNode[0].Error()
	}
	return ""
}

type BrowseNode struct {
	ID         string `xml:"BrowseNodeId"`
	Name       string
	Ancestors  BrowseNodes
	Children   BrowseNodes
	TopSellers TopSellers
	TopItemSet []TopItemSet
}

type Items struct {
	TotalResults         int
	TotalPages           int
	MoreSearchResultsURL string `xml:"MoreSearchResultsUrl"`
	Item                 []Item
}

type ItemLookupResponse struct {
	XMLName xml.Name `xml:"ItemLookupResponse"`
	Items   Items    `xml:"Items"`
}

type AmazonAssociate struct {
	Tag  string
	Host string
}
