package faker

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLorem(t *testing.T) {
	f := New(9)

	if w := f.Lorem.Word(); !contains(loremWords, w) {
		t.Fatalf("Word %q is not a lorem word", w)
	}

	for _, n := range []int{-1, 0} {
		if got := f.Lorem.Words(n); got != "" {
			t.Errorf("Words(%d) = %q, want empty", n, got)
		}
		if got := f.Lorem.Sentences(n); got != "" {
			t.Errorf("Sentences(%d) = %q, want empty", n, got)
		}
		if got := f.Lorem.Paragraphs(n, "\n"); got != "" {
			t.Errorf("Paragraphs(%d) = %q, want empty", n, got)
		}
	}

	if got := len(strings.Fields(f.Lorem.Words(7))); got != 7 {
		t.Errorf("Words(7) has %d words", got)
	}

	for i := 0; i < 50; i++ {
		s := f.Lorem.Sentence()
		words := strings.Fields(strings.TrimSuffix(s, "."))
		if !strings.HasSuffix(s, ".") || len(words) < 5 || len(words) > 12 {
			t.Fatalf("Sentence %q: want 5..12 words ending with a dot", s)
		}
		if first := []rune(s)[0]; strings.ToUpper(string(first)) != string(first) {
			t.Fatalf("Sentence %q should start with a capital letter", s)
		}
	}

	if got := strings.Count(f.Lorem.Sentences(4), "."); got != 4 {
		t.Errorf("Sentences(4) has %d dots", got)
	}

	p := f.Lorem.Paragraph()
	if n := strings.Count(p, "."); n < 3 || n > 6 {
		t.Errorf("Paragraph has %d sentences, want 3..6", n)
	}

	if got := len(strings.Split(f.Lorem.Paragraphs(3, "\n\n"), "\n\n")); got != 3 {
		t.Errorf("Paragraphs(3) has %d paragraphs", got)
	}
}

func TestDate(t *testing.T) {
	f := New(9)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(48 * time.Hour)
	for i := 0; i < 100; i++ {
		d := f.Date.Between(start, end)
		if d.Before(start) || !d.Before(end) {
			t.Fatalf("Between returned %v outside [%v, %v)", d, start, end)
		}
	}
	if got := f.Date.Between(end, start); !got.Equal(end) {
		t.Errorf("Between(end, start) = %v, want start argument %v", got, end)
	}
	if got := f.Date.Between(start, start); !got.Equal(start) {
		t.Errorf("Between(x, x) = %v, want %v", got, start)
	}

	now := time.Now()
	for i := 0; i < 100; i++ {
		if d := f.Date.Past(2); d.After(time.Now()) || d.Before(now.AddDate(-2, 0, -1)) {
			t.Fatalf("Past(2) = %v", d)
		}
		if d := f.Date.Future(2); d.Before(now) || d.After(time.Now().AddDate(2, 0, 1)) {
			t.Fatalf("Future(2) = %v", d)
		}
		b := f.Date.Birthday(18, 30)
		age := now.Sub(b).Hours() / 24 / 365.25
		if age < 18 || age > 31.1 {
			t.Fatalf("Birthday(18, 30) gives age %.1f", age)
		}
	}
}

func TestRangesEdgeCases(t *testing.T) {
	f := New(1)
	if got := f.IntRange(5, 5); got != 5 {
		t.Errorf("IntRange(5,5) = %d", got)
	}
	if got := f.IntRange(10, 1); got != 10 {
		t.Errorf("IntRange(10,1) = %d, want min", got)
	}
	if got := f.FloatRange(2, 2); got != 2 {
		t.Errorf("FloatRange(2,2) = %v", got)
	}
	if got := f.FloatRange(3, 1); got != 3 {
		t.Errorf("FloatRange(3,1) = %v, want min", got)
	}
	if got := f.intn(0); got != 0 {
		t.Errorf("intn(0) = %d", got)
	}

	seen := map[int]bool{}
	for i := 0; i < 500; i++ {
		seen[f.IntRange(1, 3)] = true
	}
	if len(seen) != 3 {
		t.Errorf("IntRange(1,3) should produce 1, 2 and 3, got %v", seen)
	}

	trues := 0
	for i := 0; i < 1000; i++ {
		if f.Bool() {
			trues++
		}
	}
	if trues < 400 || trues > 600 {
		t.Errorf("Bool() true %d times out of 1000", trues)
	}
}

