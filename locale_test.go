package faker

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяІЇЄҐіїєґ"

func TestLocales(t *testing.T) {
	got := Locales()
	for _, code := range got {
		if _, ok := locales[code]; !ok {
			t.Errorf("Locales() lists %q, but there is no data for it", code)
		}
	}
	if len(got) != len(locales) {
		t.Errorf("Locales() = %v, data has %d locales", got, len(locales))
	}
}

func TestLocaleFallbackAndAlias(t *testing.T) {
	f := New(1)

	f.Locale = "xx"
	if f.loc() != locales["en"] {
		t.Error("unknown locale should fall back to en")
	}
	f.Locale = ""
	if f.loc() != locales["en"] {
		t.Error("empty locale should fall back to en")
	}
	f.Locale = "uk"
	if f.loc() != locales["ua"] {
		t.Error(`"uk" should be an alias for "ua"`)
	}
}

// For each locale: data is non-empty, names come from its own set, phone and zip code formats match.
func TestPersonAndAddressPerLocale(t *testing.T) {
	cases := []struct {
		locale   string
		phone    *regexp.Regexp
		zip      *regexp.Regexp
		street   *regexp.Regexp
		cyrillic bool
	}{
		{"en", regexp.MustCompile(`^\+1-[2-9]\d{2}-[2-9]\d{2}-\d{4}$`), regexp.MustCompile(`^\d{5}$`), regexp.MustCompile(`^\d{1,4} .+$`), false},
		{"ru", regexp.MustCompile(`^\+7 \(9\d{2}\) \d{3}-\d{2}-\d{2}$`), regexp.MustCompile(`^\d{6}$`), regexp.MustCompile(`^.+, д\. \d{1,3}$`), true},
		{"ua", regexp.MustCompile(`^\+380 \((50|63|66|67|68|73|93|95|96|97|98|99)\) \d{3}-\d{2}-\d{2}$`), regexp.MustCompile(`^\d{5}$`), regexp.MustCompile(`^.+, буд\. \d{1,3}$`), true},
	}

	for _, tc := range cases {
		t.Run(tc.locale, func(t *testing.T) {
			f := New(3)
			f.Locale = tc.locale
			l := locales[tc.locale]

			for _, list := range [][]string{l.firstMale, l.firstFemale, l.lastMale, l.lastFemale, l.cities, l.streets} {
				if len(list) == 0 {
					t.Fatal("locale has an empty data list")
				}
			}

			for i := 0; i < 100; i++ {
				first := f.Person.FirstName()
				if !contains(l.firstMale, first) && !contains(l.firstFemale, first) {
					t.Fatalf("FirstName %q is not from the %s data", first, tc.locale)
				}
				last := f.Person.LastName()
				if !contains(l.lastMale, last) && !contains(l.lastFemale, last) {
					t.Fatalf("LastName %q is not from the %s data", last, tc.locale)
				}

				full := f.Person.FullName()
				parts := strings.SplitN(full, " ", 2)
				if len(parts) != 2 {
					t.Fatalf("FullName %q should be 'first last'", full)
				}
				maleFirst := contains(l.firstMale, parts[0])
				if maleFirst && !contains(l.lastMale, parts[1]) || !maleFirst && !contains(l.lastFemale, parts[1]) {
					t.Fatalf("FullName %q: first and last name genders do not match", full)
				}
				if got := strings.ContainsAny(full, cyrillic); got != tc.cyrillic {
					t.Fatalf("FullName %q: cyrillic = %v, want %v", full, got, tc.cyrillic)
				}

				if p := f.Person.Phone(); !tc.phone.MatchString(p) {
					t.Fatalf("Phone %q does not match %s", p, tc.phone)
				}
				if z := f.Address.ZipCode(); !tc.zip.MatchString(z) {
					t.Fatalf("ZipCode %q does not match %s", z, tc.zip)
				}
				if s := f.Address.Street(); !tc.street.MatchString(s) {
					t.Fatalf("Street %q does not match %s", s, tc.street)
				}
				if c := f.Address.City(); !contains(l.cities, c) {
					t.Fatalf("City %q is not from the %s data", c, tc.locale)
				}
				if c := f.Address.Country(); !contains(countries, c) {
					t.Fatalf("Country %q is unknown", c)
				}
				if a := f.Address.FullAddress(); strings.Count(a, ", ") < 2 {
					t.Fatalf("FullAddress %q has too few parts", a)
				}
			}
		})
	}
}

