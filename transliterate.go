package faker

import "strings"

var cyrillicToLatin = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
	'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	// Ukrainian letters
	'і': "i", 'ї': "yi", 'є': "ye", 'ґ': "g",
	// the apostrophe in Ukrainian names (Мар'яна) is not needed in email/username
	'\'': "", 'ʼ': "", '’': "",
}

// transliterate converts Cyrillic to Latin (for emails/usernames built from
// Russian and Ukrainian names) and lowercases the result; Latin characters are left as is.
func transliterate(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if lat, ok := cyrillicToLatin[r]; ok {
			b.WriteString(lat)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
