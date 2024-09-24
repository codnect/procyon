package http

import (
	"mime"
	"strings"
)

// MediaType represents an HTTP media type (also known as MIME type).
type MediaType string

// These constants represent common HTTP media types.
const (
	MediaTypeAll               MediaType = "*/*"
	MediaTypeAtomXml           MediaType = "application/atom+xml"
	MediaTypeCbor              MediaType = "application/cbor"
	MediaTypeFormUrlencoded    MediaType = "application/x-www-form-urlencoded"
	MediaTypeGraphqlResponse   MediaType = "application/graphql-response+json"
	MediaTypeJson              MediaType = "application/json"
	MediaTypeJsonUtf8          MediaType = "application/json;charset=UTF-8"
	MediaTypeYaml              MediaType = "application/yaml"
	MediaTypeYamlUtf8          MediaType = "application/yaml;charset=UTF-8"
	MediaTypeOctetStream       MediaType = "application/octet-stream"
	MediaTypePdf               MediaType = "application/pdf"
	MediaTypeProblemJson       MediaType = "application/problem+json"
	MediaTypeProblemJsonUtf8   MediaType = "application/problem+json;charset=UTF-8"
	MediaTypeProblemXml        MediaType = "application/problem+xml"
	MediaTypeProtobuf          MediaType = "application/x-protobuf"
	MediaTypeRssXml            MediaType = "application/rss+xml"
	MediaTypeStreamJson        MediaType = "application/stream+json"
	MediaTypeXhtmlXml          MediaType = "application/xhtml+xml"
	MediaTypeXml               MediaType = "application/xml"
	MediaTypeImageGif          MediaType = "image/gif"
	MediaTypeImageJpeg         MediaType = "image/jpeg"
	MediaTypeImagePng          MediaType = "image/png"
	MediaTypeMultipartFormData MediaType = "multipart/form-data"
	MediaTypeMultipartMixed    MediaType = "multipart/mixed"
	MediaTypeMultipartRelated  MediaType = "multipart/related"
	MediaTypeTextEventStream   MediaType = "text/event-stream"
	MediaTypeTextHtml          MediaType = "text/html"
	MediaTypeTextMarkdown      MediaType = "text/markdown"
	MediaTypeTextPlain         MediaType = "text/plain"
	MediaTypeTextXml           MediaType = "text/xml"
)

// cacheMediaType is a struct used to cache parsed media type information.
type cacheMediaType struct {
	// typ is the type of the media type.
	typ string
	// subtype is the subtype of the media type.
	subtype string
	// charset is the charset parameter of the media type.
	charset string
	// params is a map of all other parameters of the media type.
	params map[string]string
}

// cacheMediaTypes is a map used to cache parsed media types.
var (
	cacheMediaTypes = map[MediaType]cacheMediaType{}
)

// init function parses all predefined media types and caches their information.
func init() {
	parseMediaType(MediaTypeAll)
	parseMediaType(MediaTypeAtomXml)
	parseMediaType(MediaTypeCbor)
	parseMediaType(MediaTypeFormUrlencoded)
	parseMediaType(MediaTypeGraphqlResponse)
	parseMediaType(MediaTypeJson)
	parseMediaType(MediaTypeJsonUtf8)
	parseMediaType(MediaTypeYaml)
	parseMediaType(MediaTypeYamlUtf8)
	parseMediaType(MediaTypeOctetStream)
	parseMediaType(MediaTypePdf)
	parseMediaType(MediaTypeProblemJson)
	parseMediaType(MediaTypeProblemJsonUtf8)
	parseMediaType(MediaTypeProblemXml)
	parseMediaType(MediaTypeProtobuf)
	parseMediaType(MediaTypeRssXml)
	parseMediaType(MediaTypeStreamJson)
	parseMediaType(MediaTypeXhtmlXml)
	parseMediaType(MediaTypeXml)
	parseMediaType(MediaTypeImageGif)
	parseMediaType(MediaTypeImageJpeg)
	parseMediaType(MediaTypeImagePng)
	parseMediaType(MediaTypeMultipartFormData)
	parseMediaType(MediaTypeMultipartFormData)
	parseMediaType(MediaTypeMultipartMixed)
	parseMediaType(MediaTypeMultipartMixed)
	parseMediaType(MediaTypeMultipartRelated)
	parseMediaType(MediaTypeTextEventStream)
	parseMediaType(MediaTypeTextEventStream)
	parseMediaType(MediaTypeTextHtml)
	parseMediaType(MediaTypeTextMarkdown)
	parseMediaType(MediaTypeTextPlain)
	parseMediaType(MediaTypeTextXml)
}

// parseMediaType parses a media type string and caches its information.
func parseMediaType(mediaType MediaType) {
	fullType, params, _ := mime.ParseMediaType(string(mediaType))
	typ, subtype, _ := strings.Cut(fullType, "/")

	cache := cacheMediaType{
		typ:     typ,
		subtype: subtype,
		charset: "",
		params:  params,
	}

	charset, exists := params["charset"]
	if exists {
		cache.charset = charset
	}

	cacheMediaTypes[mediaType] = cache
}

// Type method returns the type of the media type.
func (m MediaType) Type() string {
	if val, ok := cacheMediaTypes[m]; ok {
		return val.typ
	}

	return ""
}

// IsWildcardType method checks if the type of the media type is a wildcard.
func (m MediaType) IsWildcardType() bool {
	if val, ok := cacheMediaTypes[m]; ok {
		return val.typ == "*"
	}

	return false
}

// Subtype method returns the subtype of the media type.
func (m MediaType) Subtype() string {
	if val, ok := cacheMediaTypes[m]; ok {
		return val.subtype
	}

	return ""
}

// IsWildcardSubtype method checks if the subtype of the media type is a wildcard.
func (m MediaType) IsWildcardSubtype() bool {
	if val, ok := cacheMediaTypes[m]; ok {
		return val.subtype == "*"
	}

	return false
}

// Charset method returns the charset parameter of the media type.
func (m MediaType) Charset() string {
	if val, ok := cacheMediaTypes[m]; ok {
		return val.charset
	}

	return ""
}

// Parameter method returns the value of a parameter of the media type.
func (m MediaType) Parameter(name string) (string, bool) {
	if val, ok := cacheMediaTypes[m]; ok {
		param, exists := val.params[name]
		return param, exists
	}

	return "", false
}