// Email and username are always ASCII, for any locale, including Ukrainian letters and the apostrophe.
func TestEmailAndUsernameAreASCII(t *testing.T) {
	emailRe := regexp.MustCompile(`^[a-z]+\.[a-z]+\d{1,3}@[a-z0-9.-]+$`)
	userRe := regexp.MustCompile(`^[a-z]+\d{1,4}$`)

	for _, locale := range Locales() {
		f := New(11)
		f.Locale = locale
		for i := 0; i < 300; i++ {
			if e := f.Internet.Email(); !emailRe.MatchString(e) {
				t.Fatalf("locale %s: email %q is not plain ASCII local@domain", locale, e)
			}
			if u := f.Person.Username(); !userRe.MatchString(u) {
				t.Fatalf("locale %s: username %q", locale, u)
			}
			if u := f.Internet.Username(); !userRe.MatchString(u) {
				t.Fatalf("locale %s: Internet.Username %q", locale, u)
			}
		}
	}
}

// All names of all locales transliterate to Latin with nothing left over.
func TestTransliterateCoversAllNames(t *testing.T) {
	for code, l := range locales {
		for _, list := range [][]string{l.firstMale, l.firstFemale, l.lastMale, l.lastFemale} {
			for _, name := range list {
				out := transliterate(name)
				for _, r := range out {
					if r > unicode.MaxASCII {
						t.Errorf("locale %s: transliterate(%q) = %q keeps non-ASCII %q", code, name, out, r)
						break
					}
				}
			}
		}
	}
}

func TestTransliterate(t *testing.T) {
	cases := map[string]string{
		"Иван":        "ivan",
		"Щукин":       "schukin",
		"Объём":       "obem",
		"Їжак":        "yizhak",
		"Євген":       "yevgen",
		"Ґанок":       "ganok",
		"Ірина":       "irina",
		"Мар'яна":     "maryana",
		"Мар’яна":     "maryana",
		"John Smith":  "john smith",
		"":            "",
		"Тест-123_ok": "test-123_ok",
	}
	for in, want := range cases {
		if got := transliterate(in); got != want {
			t.Errorf("transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestInternet(t *testing.T) {
	f := New(5)
	for i := 0; i < 100; i++ {
		u := f.Internet.URL()
		parsed, err := url.Parse(u)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			t.Fatalf("URL %q is not a valid https URL", u)
		}

		ip := f.Internet.IPv4()
		if parsed := net.ParseIP(ip); parsed == nil || parsed.To4() == nil {
			t.Fatalf("IPv4 %q is not a valid IPv4 address", ip)
		}

		if ua := f.Internet.UserAgent(); !contains(userAgents, ua) {
			t.Fatalf("UserAgent %q is not from the data", ua)
		}
	}
}

func TestPassword(t *testing.T) {
	f := New(5)
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	cases := []struct {
		args []int
		want int
	}{
		{nil, 12},
		{[]int{20}, 20},
		{[]int{1}, 1},
		{[]int{0}, 12},
		{[]int{-5}, 12},
	}
	for _, tc := range cases {
		p := f.Internet.Password(tc.args...)
		if len(p) != tc.want {
			t.Errorf("Password(%v) length = %d, want %d", tc.args, len(p), tc.want)
		}
		for _, r := range p {
			if !strings.ContainsRune(chars, r) {
				t.Errorf("Password(%v) = %q contains unexpected %q", tc.args, p, r)
			}
		}
	}
}

func TestCompany(t *testing.T) {
	f := New(5)
	for i := 0; i < 50; i++ {
		name := f.Company.Name()
		parts := strings.SplitN(name, " ", 2)
		if len(parts) != 2 || !contains(companyPrefixes, parts[0]) || !contains(companySuffixes, parts[1]) {
			t.Fatalf("Company.Name %q is not prefix+suffix", name)
		}
		if j := f.Company.JobTitle(); !contains(jobTitles, j) {
			t.Fatalf("JobTitle %q is not from the data", j)
		}
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
