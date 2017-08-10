package qiaoqiao

const baseWikiUrl = "https://%s.wikipedia.org/%s%s"
const urlWikiGeosearch = "w/api.php?format=json&action=query&list=geosearch&gsradius=10000&gslimit=max&gscoord="
const urlWikiImages = "w/api.php?format=json&action=query&prop=pageimages&piprop=original|thumbnail&pilimit=1&redirects=titles&titles="
const urlWikiThumbnails = "w/api.php?format=json&action=query&prop=pageimages&piprop=original|thumbnail&pilimit=1&redirects=titles&titles="
const urlWikiId = "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&pageids="
const urlWikiDocuments = "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&titles="

const defaultImage = "http://www.skillsforlifefoundation.com/images/default-thumbnail.jpg"

const eandataUrl = "http://eandata.com/feed/?v=3&mode=json&find=%s&keycode=%s"

const awsSerachResponseGroup = "Medium"

const searchupcUrl = "http://www.searchupc.com/handlers/upcsearch.ashx?request_type=3&upc=%s&access_token=%s"

const upcitemdbUrl = "https://api.upcitemdb.com/prod/trial/lookup?upc=%s"

const barcodableUrl = "https://www.barcodable.com/api/v1/upc/%s"
