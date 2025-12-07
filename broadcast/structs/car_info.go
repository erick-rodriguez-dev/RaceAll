package structs

// CarInfo contains information about a car
type CarInfo struct {
	CarIndex           uint16
	CarModelType       uint8
	TeamName           string
	RaceNumber         int
	CupCategory        uint8
	CurrentDriverIndex int
	Drivers            []DriverInfo
	Nationality        NationalityEnum
}

// NewCarInfo creates a new CarInfo instance
func NewCarInfo(carIndex uint16) *CarInfo {
	return &CarInfo{
		CarIndex: carIndex,
		Drivers:  make([]DriverInfo, 0),
	}
}

// AddDriver adds a driver to the car
func (c *CarInfo) AddDriver(driverInfo DriverInfo) {
	c.Drivers = append(c.Drivers, driverInfo)
}

// GetCurrentDriverName returns the current driver's last name
func (c *CarInfo) GetCurrentDriverName() string {
	if c.CurrentDriverIndex < len(c.Drivers) {
		return c.Drivers[c.CurrentDriverIndex].LastName
	}
	return "nobody(?)"
}
