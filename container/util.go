package container

import (
	"strings"
	"unicode"
)

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}
