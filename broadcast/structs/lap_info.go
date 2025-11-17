package structs

import "fmt"

// LapInfo contains information about a lap
type LapInfo struct {
	// LaptimeMS - Lap time in milliseconds. Use "GetLapTimeMS()"
	// as this value sometime is reported wrong by the server.
	LaptimeMS *int

	// Splits - Sector time in milliseconds.
	Splits []int

	// CarIndex - Internal server identifier assigned to the car.
	CarIndex uint16

	// DriverIndex - Driver index
	DriverIndex uint16

	// IsInvalid - Lap is not valid due corner cut, out of track, etc.
	IsInvalid bool

	// IsValidForBest - Improving its own the best time.
	IsValidForBest bool

	// Type - Lap type
	Type LapType
}

// NewLapInfo creates a new LapInfo instance
func NewLapInfo() *LapInfo {
	return &LapInfo{
		Splits: make([]int, 0),
	}
}

// GetLapTimeMS returns the lap time in milliseconds.
// Preferred to use this method instead of LapTimeMS property.
// The LapTimeMS Property sometimes doesn't return the correct laptime.
func (l *LapInfo) GetLapTimeMS() int {
	lapTimeMs := 0
	for _, split := range l.Splits {
		lapTimeMs += split
	}
	return lapTimeMs
}

// String returns a string representation of the lap info
func (l *LapInfo) String() string {
	lapTime := 0
	if l.LaptimeMS != nil {
		lapTime = *l.LaptimeMS
	}
	return fmt.Sprintf("%5d|%v", lapTime, l.Splits)
}
