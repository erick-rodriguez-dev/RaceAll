package structs

// ConnectionState represents the state of a broadcast connection
type ConnectionState struct {
	ConnectionId      int
	ConnectionSuccess bool
	IsReadonly        bool
	Error             string
}
