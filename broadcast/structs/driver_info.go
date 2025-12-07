package structs

// DriverInfo contains information about a driver
type DriverInfo struct {
	FirstName   string
	LastName    string
	ShortName   string
	Category    DriverCategory
	Nationality NationalityEnum
}
