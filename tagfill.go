package faker

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FillStruct заполняет dest (указатель на структуру) фейковыми данными по
// тегам `fake:"..."`. Поддерживает встроенные (анонимные) структуры —
// например core.BaseModel останется нетронутым (у него нет тега fake, поле
// id обычно генерирует БД, а не сидер).
//
// Поддерживаемые генераторы (имя тега -> что генерируется):
//
//	first_name, last_name, full_name, username, phone
//	email, url, ipv4, password[:длина]
//	city, street, country, zip_code, address
//	company, job_title
//	word, words:N, sentence, sentences:N, paragraph, paragraphs:N
//	bool, uuid
//	int:min,max, float:min,max
//	date_past[:лет], date_future[:лет], birthday[:minAge,maxAge]
//
// Пример:
//
//	type User struct {
//	    Name  string    `fake:"full_name"`
//	    Email string    `fake:"email"`
//	    Age   int       `fake:"int:18,65"`
//	    Bio   string    `fake:"paragraph"`
//	    Born  time.Time `fake:"date_past:60"`
//	}
//	var u User
//	f.FillStruct(&u)
func (f *Faker) FillStruct(dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("faker: FillStruct ожидает указатель на структуру, получено %T", dest)
	}
	return f.fillStructValue(v.Elem())
}

// FillSlice заполняет dest (указатель на срез структур) n фейковыми записями.
func (f *Faker) FillSlice(dest any, n int) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("faker: FillSlice ожидает указатель на срез, получено %T", dest)
	}
	sliceVal := v.Elem()
	elemType := sliceVal.Type().Elem()

	out := reflect.MakeSlice(sliceVal.Type(), 0, n)
	for i := 0; i < n; i++ {
		itemPtr := reflect.New(elemType)
		if err := f.fillStructValue(itemPtr.Elem()); err != nil {
			return err
		}
		out = reflect.Append(out, itemPtr.Elem())
	}
	sliceVal.Set(out)
	return nil
}

func (f *Faker) fillStructValue(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			if err := f.fillStructValue(v.Field(i)); err != nil {
				return err
			}
			continue
		}

		tag := field.Tag.Get("fake")
		if tag == "" || tag == "-" {
			continue
		}
		if !v.Field(i).CanSet() {
			continue
		}

		value, err := f.generate(tag)
		if err != nil {
			return fmt.Errorf("faker: поле %q (fake:%q): %w", field.Name, tag, err)
		}
		if err := setFieldValue(v.Field(i), value); err != nil {
			return fmt.Errorf("faker: поле %q (fake:%q): %w", field.Name, tag, err)
		}
	}
	return nil
}

// generate разбирает тег вида "generator" или "generator:arg1,arg2" и
// возвращает сырое сгенерированное значение (string/int/float64/bool/time.Time).
func (f *Faker) generate(tag string) (any, error) {
	name, argsStr, _ := strings.Cut(tag, ":")
	var args []string
	if argsStr != "" {
		args = strings.Split(argsStr, ",")
	}
	intArg := func(idx, def int) int {
		if idx >= len(args) {
			return def
		}
		n, err := strconv.Atoi(strings.TrimSpace(args[idx]))
		if err != nil {
			return def
		}
		return n
	}
	floatArg := func(idx int, def float64) float64 {
		if idx >= len(args) {
			return def
		}
		n, err := strconv.ParseFloat(strings.TrimSpace(args[idx]), 64)
		if err != nil {
			return def
		}
		return n
	}

	switch name {
	case "first_name":
		return f.Person.FirstName(), nil
	case "last_name":
		return f.Person.LastName(), nil
	case "full_name":
		return f.Person.FullName(), nil
	case "username":
		return f.Person.Username(), nil
	case "phone":
		return f.Person.Phone(), nil
	case "email":
		return f.Internet.Email(), nil
	case "url":
		return f.Internet.URL(), nil
	case "ipv4":
		return f.Internet.IPv4(), nil
	case "password":
		return f.Internet.Password(intArg(0, 12)), nil
	case "city":
		return f.Address.City(), nil
	case "street":
		return f.Address.Street(), nil
	case "country":
		return f.Address.Country(), nil
	case "zip_code":
		return f.Address.ZipCode(), nil
	case "address":
		return f.Address.FullAddress(), nil
	case "company":
		return f.Company.Name(), nil
	case "job_title":
		return f.Company.JobTitle(), nil
	case "word":
		return f.Lorem.Word(), nil
	case "words":
		return f.Lorem.Words(intArg(0, 5)), nil
	case "sentence":
		return f.Lorem.Sentence(), nil
	case "sentences":
		return f.Lorem.Sentences(intArg(0, 3)), nil
	case "paragraph":
		return f.Lorem.Paragraph(), nil
	case "paragraphs":
		return f.Lorem.Paragraphs(intArg(0, 2), "\n\n"), nil
	case "bool":
		return f.Bool(), nil
	case "uuid":
		return f.UUID(), nil
	case "int":
		return f.IntRange(intArg(0, 0), intArg(1, 100)), nil
	case "float":
		return f.FloatRange(floatArg(0, 0), floatArg(1, 100)), nil
	case "date_past":
		return f.Date.Past(intArg(0, 2)), nil
	case "date_future":
		return f.Date.Future(intArg(0, 2)), nil
	case "birthday":
		return f.Date.Birthday(intArg(0, 18), intArg(1, 65)), nil
	default:
		return nil, fmt.Errorf("неизвестный генератор %q", name)
	}
}

func setFieldValue(field reflect.Value, value any) error {
	if !field.CanSet() {
		return nil
	}
	rv := reflect.ValueOf(value)
	if rv.Type().AssignableTo(field.Type()) {
		field.Set(rv)
		return nil
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch n := value.(type) {
		case int:
			field.SetInt(int64(n))
		case int64:
			field.SetInt(n)
		case float64:
			field.SetInt(int64(n))
		default:
			return fmt.Errorf("не удалось привести %T к %s", value, field.Kind())
		}
	case reflect.Float32, reflect.Float64:
		switch n := value.(type) {
		case float64:
			field.SetFloat(n)
		case int:
			field.SetFloat(float64(n))
		default:
			return fmt.Errorf("не удалось привести %T к %s", value, field.Kind())
		}
	case reflect.String:
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("ожидалась строка, получено %T", value)
		}
		field.SetString(s)
	case reflect.Bool:
		b, ok := value.(bool)
		if !ok {
			return fmt.Errorf("ожидался bool, получено %T", value)
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("тип поля %s не поддерживается FillStruct", field.Type())
	}
	return nil
}
