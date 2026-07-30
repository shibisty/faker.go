// package main generates realistic fake data for database seeders and tests —
// a lightweight alternative to faker.js/gofakeit with zero external dependencies.
//
// There are two ways to use it:
//
//  1. Direct API calls, similar to faker.js (faker.person.firstName() -> f.Person.FirstName()):
//
//     f := faker.New()
//     f.Person.FullName()
//     f.Internet.Email()
//     f.Address.City()
//
//  2. Tag-based struct population (especially useful for seeders — your model
//     is already annotated with `db:"..."` tags for the ORM, and now you can
//     add `fake:"..."` tags as well):
//
//     type User struct {
//         Name  string `db:"name"  fake:"full_name"`
//         Email string `db:"email" fake:"email"`
//         Age   int    `db:"age"   fake:"int:18,65"`
//     }
//
//     var u User
//     f.FillStruct(&u)
package faker

import (
	"math/rand"
	"sync"
	"time"
)

// Faker generates fake data.
//
// A single Faker instance is NOT safe for concurrent use without external
// synchronization. Create one Faker per goroutine when seeding in parallel,
// or protect access with your own mutex.
type Faker struct {
	rnd *rand.Rand
	mu  sync.Mutex

	// Locale affects Person and Address generators: "en" (default) or "ru".
	// Lorem always generates classic pseudo-Latin lorem ipsum text, since
	// that's the industry standard placeholder and is not localized.
	Locale string

	Person   *PersonGen
	Internet *InternetGen
	Address  *AddressGen
	Company  *CompanyGen
	Lorem    *LoremGen
	Date     *DateGen

	uniqueMu   sync.Mutex
	uniqueSeen map[string]map[string]struct{}
}

// New creates a new Faker instance.
//
//	faker.New()      // non-deterministic (seeded with the current time)
//	faker.New(42)    // deterministic: the same seed always produces the same
//	                 // sequence, which is useful for reproducible tests and seed data
func New(seed ...int64) *Faker {
	s := time.Now().UnixNano()
	if len(seed) > 0 {
		s = seed[0]
	}

	f := &Faker{
		rnd:        rand.New(rand.NewSource(s)),
		Locale:     "en",
		uniqueSeen: map[string]map[string]struct{}{},
	}

	f.Person = &PersonGen{f: f}
	f.Internet = &InternetGen{f: f}
	f.Address = &AddressGen{f: f}
	f.Company = &CompanyGen{f: f}
	f.Lorem = &LoremGen{f: f}
	f.Date = &DateGen{f: f}

	return f
}

func (f *Faker) intn(n int) int {
	if n <= 0 {
		return 0
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	return f.rnd.Intn(n)
}

func (f *Faker) float01() float64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.rnd.Float64()
}

func (f *Faker) randomBytes(n int) []byte {
	b := make([]byte, n)

	f.mu.Lock()
	f.rnd.Read(b) //nolint:errcheck // math/rand.Rand.Read never returns an error.
	f.mu.Unlock()

	return b
}

// PickOne returns a random element from a non-empty slice of any type.
// Useful for custom value lists (for example, order statuses) directly
// inside your seeder code.
func PickOne[T any](f *Faker, items []T) T {
	return items[f.intn(len(items))]
}

// Bool returns a random boolean value.
func (f *Faker) Bool() bool { return f.intn(2) == 0 }

// IntRange returns a random integer in the inclusive range [min, max].
func (f *Faker) IntRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + f.intn(max-min+1)
}

// FloatRange returns a random floating-point number in the range [min, max).
func (f *Faker) FloatRange(min, max float64) float64 {
	if max <= min {
		return min
	}
	return min + f.float01()*(max-min)
}

// UUID generates a random UUID v4.
func (f *Faker) UUID() string {
	b := f.randomBytes(16)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return formatUUID(b)
}

// unique wraps a generator function and guarantees that values generated
// under the given name (for example, "email") will not repeat within this
// Faker instance. After a limited number of attempts, it returns the last
// generated value anyway to avoid an infinite loop when the source dataset
// is too small.
func (f *Faker) unique(name string, gen func() string) string {
	const maxAttempts = 50

	f.uniqueMu.Lock()
	defer f.uniqueMu.Unlock()

	seen, ok := f.uniqueSeen[name]
	if !ok {
		seen = map[string]struct{}{}
		f.uniqueSeen[name] = seen
	}

	var value string
	for i := 0; i < maxAttempts; i++ {
		value = gen()
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			return value
		}
	}

	seen[value] = struct{}{}
	return value
}

// UniqueFaker wraps Faker and guarantees that generators for unique fields
// (email, username, phone) never produce duplicates within the same Faker
// instance. This is especially useful when seeding tables with UNIQUE
// constraints.
type UniqueFaker struct{ f *Faker }

// Unique returns a UniqueFaker wrapper for this Faker.
func (f *Faker) Unique() *UniqueFaker { return &UniqueFaker{f: f} }

func (u *UniqueFaker) Email() string    { return u.f.unique("email", u.f.Internet.Email) }
func (u *UniqueFaker) Username() string { return u.f.unique("username", u.f.Person.Username) }
func (u *UniqueFaker) Phone() string    { return u.f.unique("phone", u.f.Person.Phone) }

func main() {}
