package http

// These constants represent common HTTP header fields.
const (
	// HeaderAccept specifies the media types that are acceptable for the response.
	HeaderAccept = "Accept"
	// HeaderCharset specifies the character sets that are acceptable.
	HeaderCharset = "Accept-Charset"
	// HeaderAcceptEncoding specifies the content codings that are acceptable in the response.
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderAcceptLanguage specifies the natural languages that are preferred in the response.
	HeaderAcceptLanguage = "Accept-Language"
	// HeaderAcceptPatch specifies the patch document formats that are acceptable.
	HeaderAcceptPatch = "Accept-Patch"
	// HeaderAcceptRanges allows the server to indicate its acceptance of range requests for a resource.
	HeaderAcceptRanges = "Accept-Ranges"

	// HeaderAccessControlAllowCredentials indicates whether the response to the request can be exposed when the credentials flag is true.
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// HeaderAccessControlAllowHeaders specifies the headers that are allowed in the actual request.
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderAccessControlAllowMethods specifies the methods that are allowed when accessing the resource.
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// HeaderAccessControlAllowOrigin specifies the origin that is allowed to access the resource.
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderAccessControlExposeHeaders specifies the headers that are exposed to the client.
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderAccessControlMaxAge specifies the maximum amount of time that the results of a preflight request can be cached.
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"
	// HeaderAccessControlRequestHeaders is used when issuing a preflight request to let the server know what HTTP headers will be used in the actual request.
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	// HeaderAccessControlRequestMethod is used when issuing a preflight request to let the server know what HTTP method will be used in the actual request.
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"

	// HeaderAge indicates the age of the response.
	HeaderAge = "Age"
	// HeaderAllow lists the set of methods supported by the resource.
	HeaderAllow = "Allow"
	// HeaderAuthorization contains the credentials to authenticate a user agent with a server.
	HeaderAuthorization = "Authorization"
	// HeaderCacheControl is used to specify directives for caching mechanisms.
	HeaderCacheControl = "Cache-Control"
	// HeaderConnection controls whether the network connection stays open after the current transaction finishes.
	HeaderConnection = "Connection"
	// HeaderContentDisposition is an extension header used in HTTP and MIME email to specify certain parameters related to the disposition of the message content.
	HeaderContentDisposition = "Content-Disposition"
	// HeaderContentEncoding is used to specify the content encodings applied to the entity-body.
	HeaderContentEncoding = "Content-Encoding"
	// HeaderContentLanguage describes the language(s) intended for the audience.
	HeaderContentLanguage = "Content-Language"
	// HeaderContentLength indicates the size of the entity-body in bytes.
	HeaderContentLength = "Content-Length"
	// HeaderContentType indicates the media type of the entity-body.
	HeaderContentType = "Content-Type"
	// HeaderCookie contains stored HTTP cookies previously sent by the server with the Set-Cookie header.
	HeaderCookie = "Cookie"
	// HeaderDate represents the date and time at which the message was originated.
	HeaderDate = "Date"
	// HeaderETag provides the current value of the entity tag for the requested variant.
	HeaderETag = "ETag"
	// HeaderExpect is used to indicate that particular server behaviors are required by the client.
	HeaderExpect = "Expect"
	// HeaderExpires gives the date/time after which the response is considered stale.
	HeaderExpires = "Expires"
	// HeaderFrom is an email address of the human user who controls the requesting user agent.
	HeaderFrom = "From"
	// HeaderHost specifies the domain name of the server and optionally the TCP port number.
	HeaderHost = "Host"

	// HeaderIfMatch is used to make a request method conditional.
	HeaderIfMatch = "If-Match"
	// HeaderIfModifiedSince is used to make a GET or HEAD request method conditional.
	HeaderIfModifiedSince = "If-Modified-Since"
	// HeaderIfNoneMatch is used to make a request method conditional.
	HeaderIfNoneMatch = "If-None-Match"
	// HeaderIfRange is used to make a partial GET request conditional.
	HeaderIfRange = "If-Range"
	// HeaderIfUnmodifiedSince is used to make a request method conditional.
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	// HeaderLastModified indicates the date and time at which the server believes the variant was last modified.
	HeaderLastModified = "Last-Modified"
	// HeaderLink indicates that the response is part of a series of responses.
	HeaderLink = "Link"
	// HeaderLocation is used in redirection, or when a new resource has been created.
	HeaderLocation = "Location"
	// HeaderMaxForwards limits the number of times that the message can be forwarded through proxies or gateways.
	HeaderMaxForwards = "Max-Forwards"
	// HeaderOrigin indicates where a fetch originates from.
	HeaderOrigin = "Origin"
	// HeaderPragma allows backwards compatibility with HTTP/1.0 caches where the Cache-Control header is not yet present.
	HeaderPragma = "Pragma"
	// HeaderProxyAuthenticate must be included as part of a 407 Proxy Authentication Required response.
	HeaderProxyAuthenticate = "Proxy-Authenticate"
	// HeaderProxyAuthorization allows the client to identify itself (or its user) to a proxy which requires authentication.
	HeaderProxyAuthorization = "Proxy-Authorization"
	// HeaderRange is used in an HTTP request to request only part of a document.
	HeaderRange = "Range"
	// HeaderReferer allows the client to specify, for the server's benefit, the address of the document (or element within the document) from which the URI in the request was obtained.
	HeaderReferer = "Referer"
	// HeaderRetryAfter indicates how long the user agent should wait before making a follow-up request.
	HeaderRetryAfter = "Retry-After"
	// HeaderServer contains information about the software used by the origin server to handle the request.
	HeaderServer = "Server"
	// HeaderSetCookie is sent by the server to the user agent with an HTTP response.
	HeaderSetCookie = "Set-Cookie"
	// HeaderSetCookie2 is the updated version of Set-Cookie header.
	HeaderSetCookie2 = "Set-Cookie2"
	// HeaderTE specifies the transfer encodings the user agent is willing to accept.
	HeaderTE = "TE"
	// HeaderTrailer allows the sender to include additional fields at the end of chunked messages.
	HeaderTrailer = "Trailer"
	// HeaderTransferEncoding indicates what (if any) type of transformation has been applied to the message body.
	HeaderTransferEncoding = "TransferEncoding"
	// HeaderUpgrade allows the client to specify what additional communication protocols it supports and would like to use if the server finds it appropriate to switch protocols.
	HeaderUpgrade = "Upgrade"
	// HeaderUserAgent contains information about the user agent originating the request.
	HeaderUserAgent = "UserAgent"
	// HeaderVary determines how to match future request headers to decide whether a cached response can be used rather than requesting a fresh one from the origin server.
	HeaderVary = "Vary"
	// HeaderVia is used by gateways and proxies to indicate the intermediate protocols and recipients between the user agent and the server on requests, and between the origin server and the client on responses.
	HeaderVia = "Via"
	// HeaderWarning is used to carry additional information about the status or transformation of a message which might not be reflected in the message.
	HeaderWarning = "Warning"
	// HeaderWWWAuthenticate must be included in 401 Unauthorized responses.
	HeaderWWWAuthenticate = "WWW-Authenticate"
)
