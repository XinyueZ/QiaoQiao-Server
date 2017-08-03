package qiaoqiao

const baseWikiUrl = "https://%s.wikipedia.org/%s%s"
const urlWikiGeosearch = "w/api.php?format=json&action=query&list=geosearch&gsradius=10000&gslimit=max&gscoord="
const urlWikiImages = "w/api.php?format=json&action=query&prop=pageimages&piprop=original|thumbnail&pilimit=1&redirects=titles&titles="
const urlWikiThumbnails = "w/api.php?format=json&action=query&prop=pageimages&piprop=original|thumbnail&pilimit=1&redirects=titles&titles="
const urlWikiId = "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&pageids="
const urlWikiDocuments = "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&titles="

const defaultImage = "http://www.skillsforlifefoundation.com/images/default-thumbnail.jpg"

const eandataUrl = "http://eandata.com/feed/?v=3&mode=json&find=%s&keycode=%s"

const awsUrl = "http://%s%s?AWSAccessKeyId=%s&AssociateTag=%s&IdType=EAN&ItemId=%s&Operation=ItemLookup&SearchIndex=All&Service=AWSECommerceService&Timestamp=%s&Signature=%s"
