package masklinks

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaskRune = '*'

var DefaultSchemes = []string{"http://", "https://"}

func Mask(text string) string {
	return MaskSchemes(text, DefaultSchemes...)
}

func MaskSchemes(text string, schemes ...string) string {
	if text == "" || len(schemes) == 0 {
		return text
	}

	lower := strings.ToLower(text)
	lowered := make([]string, len(schemes))
	for i, s := range schemes {
		lowered[i] = strings.ToLower(s)
	}

	var b strings.Builder
	b.Grow(len(text))

	for i := 0; i < len(text); {
		n := schemeLen(lower[i:], lowered)
		if n == 0 {
			b.WriteByte(text[i])
			i++
			continue
		}

		b.WriteString(text[i : i+n])
		i += n

		for i < len(text) {
			r, size := utf8.DecodeRuneInString(text[i:])
			if unicode.IsSpace(r) {
				break
			}
			b.WriteRune(MaskRune)
			i += size
		}
	}

	return b.String()
}

func schemeLen(s string, lowered []string) int {
	longest := 0
	for _, scheme := range lowered {
		if scheme == "" || len(scheme) <= longest {
			continue
		}
		if strings.HasPrefix(s, scheme) {
			longest = len(scheme)
		}
	}
	return longest
}
