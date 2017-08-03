//Package amazonproduct provides methods for interacting with the Amazon Product Advertising API
package qiaoqiao

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
ItemLookup takes a product ID (ASIN) and returns the result
*/
func (api AmazonProductAPI) ItemLookup(ItemId string) (string, error) {
	params := map[string]string{
		"ItemId":        ItemId,
		"ResponseGroup": "Images,ItemAttributes,Small,EditorialReview",
	}

	return api.genSignAndFetch("ItemLookup", params)
}

/*
ItemLookupWithResponseGroup takes a product ID (ASIN) and a ResponseGroup and returns the result
*/
func (api AmazonProductAPI) ItemLookupWithResponseGroup(ItemId string, ResponseGroup string) (string, error) {
	params := map[string]string{
		"ItemId":        ItemId,
		"ResponseGroup": ResponseGroup,
	}

	return api.genSignAndFetch("ItemLookup", params)
}

/*
ItemLookupWithParams takes the params for ItemLookup and returns the result
*/
func (api AmazonProductAPI) ItemLookupWithParams(params map[string]string) (string, error) {
	_, present := params["ItemId"]
	if !present {
		return "", errors.New("ItemId property is required in the params map")
	}

	return api.genSignAndFetch("ItemLookup", params)
}

/*
MultipleItemLookup takes an array of product IDs (ASIN) and returns the result
*/
func (api AmazonProductAPI) MultipleItemLookup(ItemIds []string) (string, error) {
	params := map[string]string{
		"ItemId":        strings.Join(ItemIds, ","),
		"ResponseGroup": "Images,ItemAttributes,Small,EditorialReview",
	}

	return api.genSignAndFetch("ItemLookup", params)
}

/*
MultipleItemLookupWithResponseGroup takes an array of product IDs (ASIN) as well as a ResponseGroup and returns the result
*/
func (api AmazonProductAPI) MultipleItemLookupWithResponseGroup(ItemIds []string, ResponseGroup string) (string, error) {
	params := map[string]string{
		"ItemId":        strings.Join(ItemIds, ","),
		"ResponseGroup": ResponseGroup,
	}

	return api.genSignAndFetch("ItemLookup", params)
}

/*
ItemSearchByKeyword takes a string containing keywords and returns the search results
*/
func (api AmazonProductAPI) ItemSearchByKeyword(Keywords string, page int) (string, error) {
	params := map[string]string{
		"Keywords":      Keywords,
		"ResponseGroup": "Images,ItemAttributes,Small,EditorialReview",
		"ItemPage":      strconv.FormatInt(int64(page), 10),
	}
	return api.ItemSearch("All", params)
}

func (api AmazonProductAPI) ItemSearchByKeywordWithResponseGroup(Keywords string, ResponseGroup string) (string, error) {
	params := map[string]string{
		"Keywords":      Keywords,
		"ResponseGroup": ResponseGroup,
	}
	return api.ItemSearch("All", params)
}

func (api AmazonProductAPI) ItemSearch(SearchIndex string, Parameters map[string]string) (string, error) {
	Parameters["SearchIndex"] = SearchIndex
	return api.genSignAndFetch("ItemSearch", Parameters)
}

/*
CartCreate takes a map containing ASINs and quantities. Up to 10 items are allowed
*/
func (api AmazonProductAPI) CartCreate(items map[string]int) (string, error) {

	params := make(map[string]string)

	i := 1
	for k, v := range items {
		if i < 11 {
			key := fmt.Sprintf("Item.%d.ASIN", i)
			params[key] = string(k)

			key = fmt.Sprintf("Item.%d.Quantity", i)
			params[key] = strconv.Itoa(v)

			i++
		} else {
			break
		}
	}
	return api.genSignAndFetch("CartCreate", params)
}

/*
CartAdd takes a map containing ASINs and quantities and adds them to the given cart.
Up to 10 items are allowed
*/
func (api AmazonProductAPI) CartAdd(items map[string]int, cartid, HMAC string) (string, error) {

	params := map[string]string{
		"CartId": cartid,
		"HMAC":   HMAC,
	}

	i := 1
	for k, v := range items {
		if i < 11 {
			key := fmt.Sprintf("Item.%d.ASIN", i)
			params[key] = string(k)

			key = fmt.Sprintf("Item.%d.Quantity", i)
			params[key] = strconv.Itoa(v)

			i++
		} else {
			break
		}
	}
	return api.genSignAndFetch("CartAdd", params)
}

/*
CartClear takes a CartId and HMAC that were returned when generating a cart
It then removes the contents of the cart
*/
func (api AmazonProductAPI) CartClear(CartId, HMAC string) (string, error) {

	params := map[string]string{
		"CartId": CartId,
		"HMAC":   HMAC,
	}

	return api.genSignAndFetch("CartClear", params)
}

/*
Cart get takes a CartID and HMAC that were returned when generating a cart
Returns the contents of the specified cart
*/
func (api AmazonProductAPI) CartGet(CartId, HMAC string) (string, error) {

	params := map[string]string{
		"CartId": CartId,
		"HMAC":   HMAC,
	}

	return api.genSignAndFetch("CartGet", params)
}

/*
BrowseNodeLookup takes a BrowseNodeId and returns the result.
*/
func (api AmazonProductAPI) BrowseNodeLookup(nodeId string) (string, error) {
	params := map[string]string{
		"BrowseNodeId": nodeId,
	}
	return api.genSignAndFetch("BrowseNodeLookup", params)
}

