package faker

// CompanyGen — namespace faker.Company, аналог faker.company из faker.js.
type CompanyGen struct{ f *Faker }

func (c *CompanyGen) Name() string {
	return PickOne(c.f, companyPrefixes) + " " + PickOne(c.f, companySuffixes)
}

func (c *CompanyGen) JobTitle() string {
	return PickOne(c.f, jobTitles)
}
