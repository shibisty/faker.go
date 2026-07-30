package faker

import (
	"fmt"
	"strings"
)

// PersonGen — namespace faker.Person, аналог faker.person из faker.js.
type PersonGen struct{ f *Faker }

func (p *PersonGen) FirstName() string {
	male := p.f.Bool()
	switch p.f.Locale {
	case "ru":
		if male {
			return PickOne(p.f, ruFirstNamesMale)
		}
		return PickOne(p.f, ruFirstNamesFemale)
	default:
		if male {
			return PickOne(p.f, enFirstNamesMale)
		}
		return PickOne(p.f, enFirstNamesFemale)
	}
}

// LastName возвращает фамилию. Для "ru" род фамилии не согласован с родом
// FirstName() при отдельных вызовах — используйте FullName(), если нужна
// согласованная по роду пара имя+фамилия.
func (p *PersonGen) LastName() string {
	if p.f.Locale == "ru" {
		if p.f.Bool() {
			return PickOne(p.f, ruLastNamesMale)
		}
		return PickOne(p.f, ruLastNamesFemale)
	}
	return PickOne(p.f, enLastNames)
}

// FullName возвращает согласованные по роду имя и фамилию.
func (p *PersonGen) FullName() string {
	male := p.f.Bool()
	var first, last string
	switch p.f.Locale {
	case "ru":
		if male {
			first, last = PickOne(p.f, ruFirstNamesMale), PickOne(p.f, ruLastNamesMale)
		} else {
			first, last = PickOne(p.f, ruFirstNamesFemale), PickOne(p.f, ruLastNamesFemale)
		}
	default:
		if male {
			first = PickOne(p.f, enFirstNamesMale)
		} else {
			first = PickOne(p.f, enFirstNamesFemale)
		}
		last = PickOne(p.f, enLastNames)
	}
	return first + " " + last
}

// Username генерирует ASCII-безопасный username (транслитерирует кириллицу).
func (p *PersonGen) Username() string {
	first := transliterate(p.FirstName())
	return fmt.Sprintf("%s%d", strings.ToLower(first), p.f.IntRange(1, 9999))
}

func (p *PersonGen) Phone() string {
	if p.f.Locale == "ru" {
		return fmt.Sprintf("+7 (9%02d) %03d-%02d-%02d",
			p.f.IntRange(0, 99), p.f.IntRange(0, 999), p.f.IntRange(0, 99), p.f.IntRange(0, 99))
	}
	return fmt.Sprintf("+1-%03d-%03d-%04d",
		p.f.IntRange(200, 999), p.f.IntRange(200, 999), p.f.IntRange(0, 9999))
}
