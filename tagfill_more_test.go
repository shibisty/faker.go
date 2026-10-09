package faker

import (
	"strings"
	"testing"
	"time"
)

// Every generator documented in FillStruct returns a value of the expected type.
func TestGenerateAllTags(t *testing.T) {
	f := New(1)
	cases := []struct {
		tag   string
		check func(any) bool
	}{
		{"first_name", nonEmptyString},
		{"last_name", nonEmptyString},
		{"full_name", func(v any) bool { return strings.Contains(v.(string), " ") }},
		{"username", nonEmptyString},
		{"phone", nonEmptyString},
		{"email", func(v any) bool { return strings.Contains(v.(string), "@") }},
		{"url", func(v any) bool { return strings.HasPrefix(v.(string), "https://") }},
		{"ipv4", func(v any) bool { return strings.Count(v.(string), ".") == 3 }},
		{"password", func(v any) bool { return len(v.(string)) == 12 }},
		{"password:30", func(v any) bool { return len(v.(string)) == 30 }},
		{"city", nonEmptyString},
		{"street", nonEmptyString},
		{"country", nonEmptyString},
		{"zip_code", nonEmptyString},
		{"address", nonEmptyString},
		{"company", nonEmptyString},
		{"job_title", nonEmptyString},
		{"word", nonEmptyString},
		{"words", func(v any) bool { return len(strings.Fields(v.(string))) == 5 }},
		{"words:3", func(v any) bool { return len(strings.Fields(v.(string))) == 3 }},
		{"sentence", func(v any) bool { return strings.HasSuffix(v.(string), ".") }},
		{"sentences", func(v any) bool { return strings.Count(v.(string), ".") == 3 }},
		{"sentences:2", func(v any) bool { return strings.Count(v.(string), ".") == 2 }},
		{"paragraph", nonEmptyString},
		{"paragraphs", func(v any) bool { return strings.Count(v.(string), "\n\n") == 1 }},
		{"paragraphs:3", func(v any) bool { return strings.Count(v.(string), "\n\n") == 2 }},
		{"bool", func(v any) bool { _, ok := v.(bool); return ok }},
		{"uuid", func(v any) bool { return len(v.(string)) == 36 }},
		{"int", func(v any) bool { n := v.(int); return n >= 0 && n <= 100 }},
		{"int:5,7", func(v any) bool { n := v.(int); return n >= 5 && n <= 7 }},
		{"int: 5 , 7 ", func(v any) bool { n := v.(int); return n >= 5 && n <= 7 }},
		{"int:abc,xyz", func(v any) bool { n := v.(int); return n >= 0 && n <= 100 }}, // invalid arguments → default values
		{"float", func(v any) bool { n := v.(float64); return n >= 0 && n < 100 }},
		{"float:1.5,2.5", func(v any) bool { n := v.(float64); return n >= 1.5 && n < 2.5 }},
		{"float:x", func(v any) bool { n := v.(float64); return n >= 0 && n < 100 }},
		{"date_past", func(v any) bool { return v.(time.Time).Before(time.Now()) }},
		{"date_past:10", func(v any) bool { return v.(time.Time).After(time.Now().AddDate(-10, 0, -1)) }},
		{"date_future", func(v any) bool { return v.(time.Time).After(time.Now().Add(-time.Second)) }},
		{"date_future:1", func(v any) bool { return v.(time.Time).Before(time.Now().AddDate(1, 0, 1)) }},
		{"birthday", func(v any) bool { return v.(time.Time).Before(time.Now().AddDate(-18, 0, 1)) }},
		{"birthday:30,40", func(v any) bool { return v.(time.Time).Before(time.Now().AddDate(-30, 0, 1)) }},
	}
	for _, tc := range cases {
		t.Run(tc.tag, func(t *testing.T) {
			v, err := f.generate(tc.tag)
			if err != nil {
				t.Fatalf("generate(%q) error: %v", tc.tag, err)
			}
			if !tc.check(v) {
				t.Fatalf("generate(%q) = %#v failed the check", tc.tag, v)
			}
		})
	}
}

func nonEmptyString(v any) bool { s, ok := v.(string); return ok && s != "" }

// Numeric generators are converted to the various int/float sizes; type errors are clear.
func TestFillStructTypeConversions(t *testing.T) {
	type Numbers struct {
		I8   int8    `fake:"int:1,9"`
		I16  int16   `fake:"int:1,9"`
		I32  int32   `fake:"int:1,9"`
		I64  int64   `fake:"int:1,9"`
		F32  float32 `fake:"float:1,2"`
		F64  float64 `fake:"int:3,3"` // int → float
		IntF int     `fake:"float:4,4.5"`
	}
	f := New(1)
	var n Numbers
	if err := f.FillStruct(&n); err != nil {
		t.Fatal(err)
	}
	if n.I8 < 1 || n.I8 > 9 || n.I16 < 1 || n.I32 < 1 || n.I64 < 1 {
		t.Fatalf("ints not filled: %+v", n)
	}
	if n.F32 < 1 || n.F32 >= 2 || n.F64 != 3 || n.IntF != 4 {
		t.Fatalf("floats not converted: %+v", n)
	}

	errorCases := []struct {
		name string
		dest any
		want string
	}{
		{"string into int", &struct {
			X int `fake:"email"`
		}{}, "cannot convert"},
		{"string into float", &struct {
			X float64 `fake:"email"`
		}{}, "cannot convert"},
		{"int into string", &struct {
			X string `fake:"int"`
		}{}, "expected string"},
		{"string into bool", &struct {
			X bool `fake:"email"`
		}{}, "expected bool"},
		{"unsupported type", &struct {
			X []string `fake:"email"`
		}{}, "not supported"},
		{"bool into int", &struct {
			X int `fake:"bool"`
		}{}, "cannot convert"},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			err := f.FillStruct(tc.dest)
			if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), `field "X"`) {
				t.Fatalf("want error containing %q and the field name, got %v", tc.want, err)
			}
		})
	}
}

