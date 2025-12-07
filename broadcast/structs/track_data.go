package structs

// TrackData contains information about the track
type TrackData struct {
	TrackName   string
	TrackId     int
	TrackMeters float32
	CameraSets  map[string][]string
	HUDPages    []string
}
