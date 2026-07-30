package faker

import "strings"

var cyrillicToLatin = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
	'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

// transliterate переводит кириллицу в латиницу (для email/username из
// русских имён) и приводит к нижнему регистру; латинские символы не трогает.
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