func TestUUIDUniqueAndVariant(t *testing.T) {
	f := New(1)
	seen := map[string]bool{}
	for i := 0; i < 1000; i++ {
		id := f.UUID()
		if seen[id] {
			t.Fatalf("duplicate UUID %s", id)
		}
		seen[id] = true
		if v := id[19]; v != '8' && v != '9' && v != 'a' && v != 'b' {
			t.Fatalf("UUID %s: variant nibble %q, want 8/9/a/b (RFC 4122)", id, v)
		}
	}
}

func TestPickOneEmptyPanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil || !strings.Contains(r.(string), "empty slice") {
			t.Fatalf("PickOne on empty slice should panic with a clear message, got %v", r)
		}
	}()
	PickOne(New(1), []int{})
}

func TestUniqueAllKinds(t *testing.T) {
	f := New(3)
	u := f.Unique()
	for _, gen := range []func() string{u.Email, u.Username, u.Phone} {
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			v := gen()
			if seen[v] {
				t.Fatalf("duplicate unique value %q", v)
			}
			seen[v] = true
		}
	}
}

// When the source set is exhausted, unique does not loop forever but returns a duplicate.
func TestUniqueExhaustedDoesNotHang(t *testing.T) {
	f := New(1)
	calls := 0
	gen := func() string { calls++; return "same" }

	if got := f.unique("k", gen); got != "same" || calls != 1 {
		t.Fatalf("first call: got %q after %d calls", got, calls)
	}
	calls = 0
	if got := f.unique("k", gen); got != "same" || calls != 50 {
		t.Fatalf("exhausted: got %q after %d calls, want 50 attempts", got, calls)
	}
	// A different name has its own set of values.
	if got := f.unique("other", gen); got != "same" {
		t.Fatalf("unique sets should be per name, got %q", got)
	}
}

// One Faker from several goroutines: no races (run with -race) and no panics.
func TestConcurrentUse(t *testing.T) {
	f := New(1)
	f.Locale = "ua"
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				_ = f.Person.FullName()
				_ = f.Internet.Email()
				_ = f.Unique().Email()
				_ = f.UUID()
				_ = f.Lorem.Sentence()
				var u testUser
				if err := f.FillStruct(&u); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()
}

// Determinism for all generators and all locales, not just FullName.
func TestDeterministicAllGenerators(t *testing.T) {
	sample := func(f *Faker) []string {
		return []string{
			f.Person.FirstName(), f.Person.LastName(), f.Person.FullName(), f.Person.Username(), f.Person.Phone(),
			f.Internet.Email(), f.Internet.URL(), f.Internet.IPv4(), f.Internet.Password(), f.Internet.UserAgent(),
			f.Address.FullAddress(), f.Company.Name(), f.Company.JobTitle(), f.Lorem.Paragraph(), f.UUID(),
		}
	}
	for _, locale := range Locales() {
		a, b := New(77), New(77)
		a.Locale, b.Locale = locale, locale
		for i := 0; i < 10; i++ {
			sa, sb := sample(a), sample(b)
			for j := range sa {
				if sa[j] != sb[j] {
					t.Fatalf("locale %s: same seed gave %q and %q", locale, sa[j], sb[j])
				}
			}
		}
	}
}

func TestRefDateMakesDatesReproducible(t *testing.T) {
	ref := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	type row struct {
		Born time.Time `fake:"birthday:18,65"`
		Seen time.Time `fake:"date_past:2"`
		Next time.Time `fake:"date_future:1"`
	}
	gen := func() row {
		f := New(5)
		f.RefDate = ref
		var r row
		if err := f.FillStruct(&r); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond)
		return r
	}
	a, b := gen(), gen()
	if a != b {
		t.Fatalf("same seed and RefDate must give the same dates: %+v vs %+v", a, b)
	}
	if !a.Seen.Before(ref) || a.Seen.Before(ref.AddDate(-2, 0, 0)) || !a.Next.After(ref) || a.Next.After(ref.AddDate(1, 0, 0)) {
		t.Fatalf("dates must count from RefDate: %+v", a)
	}
	if age := ref.Sub(a.Born).Hours() / 24 / 365.25; age < 18 || age > 66.1 {
		t.Fatalf("birthday age %.1f", age)
	}
}
