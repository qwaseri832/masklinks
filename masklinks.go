package masklinks

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaskRune = '*'

var defaultSchemes = []string{"http://", "https://"}

func DefaultSchemes() []string {
	return slices.Clone(defaultSchemes)
}

func Mask(text string) string {
	return MaskSchemes(text, defaultSchemes...)
}

func MaskSchemes(text string, schemes ...string) string {
	if text == "" || len(schemes) == 0 {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))

	for i := 0; i < len(text); {
		if n := schemeLen(text[i:], schemes); n > 0 {
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

			continue
		}

		_, size := utf8.DecodeRuneInString(text[i:])
		b.WriteString(text[i : i+size])
		i += size
	}

	return b.String()
}

func schemeLen(s string, schemes []string) int {
	longest := 0
	for _, scheme := range schemes {
		n := len(scheme)
		if n == 0 || n <= longest || n > len(s) {
			continue
		}
		if strings.EqualFold(s[:n], scheme) {
			longest = n
		}
	}
	return longest
}
