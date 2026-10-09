package faker

import (
	"fmt"
	"strings"
)

// InternetGen is the faker.Internet namespace, the counterpart of faker.internet in faker.js.
type InternetGen struct{ f *Faker }

func (i *InternetGen) Email() string {
	local := transliterate(i.f.Person.FirstName()) + "." + transliterate(i.f.Person.LastName())
	local = strings.ReplaceAll(local, " ", "")
	domain := PickOne(i.f, emailDomains)
	return fmt.Sprintf("%s%d@%s", strings.ToLower(local), i.f.IntRange(1, 999), domain)
}

// Username is an alias from faker.Internet.Username() to faker.Person.Username(),
// as in faker.js, where internet.userName() and person.firstName() overlap.
func (i *InternetGen) Username() string { return i.f.Person.Username() }

func (i *InternetGen) URL() string {
	return fmt.Sprintf("https://%s.%s/%s",
		strings.ToLower(PickOne(i.f, loremWords)),
		PickOne(i.f, tlds),
		strings.ToLower(PickOne(i.f, loremWords)))
}

func (i *InternetGen) IPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		i.f.IntRange(1, 254), i.f.IntRange(0, 255), i.f.IntRange(0, 255), i.f.IntRange(1, 254))
}

// Password generates a random password. The default length is 12 characters;
// optionally pass the desired length: Password(20).
func (i *InternetGen) Password(length ...int) string {
	n := 12
	if len(length) > 0 && length[0] > 0 {
		n = length[0]
	}
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, n)
	for idx := range b {
		b[idx] = chars[i.f.intn(len(chars))]
	}
	return string(b)
}

func (i *InternetGen) UserAgent() string {
	return PickOne(i.f, userAgents)
}
