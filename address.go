package faker

import "fmt"

// AddressGen — namespace faker.Address, аналог faker.location из faker.js.
type AddressGen struct{ f *Faker }

func (a *AddressGen) City() string {
	if a.f.Locale == "ru" {
		return PickOne(a.f, ruCities)
	}
	
	if a.f.Locale == "ua" {
		return PickOne(a.f, uaCities)
	}
	
	return PickOne(a.f, enCities)
}

func (a *AddressGen) Street() string {
	if a.f.Locale == "ru" {
		return fmt.Sprintf("%s, д. %d", PickOne(a.f, ruStreets), a.f.IntRange(1, 150))
	}
	
	if a.f.Locale == "ua" {
		return fmt.Sprintf("%s, д. %d", PickOne(a.f, uaStreets), a.f.IntRange(1, 150))
	}
	
	return fmt.Sprintf("%d %s", a.f.IntRange(1, 9999), PickOne(a.f, enStreets))
}

func (a *AddressGen) Country() string {
	return PickOne(a.f, countries)
}

func (a *AddressGen) ZipCode() string {
	if a.f.Locale == "ru" || a.f.Locale == "ua" {
		return fmt.Sprintf("%06d", a.f.IntRange(100000, 999999))
	}
	return fmt.Sprintf("%05d", a.f.IntRange(10000, 99999))
}

// FullAddress собирает улицу, город, страну и индекс в одну строку.
func (a *AddressGen) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s", a.Street(), a.City(), a.Country(), a.ZipCode())
}
