package handlers

import (
	"regexp"
	"strings"
)

var (
	reNonAlnum  = regexp.MustCompile(`[^a-z0-9]+`)
	reMultiDash = regexp.MustCompile(`-{2,}`)
)

// slugify converts a display name to a URL-safe slug.
// Handles Cyrillic (Russian/Turkmen) and Latin diacritics.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = transliterateCyrillic(s)
	s = reNonAlnum.ReplaceAllString(s, "-")
	s = reMultiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "item"
	}
	return s
}

var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "j", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	// Turkmen
	'ä': "ae", 'ň': "ng", 'ö': "o", 'ü': "u", 'ý': "y", 'ş': "sh", 'ž': "zh",
	// common Latin diacritics
	'á': "a", 'à': "a", 'â': "a", 'ã': "a", 'å': "a",
	'é': "e", 'è': "e", 'ê': "e", 'ë': "e",
	'í': "i", 'ì': "i", 'î': "i", 'ï': "i",
	'ó': "o", 'ò': "o", 'ô': "o", 'õ': "o", 'ø': "o",
	'ú': "u", 'ù': "u", 'û': "u",
	'ñ': "n", 'ç': "c", 'ß': "ss",
}

func transliterateCyrillic(s string) string {
	var b strings.Builder
	for _, r := range s {
		if mapped, ok := cyrillicMap[r]; ok {
			b.WriteString(mapped)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
