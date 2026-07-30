package faker

import (
	"strings"
	"testing"
	"time"
)

func TestDeterministicSeed(t *testing.T) {
	a := New(42)
	b := New(42)

	for i := 0; i < 20; i++ {
		na, nb := a.Person.FullName(), b.Person.FullName()
		if na != nb {
			t.Fatalf("the same seed should produce the same sequence: %q != %q", na, nb)
		}
	}
}

func TestDifferentSeedsDiffer(t *testing.T) {
	a := New(1)
	b := New(2)

	same := 0
	const n = 30
	for i := 0; i < n; i++ {
		if a.Person.FullName() == b.Person.FullName() {
			same++
		}
	}
	if same == n {
		t.Fatal("different seeds should not produce an identical sequence of 30 consecutive values")
	}
}

func TestIntRange(t *testing.T) {
	f := New(1)
	for i := 0; i < 200; i++ {
		v := f.IntRange(10, 20)
		if v < 10 || v > 20 {
			t.Fatalf("IntRange(10,20) returned a value outside the expected range: %d", v)
		}
	}
}

func TestFloatRange(t *testing.T) {
	f := New(1)
	for i := 0; i < 200; i++ {
		v := f.FloatRange(1.5, 2.5)
		if v < 1.5 || v >= 2.5 {
			t.Fatalf("FloatRange(1.5,2.5) returned a value outside the expected range: %v", v)
		}
	}
}

func TestEmailFormat(t *testing.T) {
	f := New(1)
	for i := 0; i < 50; i++ {
		email := f.Internet.Email()
		if !strings.Contains(email, "@") {
			t.Fatalf("email is missing '@': %q", email)
		}
		if strings.ContainsAny(email, "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя") {
			t.Fatalf("email should be transliterated into Latin characters: %q", email)
		}
	}
}

func TestEmailRuLocale(t *testing.T) {
	f := New(1)
	f.Locale = "ru"
	for i := 0; i < 50; i++ {
		email := f.Internet.Email()
		if strings.ContainsAny(email, "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя") {
			t.Fatalf("email should still use Latin characters when locale=ru: %q", email)
		}
	}
}

func TestPersonRuLocale(t *testing.T) {
	f := New(1)
	f.Locale = "ru"
	name := f.Person.FullName()
	if !strings.ContainsAny(name, "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя") {
		t.Fatalf("FullName should use Cyrillic characters when locale=ru: %q", name)
	}
}

func TestUniqueEmailNoDuplicates(t *testing.T) {
	f := New(7)
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		email := f.Unique().Email()
		if seen[email] {
			t.Fatalf("Unique().Email() returned a duplicate: %q", email)
		}
		seen[email] = true
	}
}

func TestUUIDFormat(t *testing.T) {
	f := New(1)
	id := f.UUID()
	parts := strings.Split(id, "-")
	if len(parts) != 5 {
		t.Fatalf("UUID should consist of 5 hyphen-separated parts, got %d: %q", len(parts), id)
	}
	wantLens := []int{8, 4, 4, 4, 12}
	for i, p := range parts {
		if len(p) != wantLens[i] {
			t.Fatalf("UUID part %d has length %d, expected %d (%q)", i, len(p), wantLens[i], id)
		}
	}
	if id[14] != '4' {
		t.Fatalf("UUID should be version 4, got %q", id)
	}
}

type testUser struct {
	Name  string    `fake:"full_name"`
	Email string    `fake:"email"`
	Age   int       `fake:"int:18,65"`
	Bio   string    `fake:"paragraph"`
	Admin bool      `fake:"bool"`
	Born  time.Time `fake:"date_past:60"`
}

func TestFillStruct(t *testing.T) {
	f := New(1)
	var u testUser
	if err := f.FillStruct(&u); err != nil {
		t.Fatalf("FillStruct returned an error: %v", err)
	}
	if u.Name == "" || u.Email == "" || u.Bio == "" {
		t.Fatalf("string fields should not be empty: %+v", u)
	}
	if u.Age < 18 || u.Age > 65 {
		t.Fatalf("Age is outside the range specified by tag int:18,65: %d", u.Age)
	}
	if !strings.Contains(u.Email, "@") {
		t.Fatalf("Email does not look like a valid email address: %q", u.Email)
	}
	if u.Born.After(time.Now()) {
		t.Fatalf("date_past should generate a date in the past: %v", u.Born)
	}
}

func TestFillSlice(t *testing.T) {
	f := New(1)
	var users []testUser
	if err := f.FillSlice(&users, 10); err != nil {
		t.Fatalf("FillSlice returned an error: %v", err)
	}
	if len(users) != 10 {
		t.Fatalf("expected 10 records, got %d", len(users))
	}
	for _, u := range users {
		if u.Name == "" {
			t.Fatal("FillSlice produced a record with an empty Name")
		}
	}
}

func TestFillStructWithEmbedded(t *testing.T) {
	type Base struct {
		ID int64 `db:"id"` // no fake tag — should remain unchanged
	}
	type WithBase struct {
		Base
		Name string `fake:"first_name"`
	}

	f := New(1)
	var w WithBase
	w.ID = 999

	if err := f.FillStruct(&w); err != nil {
		t.Fatalf("FillStruct returned an error: %v", err)
	}
	if w.ID != 999 {
		t.Fatalf("field without a fake tag should not be modified, ID = %d", w.ID)
	}
	if w.Name == "" {
		t.Fatal("Name should be populated")
	}
}

func TestUnknownGeneratorReturnsError(t *testing.T) {
	type Bad struct {
		X string `fake:"totally_unknown_generator"`
	}

	f := New(1)
	var b Bad

	if err := f.FillStruct(&b); err == nil {
		t.Fatal("expected an error for an unknown generator")
	}
}

func TestPickOne(t *testing.T) {
	f := New(1)
	items := []string{"a", "b", "c"}

	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		seen[PickOne(f, items)] = true
	}

	if len(seen) == 0 {
		t.Fatal("PickOne returned no values")
	}

	for k := range seen {
		found := false
		for _, it := range items {
			if it == k {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("PickOne returned a value not present in the source slice: %q", k)
		}
	}
}
