# faker.go

[![Patreon](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/cw/shibisty)

> Lightweight fake data generator for Go with zero dependencies.
>
> Generate realistic fake data for tests, database seeders, demos and mock APIs.

[![CI](https://github.com/shibisty/faker.go/actions/workflows/ci.workflow.yml/badge.svg)](https://github.com/shibisty/faker.go/actions/workflows/ci.workflow.yml)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Unlike most faker libraries, **faker.go** provides both a familiar API (`faker.js` style) and automatic struct population using tags, making it ideal for database seeding and testing.

---

## Features

- ✅ Zero external dependencies
- ✅ Familiar API inspired by `faker.js`
- ✅ Automatic struct population with `fake:"..."` tags
- ✅ Deterministic output for a given seed
- ✅ Locales: English, Russian, Ukrainian
- ✅ UUID v4 generator
- ✅ Unique value generation
- ✅ Generic helpers
- ✅ Perfect for database seeders
- ✅ Small and fast

---

# Installation

With [gtr](https://github.com/shibisty/gtr):

```bash
gtr add github:shibisty/faker.go
```

Then `import "faker"`.

Requires Go 1.22 or newer.

---

# Quick Start

```go
package main

import (
    "fmt"

    "faker"
)

func main() {
    f := faker.New(42)

    fmt.Println(f.Person.FullName())
    fmt.Println(f.Internet.Email())
    fmt.Println(f.Address.City())
    fmt.Println(f.Company.Name())
}
```

Example output:

```
Donna Rodriguez
mark.sanchez594@test.io
Dallas
Fusion Technologies
```

---

# Deterministic Seeds

Using the same seed always produces the same sequence.

```go
f := faker.New(42)

fmt.Println(f.Person.FullName())
fmt.Println(f.Person.FullName())
```

Perfect for reproducible tests. Determinism holds when one `Faker` is used from a single goroutine.
`Date` generators count from the current time; set `f.RefDate` to a fixed moment to make dates reproducible too:

```go
f := faker.New(42)
f.RefDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
```

---

# Locales

```go
f := faker.New()   // English ("en") by default
f.Locale = "ru"    // Russian
f.Locale = "ua"    // Ukrainian ("uk" works too)
```

`faker.Locales()` returns the list of supported codes. Unknown codes fall back to English.

Locale affects names, phone numbers, cities, streets and zip codes:

| Locale | Name | Phone | Street | Zip |
|---|---|---|---|---|
| `en` | James Brown | +1-263-812-7508 | 6347 River Rd | 5 digits |
| `ru` | Александр Зайцев | +7 (963) 412-08-72 | пр. Мира, д. 53 | 6 digits |
| `ua` | Олександр Лисенко | +380 (67) 412-08-72 | просп. Незалежності, буд. 53 | 5 digits |

Emails and usernames are always ASCII: Cyrillic names are transliterated (`Євген` → `yevgen`).
Lorem always generates classic Latin placeholder text regardless of locale.

---

# Person

```go
f.Person.FirstName()
f.Person.LastName()
f.Person.FullName()   // first and last name of the same gender
f.Person.Username()
f.Person.Phone()
```

---

# Internet

```go
f.Internet.Email()
f.Internet.Username()
f.Internet.URL()
f.Internet.IPv4()
f.Internet.Password()     // 12 characters
f.Internet.Password(20)   // custom length
f.Internet.UserAgent()
```

---

# Address

```go
f.Address.City()
f.Address.Street()
f.Address.ZipCode()
f.Address.Country()
f.Address.FullAddress()
```

---

# Company

```go
f.Company.Name()
f.Company.JobTitle()
```

---

# Lorem Ipsum

```go
f.Lorem.Word()
f.Lorem.Words(5)
f.Lorem.Sentence()
f.Lorem.Sentences(3)
f.Lorem.Paragraph()
f.Lorem.Paragraphs(2, "\n\n")
```

---

# Date & Time

```go
f.Date.Past(2)                 // within the last 2 years
f.Date.Future(1)               // within the next year
f.Date.Between(start, end)     // in [start, end)
f.Date.Birthday(18, 65)        // for an age between 18 and 65
```

---

# Utility Functions

```go
f.Bool()
f.IntRange(10, 50)       // inclusive
f.FloatRange(1.5, 10.0)  // [min, max)
f.UUID()                 // valid UUID v4

statuses := []string{"pending", "paid", "cancelled"}
status := faker.PickOne(f, statuses) // any slice type; panics on an empty slice
```

---

# Unique Values

Useful for tables with UNIQUE constraints.

```go
email := f.Unique().Email()
username := f.Unique().Username()
phone := f.Unique().Phone()
```

The generator keeps track of previously generated values and avoids duplicates. If the data set is exhausted, it returns a repeated value after 50 attempts instead of looping forever.

---

# Struct Population

Annotate your struct instead of assigning every field:

```go
type User struct {
    ID    int64     `db:"id"`                 // no tag — left untouched
    Name  string    `db:"name"  fake:"full_name"`
    Email string    `db:"email" fake:"email"`
    Age   int       `db:"age"   fake:"int:18,65"`
    Born  time.Time `db:"born"  fake:"birthday:18,65"`
}

var user User
err := f.FillStruct(&user)

var users []User   // or []*User
err = f.FillSlice(&users, 100)
```

Embedded structs are filled recursively. An embedded `*Base` is filled only if it is not nil.

| Tag | Result |
|---|---|
| `first_name`, `last_name`, `full_name`, `username`, `phone` | Person |
| `email`, `url`, `ipv4`, `password[:length]` | Internet |
| `city`, `street`, `country`, `zip_code`, `address` | Address |
| `company`, `job_title` | Company |
| `word`, `words[:n]`, `sentence`, `sentences[:n]`, `paragraph`, `paragraphs[:n]` | Lorem |
| `bool`, `uuid` | Utility |
| `int[:min,max]`, `float[:min,max]` | Numbers (default 0..100) |
| `date_past[:years]`, `date_future[:years]`, `birthday[:minAge,maxAge]` | `time.Time` |
| `-` | Skip the field |

Numeric values are converted to the field type (`int8`…`int64`, `float32`, `float64`).
An unknown tag or an incompatible field type returns an error that names the field.

---

# Example Seeder

```go
users := make([]User, 0, 1000)

for i := 0; i < 1000; i++ {
    users = append(users, User{
        Name:  f.Person.FullName(),
        Email: f.Unique().Email(),
        Age:   f.IntRange(18, 65),
    })
}
```

---

# Thread Safety

All methods are safe to call from several goroutines: access to the random source is synchronized.
The sequence is deterministic only within a single goroutine, and `Locale` must not be changed while other goroutines use the same `Faker`.
For reproducible parallel seeding, create one `Faker` per goroutine with its own seed.

---

# Development

In a checkout, run `gtr install` once (it generates `go.mod`), then:

```bash
gtr run test -- -race -cover
```

With the `go` shim (`gtr self shims`), plain go commands work too:

```bash
go vet ./...
go test -run '^$' -fuzz=FuzzGenerate -fuzztime=30s .
```

---

# Why another faker library?

- No dependencies
- Small footprint
- Deterministic output
- Generic helpers
- Struct tag support
- Database seeding focused

---

# License

MIT

[![Patreon](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/cw/shibisty)

If this project helps you, consider supporting its development on Patreon ❤️
