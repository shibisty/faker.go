package faker

import "fmt"

// localeData holds everything that depends on the locale: names, addresses, phone and zip code formats.
// To add a locale, just add an entry to locales.
type localeData struct {
	firstMale, firstFemale []string
	lastMale, lastFemale   []string
	// genderedLast means the last name depends on gender (ru: Иванов/Иванова).
	// For en there is a single last name, and the extra Bool() call is skipped so that
	// sequences for existing seeds do not change.
	genderedLast bool

	cities  []string
	streets []string
	// street picks the street itself: the order of generator calls differs between locales,
	// and it must be preserved so that old seeds produce the same addresses.
	street func(f *Faker, streets []string) string
	phone  func(f *Faker) string
	zip    func(f *Faker) string
}

var locales = map[string]*localeData{
	"en": {
		firstMale: enFirstNamesMale, firstFemale: enFirstNamesFemale,
		lastMale: enLastNames, lastFemale: enLastNames,
		cities: enCities, streets: enStreets,
		street: func(f *Faker, s []string) string { return fmt.Sprintf("%d %s", f.IntRange(1, 9999), PickOne(f, s)) },
		phone: func(f *Faker) string {
			return fmt.Sprintf("+1-%03d-%03d-%04d", f.IntRange(200, 999), f.IntRange(200, 999), f.IntRange(0, 9999))
		},
		zip: func(f *Faker) string { return fmt.Sprintf("%05d", f.IntRange(10000, 99999)) },
	},
	"ru": {
		firstMale: ruFirstNamesMale, firstFemale: ruFirstNamesFemale,
		lastMale: ruLastNamesMale, lastFemale: ruLastNamesFemale, genderedLast: true,
		cities: ruCities, streets: ruStreets,
		street: func(f *Faker, s []string) string { return fmt.Sprintf("%s, д. %d", PickOne(f, s), f.IntRange(1, 150)) },
		phone: func(f *Faker) string {
			return fmt.Sprintf("+7 (9%02d) %03d-%02d-%02d", f.IntRange(0, 99), f.IntRange(0, 999), f.IntRange(0, 99), f.IntRange(0, 99))
		},
		zip: func(f *Faker) string { return fmt.Sprintf("%06d", f.IntRange(100000, 999999)) },
	},
	"ua": {
		firstMale: uaFirstNamesMale, firstFemale: uaFirstNamesFemale,
		lastMale: uaLastNamesMale, lastFemale: uaLastNamesFemale, genderedLast: true,
		cities: uaCities, streets: uaStreets,
		street: func(f *Faker, s []string) string {
			return fmt.Sprintf("%s, буд. %d", PickOne(f, s), f.IntRange(1, 150))
		},
		phone: func(f *Faker) string {
			return fmt.Sprintf("+380 (%02d) %03d-%02d-%02d", PickOne(f, uaMobileCodes), f.IntRange(0, 999), f.IntRange(0, 99), f.IntRange(0, 99))
		},
		// Ukrainian zip codes have 5 digits.
		zip: func(f *Faker) string { return fmt.Sprintf("%05d", f.IntRange(1000, 99999)) },
	},
}

// uaMobileCodes are the mobile operator codes of Ukraine.
var uaMobileCodes = []int{50, 63, 66, 67, 68, 73, 93, 95, 96, 97, 98, 99}

// localeAliases maps alternative codes: "uk" is the ISO 639-1 code for Ukrainian.
var localeAliases = map[string]string{"uk": "ua"}

// Locales returns the list of supported locale codes.
func Locales() []string { return []string{"en", "ru", "ua"} }

// loc returns the data for the current locale; an unknown locale falls back to en.
func (f *Faker) loc() *localeData {
	code := f.Locale
	if alias, ok := localeAliases[code]; ok {
		code = alias
	}
	if l, ok := locales[code]; ok {
		return l
	}
	return locales["en"]
}
