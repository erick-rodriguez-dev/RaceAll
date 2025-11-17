package structs

// RealtimeCarUpdate contains real-time update information for a car
type RealtimeCarUpdate struct {
	CarIndex       int
	DriverIndex    int
	Gear           int
	Heading        float32
	WorldPosX      float32
	WorldPosY      float32
	CarLocation    CarLocationEnum
	Kmh            int
	Position       int
	TrackPosition  int
	SplinePosition float32
	Delta          int
	BestSessionLap *LapInfo
	LastLap        *LapInfo
	CurrentLap     *LapInfo
	Laps           int
	CupPosition    uint16
	DriverCount    uint8
}
