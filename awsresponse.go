package qiaoqiao

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
	ItemAttributes   *ItemAttributes
	OfferSummary     OfferSummary
	Offers           Offers
	SalesRank        int
	SmallImage       *Image
	MediumImage      *Image
	LargeImage       *Image
	ImageSets        *ImageSets
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
	Manufacturer    *string
	Publisher       *string
	NumberOfItems   int
	PackageQuantity int
	Feature         string
	Model           string
	ProductGroup    string
	ReleaseDate     string
	Studio          *string
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

type ImageSets struct {
	ImageSet []ImageSet
}

type ImageSet struct {
	//Category string `xml:"Category,attr"`
	Category       string `xml:",attr"`
	SwatchImage    *Image
	SmallImage     *Image
	ThumbnailImage *Image
	TinyImage      *Image
	MediumImage    *Image
	LargeImage     *Image
}

func (p *ItemLookupResponse) parse(productQuery *ProductQuery) IProductResult {
	return p
}

func (p *ItemLookupResponse) getStatus() (status int) {
	if !p.Items.Request.IsValid || p.Items.Item.ItemAttributes == nil {
		status = StatusRequestUnsuccessfully
	} else {
		status = StatusRequestSuccessfully
	}
	return
}

func (p *ItemLookupResponse) getProduct() string {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item.ItemAttributes.Title
	}
	return ""
}
func (p *ItemLookupResponse) getDescription() (desc string) {
	if p.getStatus() == StatusRequestSuccessfully {
		return p.Items.Item.ItemAttributes.Title + "\n" +
			p.Items.Item.EditorialReviews.EditorialReview.Content
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
		return generateBarcodeUrl(p.Items.Request.ItemLookupRequest.ItemID)
	}
	return ""
}
func (p *ItemLookupResponse) getCompany() Company {
	if p.getStatus() == StatusRequestSuccessfully {
		var companyName string = ""
		if p.Items.Item.ItemAttributes.Publisher != nil {
			companyName = *p.Items.Item.ItemAttributes.Publisher
		} else if p.Items.Item.ItemAttributes.Studio != nil {
			companyName = *p.Items.Item.ItemAttributes.Studio
		} else if p.Items.Item.ItemAttributes.Manufacturer != nil {
			companyName = *p.Items.Item.ItemAttributes.Manufacturer
		} else {
			companyName = ""
		}
		return Company{
			companyName,
			p.Items.Item.DetailPageURL,
		}
	}
	return Company{"", ""}
}

func (p *ItemLookupResponse) getProductImage() (imageList []ProductImage) {
	imageList = make([]ProductImage, 0)
	if p.getStatus() == StatusRequestSuccessfully {
		if p.Items.Item.ImageSets.ImageSet != nil && len(p.Items.Item.ImageSets.ImageSet) > 0 {
			pi := ProductImage{make([]string, 0), make([]string, 0), make([]string, 0), "", "aws"}
			for _, element := range p.Items.Item.ImageSets.ImageSet {
				if element.SmallImage != nil {
					pi.Small = append(pi.Small, element.SmallImage.URL)
				}
				if element.MediumImage != nil {
					pi.Medium = append(pi.Medium, element.MediumImage.URL)
				}
				if element.LargeImage != nil {
					pi.Large = append(pi.Large, element.LargeImage.URL)
				}
				if element.ThumbnailImage != nil {
					pi.Thumbnail = element.ThumbnailImage.URL
				}
				imageList = append(imageList, pi)
			}
		} else {
			pi := ProductImage{make([]string, 0), make([]string, 0), make([]string, 0), "", "aws"}
			if p.Items.Item.SmallImage != nil {
				pi.Thumbnail = p.Items.Item.SmallImage.URL
				pi.Small = append(pi.Small, p.Items.Item.SmallImage.URL)
			}
			if p.Items.Item.LargeImage != nil {
				pi.Large = append(pi.Large, p.Items.Item.LargeImage.URL)
			}
			if p.Items.Item.MediumImage != nil {
				pi.Medium = append(pi.Medium, p.Items.Item.MediumImage.URL)
			}
			imageList = append(imageList, pi)
		}
	}
	return
}
