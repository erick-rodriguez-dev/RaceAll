package structs

import "time"

// RealtimeUpdate contains real-time session update information
type RealtimeUpdate struct {
	EventIndex           int
	SessionIndex         int
	Phase                SessionPhase
	SessionTime          time.Duration
	RemainingTime        time.Duration
	TimeOfDay            time.Duration
	RainLevel            float32
	Clouds               float32
	Wetness              float32
	BestSessionLap       *LapInfo
	BestLapCarIndex      uint16
	BestLapDriverIndex   uint16
	FocusedCarIndex      int
	ActiveCameraSet      string
	ActiveCamera         string
	IsReplayPlaying      bool
	ReplaySessionTime    float32
	ReplayRemainingTime  float32
	SessionRemainingTime time.Duration
	SessionEndTime       time.Duration
	SessionType          RaceSessionType
	AmbientTemp          uint8
	TrackTemp            uint8
	CurrentHudPage       string
}