func TestFillStructSkips(t *testing.T) {
	type S struct {
		NoTag    string
		Dash     string `fake:"-"`
		private  string `fake:"email"` //nolint:unused // unexported fields must be skipped
		Exported string `fake:"email"`
	}
	f := New(1)
	s := S{NoTag: "keep", Dash: "keep"}
	if err := f.FillStruct(&s); err != nil {
		t.Fatal(err)
	}
	if s.NoTag != "keep" || s.Dash != "keep" || s.private != "" || s.Exported == "" {
		t.Fatalf("unexpected result: %+v", s)
	}
}

func TestFillStructBadInput(t *testing.T) {
	f := New(1)
	var u testUser
	for _, dest := range []any{u, nil, new(int), &[]testUser{}} {
		if err := f.FillStruct(dest); err == nil {
			t.Errorf("FillStruct(%T) should return an error", dest)
		}
	}
	for _, dest := range []any{[]testUser{}, &u, new([]int), new([]*int)} {
		if err := f.FillSlice(dest, 1); err == nil {
			t.Errorf("FillSlice(%T) should return an error", dest)
		}
	}
}

// Regression: an embedded *Base used to panic (NumField on a pointer).
func TestFillStructEmbeddedPointer(t *testing.T) {
	type Base struct {
		Code string `fake:"uuid"`
	}
	type WithPtr struct {
		*Base
		Name string `fake:"first_name"`
	}
	f := New(1)

	var nilBase WithPtr
	if err := f.FillStruct(&nilBase); err != nil {
		t.Fatal(err)
	}
	if nilBase.Base != nil || nilBase.Name == "" {
		t.Fatalf("nil embedded pointer must stay nil, Name filled: %+v", nilBase)
	}

	withBase := WithPtr{Base: &Base{}}
	if err := f.FillStruct(&withBase); err != nil {
		t.Fatal(err)
	}
	if len(withBase.Code) != 36 {
		t.Fatalf("existing embedded pointer should be filled: %+v", withBase.Base)
	}
}

// Regression: an embedded non-struct type used to panic; with a tag it is filled.
func TestFillStructEmbeddedNonStruct(t *testing.T) {
	type ID int
	type Label string
	type S struct {
		ID
		Label `fake:"word"`
		Name  string `fake:"first_name"`
	}
	f := New(1)
	s := S{ID: 7}
	if err := f.FillStruct(&s); err != nil {
		t.Fatal(err)
	}
	if s.ID != 7 || s.Label == "" || s.Name == "" {
		t.Fatalf("unexpected result: %+v", s)
	}
}

// Regression: []*T used to panic.
func TestFillSlicePointers(t *testing.T) {
	f := New(1)
	var users []*testUser
	if err := f.FillSlice(&users, 5); err != nil {
		t.Fatal(err)
	}
	if len(users) != 5 {
		t.Fatalf("got %d users", len(users))
	}
	for _, u := range users {
		if u == nil || u.Name == "" {
			t.Fatalf("pointer element not filled: %+v", u)
		}
	}

	var none []testUser
	if err := f.FillSlice(&none, -3); err != nil || len(none) != 0 {
		t.Fatalf("FillSlice(n<0) = %v, len %d", err, len(none))
	}
}

func TestFillSlicePropagatesError(t *testing.T) {
	type Bad struct {
		X string `fake:"nope"`
	}
	var items []Bad
	if err := New(1).FillSlice(&items, 3); err == nil {
		t.Fatal("FillSlice should return the element error")
	}
}

func TestFillStructEmbeddedError(t *testing.T) {
	type Inner struct {
		X string `fake:"nope"`
	}
	type Outer struct{ Inner }
	type OuterPtr struct{ *Inner }
	f := New(1)
	if err := f.FillStruct(&Outer{}); err == nil {
		t.Error("error in embedded struct should be returned")
	}
	if err := f.FillStruct(&OuterPtr{Inner: &Inner{}}); err == nil {
		t.Error("error in embedded pointer should be returned")
	}
}

// generate does not panic on any tag.
func FuzzGenerate(f *testing.F) {
	for _, seed := range []string{"int:1,2", "float:a,b", "words:-5", "password:-1", "birthday:90,1", ":", "", "paragraphs:1000000000"} {
		f.Add(seed)
	}
	fk := New(1)
	f.Fuzz(func(t *testing.T, tag string) {
		// huge N for words/paragraphs is not a generator bug but just load; cap it
		if strings.Count(tag, "0") > 4 {
			t.Skip()
		}
		_, _ = fk.generate(tag)
	})
}
