# faker.go

[![Patreon](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/cw/shibisty)

> Lightweight fake data generator for Go with zero dependencies.
>
> Generate realistic fake data for tests, database seeders, demos and mock APIs.

[![Go Version](https://img.shields.io/badge/go-1.18+-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Unlike most faker libraries, **faker.go** provides both a familiar API (`faker.js` style) and automatic struct population using tags, making it ideal for database seeding and testing. :contentReference[oaicite:0]{index=0}

---

## Features

- ✅ Zero external dependencies
- ✅ Familiar API inspired by `faker.js`
- ✅ Automatic struct population
- ✅ Deterministic random generator
- ✅ Multiple locales
- ✅ UUID v4 generator
- ✅ Unique value generation
- ✅ Generic helpers
- ✅ Perfect for database seeders
- ✅ Small and fast

---

# Installation

```bash
go get github.com/shibisty/faker.go
```

---

# Quick Start

```go
package main

import (
    "fmt"

    faker "github.com/shibisty/faker.go"
)

func main() {
    f := faker.New()

    fmt.Println(f.Person.FullName())
    fmt.Println(f.Internet.Email())
    fmt.Println(f.Address.City())
    fmt.Println(f.Company.Name())
}
```

Example output:

```
John Smith
john.smith@gmail.com
New York
Acme Corporation
```

---

# Deterministic Seeds

Using the same seed always produces the same sequence.

```go
f := faker.New(42)

fmt.Println(f.Person.FullName())
fmt.Println(f.Person.FullName())
```

Perfect for reproducible tests. :contentReference[oaicite:1]{index=1}

---

# Locales

Default locale:

```go
f := faker.New()
```

Russian:

```go
f := faker.New()
f.Locale = "ru"
```

Ukrainian:

```go
f := faker.New()
f.Locale = "ua"
```

Supported locales currently:

- English (`en`)
- Russian (`ru`)
- Ukrainian (`ua`)

---

# Person

```go
f.Person.FirstName()

f.Person.LastName()

f.Person.FullName()

f.Person.Username()

f.Person.Phone()
```

---

# Internet

```go
f.Internet.Email()

f.Internet.Domain()

f.Internet.URL()
```

---

# Address

```go
f.Address.City()

f.Address.Street()

f.Address.ZipCode()

f.Address.Country()
```

---

# Company

```go
f.Company.Name()

f.Company.Buzzword()
```

---

# Lorem Ipsum

```go
f.Lorem.Word()

f.Lorem.Sentence()

f.Lorem.Paragraph()
```

Lorem always generates classic Latin placeholder text regardless of locale. :contentReference[oaicite:2]{index=2}

---

# Date & Time

```go
f.Date.Past()

f.Date.Future()

f.Date.Between()
```

---

# Utility Functions

## Boolean

```go
f.Bool()
```

---

## Integer Range

```go
f.IntRange(10, 50)
```

---

## Float Range

```go
f.FloatRange(1.5, 10.0)
```

---

## UUID

```go
f.UUID()
```

Returns a valid UUID v4. :contentReference[oaicite:3]{index=3}

---

## Pick Random Item

```go
statuses := []string{
    "pending",
    "paid",
    "cancelled",
}

status := faker.PickOne(f, statuses)
```

Works with any slice type thanks to generics. :contentReference[oaicite:4]{index=4}

---

# Unique Values

Useful for tables with UNIQUE constraints.

```go
f := faker.New()

email := f.Unique().Email()

username := f.Unique().Username()

phone := f.Unique().Phone()
```

The generator keeps track of previously generated values and avoids duplicates whenever possible. :contentReference[oaicite:5]{index=5}

---

# Struct Population

Instead of manually assigning every field, simply annotate your struct.

```go
type User struct {
    Name  string `fake:"full_name"`
    Email string `fake:"email"`
    Age   int    `fake:"int:18,65"`
}
```

Populate it automatically:

```go
var user User

err := f.FillStruct(&user)
```

---

## Slice Population

```go
var users []User

f.FillSlice(&users, 100)
```

Creates 100 fully populated users.

---

# Example Seeder

```go
users := make([]User, 0)

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

A single `Faker` instance is **not safe** for concurrent use.

For parallel workloads create one instance per goroutine:

```go
go func() {
    f := faker.New()

    ...
}()
```

or synchronize access yourself. :contentReference[oaicite:6]{index=6}

---

# DLL Support

The library can also be compiled as a native shared library:

```bash
go build -buildmode=c-shared
```

This allows using the generator from:

- C
- C++
- C#
- Rust
- Python
- Delphi
- Go (via DLL loading)
- Any language capable of calling C ABI functions

---

# Why another faker library?

- No dependencies
- Small footprint
- Deterministic output
- Generic helpers
- Struct tag support
- Database seeding focused
- DLL compatible

---

# License

MIT

[![Patreon](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/cw/shibisty)

If this project helps you, consider supporting its development on Patreon ❤️
