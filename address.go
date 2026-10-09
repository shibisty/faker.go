package faker

import "fmt"

// AddressGen is the faker.Address namespace, the counterpart of faker.location in faker.js.
type AddressGen struct{ f *Faker }

// City returns a city for the locale.
func (a *AddressGen) City() string {
	return PickOne(a.f, a.f.loc().cities)
}

// Street returns a street with a house number in the locale format.
func (a *AddressGen) Street() string {
	l := a.f.loc()
	return l.street(a.f, l.streets)
}

// Country returns a country name (in English for all locales).
func (a *AddressGen) Country() string {
	return PickOne(a.f, countries)
}

// ZipCode returns a zip code in the locale format.
func (a *AddressGen) ZipCode() string {
	return a.f.loc().zip(a.f)
}

// FullAddress joins the street, city, country and zip code into one string.
func (a *AddressGen) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s", a.Street(), a.City(), a.Country(), a.ZipCode())
}
