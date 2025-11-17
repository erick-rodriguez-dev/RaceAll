package structs

// BroadcastingEvent represents a broadcasting event
type BroadcastingEvent struct {
	Type    BroadcastingCarEventType
	Msg     string
	TimeMs  int
	CarId   int
	CarData *CarInfo
}
