package faker

// CompanyGen is the faker.Company namespace, the counterpart of faker.company in faker.js.
type CompanyGen struct{ f *Faker }

func (c *CompanyGen) Name() string {
	return PickOne(c.f, companyPrefixes) + " " + PickOne(c.f, companySuffixes)
}

func (c *CompanyGen) JobTitle() string {
	return PickOne(c.f, jobTitles)
}
