package faker

import "time"

// DateGen — namespace faker.Date, аналог faker.date из faker.js.
type DateGen struct{ f *Faker }

// Between возвращает случайный момент времени в [start, end).
func (d *DateGen) Between(start, end time.Time) time.Time {
	if !end.After(start) {
		return start
	}
	delta := end.Sub(start)
	return start.Add(time.Duration(d.f.float01() * float64(delta)))
}

// Past возвращает случайную дату за последние maxYearsAgo лет.
func (d *DateGen) Past(maxYearsAgo int) time.Time {
	now := time.Now()
	return d.Between(now.AddDate(-maxYearsAgo, 0, 0), now)
}

// Future возвращает случайную дату в пределах ближайших maxYearsAhead лет.
func (d *DateGen) Future(maxYearsAhead int) time.Time {
	now := time.Now()
	return d.Between(now, now.AddDate(maxYearsAhead, 0, 0))
}

// Birthday возвращает дату рождения для возраста в диапазоне [minAge, maxAge].
func (d *DateGen) Birthday(minAge, maxAge int) time.Time {
	now := time.Now()
	age := d.f.IntRange(minAge, maxAge)
	return now.AddDate(-age, -d.f.IntRange(0, 11), -d.f.IntRange(0, 27))
}
