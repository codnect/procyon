package http

import (
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

// Cookie represents an HTTP cookie.
type Cookie struct {
	// Name is the name of the cookie.
	Name string
	// Value is the value of the cookie.
	Value string

	// Path is the path of the cookie.
	Path string
	// Domain is the domain of the cookie.
	Domain string
	// Expires is the expiration time of the cookie.
	Expires time.Time

	// MaxAge is the maximum age of the cookie.
	MaxAge int
	// Secure is the secure flag of the cookie.
	Secure bool
	// HttpOnly is the http only flag of the cookie.
	HttpOnly bool
	// SameSite is the same site flag of the cookie.
	SameSite SameSite
}

// SameSite represents the SameSite attribute of the cookie.
type SameSite int

const (
	// SameSiteDefaultMode represents the default mode of the SameSite attribute.
	SameSiteDefaultMode SameSite = iota + 1
	// SameSiteLaxMode represents the lax mode of the SameSite attribute.
	SameSiteLaxMode
	// SameSiteStrictMode represents the strict mode of the SameSite attribute.
	SameSiteStrictMode
	// SameSiteNoneMode represents the none mode of the SameSite attribute.
	SameSiteNoneMode
)

func parseCookies(header http.Header) []*Cookie {
	lines := header["Cookie"]

	if len(lines) == 0 {
		return []*Cookie{}
	}

	cookies := make([]*Cookie, 0, len(lines)+strings.Count(lines[0], ";"))
	for _, line := range lines {
		line = textproto.TrimString(line)

		var part string
		for len(line) > 0 {
			part, line, _ = strings.Cut(line, ";")
			part = textproto.TrimString(part)
			if part == "" {
				continue
			}
			name, val, _ := strings.Cut(part, "=")
			name = textproto.TrimString(name)
			if !isCookieNameValid(name) {
				continue
			}

			val, ok := parseCookieValue(val, true)
			if !ok {
				continue
			}
			cookies = append(cookies, &Cookie{Name: name, Value: val})
		}
	}
	return cookies
}

func isCookieNameValid(name string) bool {
	if name == "" {
		return false
	}
	return strings.IndexFunc(name, isNotToken) < 0
}

func isNotToken(r rune) bool {
	return !isTokenRune(r)
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

var isTokenTable = [127]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
}

func isTokenRune(r rune) bool {
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
}
