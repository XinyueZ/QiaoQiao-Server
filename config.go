package qiaoqiao

const baseWikiUrl = "https://%s.wikipedia.org/%s%s"
const defaultImage = "http://www.skillsforlifefoundation.com/images/default-thumbnail.jpg"
const urlWikiImages = "w/api.php?format=json&action=query&prop=pageimages&piprop=original|thumbnail&pilimit=1&redirects=titles&titles="
const urlWikiDocuments = "w/api.php?format=json&action=query&prop=extracts|pageimages|langlinks&llprop=autonym&lldir=descending&lllimit=500&piprop=original|name|thumbnail&exlimit=1&redirects=titles&titles="
const urlWikiGeosearch = "w/api.php?format=json&action=query&list=geosearch&gsradius=10000&gslimit=max&gscoord="
