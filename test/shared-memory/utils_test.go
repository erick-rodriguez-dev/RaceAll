package sharedmemory_test

import (
	"testing"

	"RaceAll/internal/sharedmemory"
)

func TestUTF16ToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint16
		expected string
	}{
		{
			name:     "Simple ASCII string",
			input:    []uint16{'H', 'e', 'l', 'l', 'o', 0},
			expected: "Hello",
		},
		{
			name:     "Empty string",
			input:    []uint16{0},
			expected: "",
		},
		{
			name:     "String with numbers",
			input:    []uint16{'1', '2', '3', '4', '5', 0},
			expected: "12345",
		},
		{
			name:     "String without null terminator",
			input:    []uint16{'T', 'e', 's', 't'},
			expected: "Test",
		},
		{
			name:     "String with early null",
			input:    []uint16{'H', 'i', 0, 'B', 'y', 'e', 0},
			expected: "Hi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sharedmemory.UTF16ToString(tt.input)
			if got != tt.expected {
				t.Errorf("UTF16ToString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGraphics_HelperMethods(t *testing.T) {
	g := &sharedmemory.Graphics{}

	// Test GetCurrentTime
	copy(g.CurrentTime[:], []uint16{'1', ':', '2', '3', '.', '4', '5', '6', 0})
	if got := g.GetCurrentTime(); got != "1:23.456" {
		t.Errorf("GetCurrentTime() = %q, want %q", got, "1:23.456")
	}

	// Test GetLastTime
	copy(g.LastTime[:], []uint16{'1', ':', '2', '5', '.', '1', '2', '3', 0})
	if got := g.GetLastTime(); got != "1:25.123" {
		t.Errorf("GetLastTime() = %q, want %q", got, "1:25.123")
	}

	// Test GetBestTime
	copy(g.BestTime[:], []uint16{'1', ':', '2', '2', '.', '9', '9', '9', 0})
	if got := g.GetBestTime(); got != "1:22.999" {
		t.Errorf("GetBestTime() = %q, want %q", got, "1:22.999")
	}

	// Test GetTyreCompound
	copy(g.TyreCompound[:], []uint16{'S', 'o', 'f', 't', 0})
	if got := g.GetTyreCompound(); got != "Soft" {
		t.Errorf("GetTyreCompound() = %q, want %q", got, "Soft")
	}

	// Test GetTrackStatus
	copy(g.TrackStatus[:], []uint16{'G', 'R', 'E', 'E', 'N', 0})
	if got := g.GetTrackStatus(); got != "GREEN" {
		t.Errorf("GetTrackStatus() = %q, want %q", got, "GREEN")
	}
}

func TestStatic_HelperMethods(t *testing.T) {
	s := &sharedmemory.Static{}

	// Test GetSMVersion
	copy(s.SMVersion[:], []uint16{'1', '.', '8', 0})
	if got := s.GetSMVersion(); got != "1.8" {
		t.Errorf("GetSMVersion() = %q, want %q", got, "1.8")
	}

	// Test GetACVersion
	copy(s.ACVersion[:], []uint16{'1', '.', '9', 0})
	if got := s.GetACVersion(); got != "1.9" {
		t.Errorf("GetACVersion() = %q, want %q", got, "1.9")
	}

	// Test GetCarModel
	copy(s.CarModel[:], []uint16{'4', '8', '8', ' ', 'G', 'T', '3', 0})
	if got := s.GetCarModel(); got != "488 GT3" {
		t.Errorf("GetCarModel() = %q, want %q", got, "488 GT3")
	}

	// Test GetTrack
	copy(s.Track[:], []uint16{'M', 'o', 'n', 'z', 'a', 0})
	if got := s.GetTrack(); got != "Monza" {
		t.Errorf("GetTrack() = %q, want %q", got, "Monza")
	}

	// Test GetPlayerName
	copy(s.PlayerName[:], []uint16{'J', 'o', 'h', 'n', 0})
	if got := s.GetPlayerName(); got != "John" {
		t.Errorf("GetPlayerName() = %q, want %q", got, "John")
	}

	// Test GetPlayerSurname
	copy(s.PlayerSurname[:], []uint16{'D', 'o', 'e', 0})
	if got := s.GetPlayerSurname(); got != "Doe" {
		t.Errorf("GetPlayerSurname() = %q, want %q", got, "Doe")
	}

	// Test GetPlayerNick
	copy(s.PlayerNick[:], []uint16{'J', 'D', '0', '1', 0})
	if got := s.GetPlayerNick(); got != "JD01" {
		t.Errorf("GetPlayerNick() = %q, want %q", got, "JD01")
	}
}

func TestSharedMemoryReader_IsConnected(t *testing.T) {
	reader := sharedmemory.NewSharedMemoryReader()

	// Should not be connected initially
	if reader.IsConnected() {
		t.Error("IsConnected() = true, want false for new reader")
	}

	// Note: We can't test actual connection without ACC running
	// This test just verifies the initial state
}

func TestNewSharedMemoryReader(t *testing.T) {
	reader := sharedmemory.NewSharedMemoryReader()

	if reader == nil {
		t.Fatal("NewSharedMemoryReader() returned nil")
	}

	// Verify initial state
	if reader.IsConnected() {
		t.Error("New reader should not be connected")
	}
}
