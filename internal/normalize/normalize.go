package normalize

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func CategoryName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	t := norm.NFD.String(name)

	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	return norm.NFC.String(b.String())
}
