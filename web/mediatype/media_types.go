package mediatype

type MediaType string

const (
	ApplicationJson           MediaType = "application/json"
	ApplicationXml            MediaType = "application/xml"
	ApplicationFormUrlEncoded MediaType = "application/x-www-form-urlencoded"
	TextPlain                 MediaType = "text/plain"
	TextHtml                  MediaType = "text/html"
	TextXml                   MediaType = "text/xml"
	MultiPartFormData         MediaType = "multipart/form-data"
)
