package broadcast

import (
	"math"
	"time"
)

// DriverInfo
type DriverInfo struct {
	FirstName   string
	LastName    string
	ShortName   string
	Category    DriverCategory
	Nationality NationalityEnum
}

func (d *DriverInfo) GetFullName() string {
	return d.FirstName + " " + d.LastName
}

// LapInfo 
type LapInfo struct {
	// LaptimeMS is the time of the lap in milliseconds
	// could be nil if no lap time is recorded
	LaptimeMS *int32

	// Splits contains the times of each sector in milliseconds
	// could contain nil values for sectors not completed
	Splits [3]*int32

	// CarIndex is the identifier of the car on the server
	CarIndex uint16

	// DriverIndex is the index of the driver within the car
	DriverIndex uint16

	// IsInvalid indicates if the lap is invalid (cut, off track, etc.)
	IsInvalid bool

	// IsValidForBest indicates if the lap is valid for best time
	IsValidForBest bool

	// Type indicates the type of lap (Outlap, Regular, Inlap)
	Type LapType
}

// GetLapTimeMS returns the lap time calculated from the splits
// This method is preferable over LaptimeMS since sometimes the server
// does not report the lap time correctly
func (l *LapInfo) GetLapTimeMS() int32 {
	var totalTime int32 = 0
	for _, split := range l.Splits {
		if split != nil {
			totalTime += *split
		}
	}
	return totalTime
}

// HasValidLapTime returns true if the lap has a valid time
func (l *LapInfo) HasValidLapTime() bool {
	return l.LaptimeMS != nil && *l.LaptimeMS != math.MaxInt32
}

// CarInfo contains information about a car
type CarInfo struct {
	// CarIndex is the internal identifier of the car on the server
	CarIndex uint16

	// CarModelType is the model of the car (byte)
	CarModelType byte

	// TeamName is the name of the team
	TeamName string

	// RaceNumber is the race number
	RaceNumber int32

	// CupCategory: Overall/Pro = 0, ProAm = 1, Am = 2, Silver = 3, National = 4
	CupCategory byte

	// CurrentDriverIndex is the index of the current driver in the Drivers array
	CurrentDriverIndex byte

	// Drivers is the list of drivers for the car
	Drivers []DriverInfo

	// Nationality is the nationality of the team
	Nationality NationalityEnum
}

// GetCurrentDriver returns the current driver
func (c *CarInfo) GetCurrentDriver() *DriverInfo {
	if int(c.CurrentDriverIndex) < len(c.Drivers) {
		return &c.Drivers[c.CurrentDriverIndex]
	}
	return nil
}

// GetCurrentDriverName returns the name of the current driver
func (c *CarInfo) GetCurrentDriverName() string {
	driver := c.GetCurrentDriver()
	if driver != nil {
		return driver.GetFullName()
	}
	return "Unknown"
}

// RealtimeUpdate contains real-time update information for the session
type RealtimeUpdate struct {
	// EventIndex is the event index
	EventIndex uint16

	// SessionIndex is the session index
	SessionIndex uint16

	// SessionType is the type of session (Practice, Qualifying, Race, etc.)
	SessionType RaceSessionType

	// Phase is the current phase of the session
	Phase SessionPhase

	// SessionTime is the elapsed time of the session
	SessionTime time.Duration

	// SessionEndTime is the end time of the session
	SessionEndTime time.Duration

	// FocusedCarIndex is the index of the focused car
	FocusedCarIndex int32

	// ActiveCameraSet is the active camera set
	ActiveCameraSet string

	// ActiveCamera is the active camera
	ActiveCamera string

	// CurrentHudPage is the current HUD page
	CurrentHudPage string

	// IsReplayPlaying indicates if a replay is currently playing	
	IsReplayPlaying bool

	// ReplaySessionTime is the session time in the replay
	ReplaySessionTime float32

	// ReplayRemainingTime is the remaining time in the replay
	ReplayRemainingTime float32

	// TimeOfDay is the time of day
	TimeOfDay time.Duration

	// AmbientTemp is the ambient temperature in °C
	AmbientTemp byte

	// TrackTemp is the track temperature in °C
	TrackTemp byte

	// Clouds is the cloud level (0.0 - 1.0)
	Clouds float32

	// RainLevel is the rain level (0.0 - 1.0)
	RainLevel float32

	// Wetness is the wetness level of the track (0.0 - 1.0)
	Wetness float32

	// BestSessionLap is the best lap of the session
	BestSessionLap LapInfo
}

// RealtimeCarUpdate contains real-time update information for a car
type RealtimeCarUpdate struct {
	// CarIndex is the index of the car
	CarIndex uint16

	// DriverIndex is the index of the driver (may change in endurance races)
	DriverIndex uint16

	// DriverCount is the number of drivers in the car
	DriverCount byte

	// Gear is the current gear (-1 = R, 0 = N, 1+ = gears)
	Gear int8

	// WorldPosX is the X position in the world
	WorldPosX float32

	// WorldPosY is the Y position in the world
	WorldPosY float32

	// Heading is the orientation of the car in radians
	Heading float32

	// CarLocation is the location of the car (Track, Pitlane, etc.)
	CarLocation CarLocationEnum

	// Kmh is the speed in km/h
	Kmh uint16

	// Position is the official position in P/Q/R (base 1)
	Position uint16

	// CupPosition is the position in the Cup category (base 1)
	CupPosition uint16

	// TrackPosition is the position on the track (base 1)
	TrackPosition uint16

	// SplinePosition is the position on the track spline (0.0 - 1.0)
	SplinePosition float32

	// Laps is the number of laps completed
	Laps uint16

	// Delta is the real-time delta to the best lap of the session (ms)
	Delta int32

	// BestSessionLap is the best lap of the session for this car
	BestSessionLap LapInfo

	// LastLap is the last lap completed
	LastLap LapInfo

	// CurrentLap is the current lap
	CurrentLap LapInfo
}

// TrackData contains information about the track
type TrackData struct {
	// TrackName is the name of the track
	TrackName string

	// TrackId is the identifier of the track
	TrackId int32

	// TrackMeters is the length of the track in meters
	TrackMeters int32

	// CameraSets is a map of camera sets and their cameras
	CameraSets map[string][]string

	// HUDPages is a list of available HUD pages
	HUDPages []string
}

// BroadcastingEvent represents a broadcasting event
type BroadcastingEvent struct {
	// Type is the type of event
	Type BroadcastingEventType

	// Msg is the message of the event
	Msg string

	// TimeMs is the time of the event in milliseconds
	TimeMs int32

	// CarId is the identifier of the car involved
	CarId int32

	// CarData is the information of the car (may be nil)
	CarData *CarInfo
}

// ConnectionState represents the state of the connection
type ConnectionState struct {
	// ConnectionId is the identifier of the connection
	ConnectionId int32

	// Success indicates if the connection was successful
	Success bool

	// IsReadonly indicates if the connection is read-only
	IsReadonly bool

	// ErrorMsg is the error message (if any)
	ErrorMsg string
}
