package export

import (
	"strings"
	"unicode"
)

func SanitizeName(name string) string {
	b := strings.Builder{}

	for _, r := range name {
		if unicode.IsSpace(r) {
			b.WriteRune('_')
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}
