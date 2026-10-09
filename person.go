package faker

import (
	"fmt"
	"strings"
)

// PersonGen is the faker.Person namespace, the counterpart of faker.person in faker.js.
type PersonGen struct{ f *Faker }

// FirstName returns a first name of a random gender.
func (p *PersonGen) FirstName() string {
	l := p.f.loc()
	if p.f.Bool() {
		return PickOne(p.f, l.firstMale)
	}
	return PickOne(p.f, l.firstFemale)
}

// LastName returns a last name. For locales with gendered last names (ru, ua), the gender is
// not matched to FirstName() across separate calls; use FullName()
// if you need a gender-consistent first and last name pair.
func (p *PersonGen) LastName() string {
	l := p.f.loc()
	if l.genderedLast {
		if p.f.Bool() {
			return PickOne(p.f, l.lastMale)
		}
		return PickOne(p.f, l.lastFemale)
	}
	return PickOne(p.f, l.lastMale)
}

// FullName returns a gender-consistent first and last name.
func (p *PersonGen) FullName() string {
	l := p.f.loc()
	var first, last string
	if p.f.Bool() {
		first, last = PickOne(p.f, l.firstMale), PickOne(p.f, l.lastMale)
	} else {
		first, last = PickOne(p.f, l.firstFemale), PickOne(p.f, l.lastFemale)
	}
	return first + " " + last
}

// Username generates an ASCII-safe username (transliterates Cyrillic).
func (p *PersonGen) Username() string {
	first := transliterate(p.FirstName())
	return fmt.Sprintf("%s%d", strings.ToLower(first), p.f.IntRange(1, 9999))
}

// Phone returns a phone number in the locale format.
func (p *PersonGen) Phone() string {
	return p.f.loc().phone(p.f)
}