func (api AmazonProductAPI) BrowseNodeLookupWithResponseGroup(nodeId string, responseGroup string) (string, error) {
	params := map[string]string{
		"BrowseNodeId":  nodeId,
		"ResponseGroup": responseGroup,
	}
	return api.genSignAndFetch("BrowseNodeLookup", params)
}

// Response describes the generic API Response
type AWSResponse struct {
	OperationRequest struct {
		RequestID             string     `xml:"RequestId"`
		Arguments             []Argument `xml:"Arguments>Argument"`
		RequestProcessingTime float64
	}
}

// Argument todo
type Argument struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

// Image todo
type Image struct {
	URL    string
	Height uint16
	Width  uint16
}

// Price describes the product price as
// Amount of cents in CurrencyCode
type Price struct {
	Amount         uint
	CurrencyCode   string
	FormattedPrice string
}

type TopSeller struct {
	ASIN  string
	Title string
}

// Item represents a product returned by the API
type Item struct {
	ASIN             string
	URL              string
	DetailPageURL    string
	ItemAttributes   ItemAttributes
	OfferSummary     OfferSummary
	Offers           Offers
	SalesRank        int
	SmallImage       Image
	MediumImage      Image
	LargeImage       Image
	EditorialReviews EditorialReviews
	BrowseNodes struct {
		BrowseNode []BrowseNode
	}
}

// BrowseNode represents a browse node returned by API
type BrowseNode struct {
	BrowseNodeID string `xml:"BrowseNodeId"`
	Name         string
	TopSellers struct {
		TopSeller []TopSeller
	}
	Ancestors struct {
		BrowseNode []BrowseNode
	}
}

// ItemAttributes response group
type ItemAttributes struct {
	Author          string
	Binding         string
	Brand           string
	Color           string
	EAN             string
	Creator         string
	Title           string
	ListPrice       Price
	Manufacturer    string
	Publisher       string
	NumberOfItems   int
	PackageQuantity int
	Feature         string
	Model           string
	ProductGroup    string
	ReleaseDate     string
	Studio          string
	Warranty        string
	Size            string
	UPC             string
}

// Offer response attribute
type Offer struct {
	Condition       string `xml:"OfferAttributes>Condition"`
	ID              string `xml:"OfferListing>OfferListingId"`
	Price           Price  `xml:"OfferListing>Price"`
	PercentageSaved uint   `xml:"OfferListing>PercentageSaved"`
	Availability    string `xml:"OfferListing>Availability"`
}

// Offers response group
type Offers struct {
	TotalOffers     int
	TotalOfferPages int
	MoreOffersURL   string  `xml:"MoreOffersUrl"`
	Offers          []Offer `xml:"Offer"`
}

// OfferSummary response group
type OfferSummary struct {
	LowestNewPrice   Price
	LowerUsedPrice   Price
	TotalNew         int
	TotalUsed        int
	TotalCollectible int
	TotalRefurbished int
}

// EditorialReview response attribute
type EditorialReview struct {
	Source  string
	Content string
}

// EditorialReviews response group
type EditorialReviews struct {
	EditorialReview EditorialReview
}

// BrowseNodeLookupRequest is the confirmation of a BrowseNodeInfo request
type BrowseNodeLookupRequest struct {
	BrowseNodeId  string
	ResponseGroup string
}

// ItemLookupRequest is the confirmation of a ItemLookup request
type ItemLookupRequest struct {
	IDType        string `xml:"IdType"`
	ItemID        string `xml:"ItemId"`
	ResponseGroup string `xml:"ResponseGroup"`
	VariationPage string
}

// ItemLookupResponse describes the API response for the ItemLookup operation
type ItemLookupResponse struct {
	AWSResponse
	Items struct {
		Request struct {
			IsValid           bool
			ItemLookupRequest ItemLookupRequest
		}
		Item Item `xml:"Item"`
	}
}

// ItemSearchRequest is the confirmation of a ItemSearch request
type ItemSearchRequest struct {
	Keywords      string `xml:"Keywords"`
	SearchIndex   string `xml:"SearchIndex"`
	ResponseGroup string `xml:"ResponseGroup"`
}

type ItemSearchResponse struct {
	AWSResponse
	Items struct {
		Request struct {
			IsValid           bool
			ItemSearchRequest ItemSearchRequest
		}
		Items                []Item `xml:"Item"`
		TotalResult          int
		TotalPages           int
		MoreSearchResultsUrl string
	}
}

type BrowseNodeLookupResponse struct {
	AWSResponse
	BrowseNodes struct {
		Request struct {
			IsValid                 bool
			BrowseNodeLookupRequest BrowseNodeLookupRequest
		}
		BrowseNode BrowseNode
	}
}

func (p *ItemLookupResponse) getStatus() (status int) {
	if p.Items.Request.IsValid {
		status = StatusRequestSuccessfully
	} else {
		status = StatusRequestUnsuccessfully
	}
	return StatusRequestSuccessfully
}

func (p *ItemLookupResponse) getProduct() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item.ItemAttributes.Title
	}
	return ""
}
func (p *ItemLookupResponse) getDescription() (desc string) {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item.ItemAttributes.Title
	}
	return ""
}
func (p *ItemLookupResponse) getPeople() (people string) {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item.ItemAttributes.Author
	}
	return ""
}
func (p *ItemLookupResponse) getBarcodeUrl() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return "http://www.searchupc.com/drawupc.aspx?q=" + p.Items.Request.ItemLookupRequest.ItemID
	}
	return ""
}
func (p *ItemLookupResponse) getCompany() Company {
	if p.getStatus() == StatusRequestSuccessfully {
		return Company{
			p.Items.Item.ItemAttributes.Publisher,
			p.Items.Item.DetailPageURL,
		}
	}
	return Company{"", ""}
}
