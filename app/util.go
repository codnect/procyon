package app

import "strings"

func sliceToDelimitedString(values []string) string {
	var builder strings.Builder
	for index, value := range values {
		builder.WriteByte('"')
		builder.WriteString(value)
		builder.WriteByte('"')
		if index != len(values)-1 {
			builder.WriteByte(',')
		}
	}
	return builder.String()
}
