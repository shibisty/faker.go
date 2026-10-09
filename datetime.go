package faker

import "time"

// DateGen is the faker.Date namespace, the counterpart of faker.date in faker.js.
type DateGen struct{ f *Faker }

// Between returns a random point in time in [start, end).
func (d *DateGen) Between(start, end time.Time) time.Time {
	if !end.After(start) {
		return start
	}
	delta := end.Sub(start)
	return start.Add(time.Duration(d.f.float01() * float64(delta)))
}

// Past returns a random date within the last maxYearsAgo years.
func (d *DateGen) Past(maxYearsAgo int) time.Time {
	now := d.f.now()
	return d.Between(now.AddDate(-maxYearsAgo, 0, 0), now)
}

// Future returns a random date within the next maxYearsAhead years.
func (d *DateGen) Future(maxYearsAhead int) time.Time {
	now := d.f.now()
	return d.Between(now, now.AddDate(maxYearsAhead, 0, 0))
}

// Birthday returns a birth date for an age in the range [minAge, maxAge].
func (d *DateGen) Birthday(minAge, maxAge int) time.Time {
	now := d.f.now()
	age := d.f.IntRange(minAge, maxAge)
	return now.AddDate(-age, -d.f.IntRange(0, 11), -d.f.IntRange(0, 27))
}
