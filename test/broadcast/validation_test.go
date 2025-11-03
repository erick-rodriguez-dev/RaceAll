package broadcast_test

import (
	"testing"

	"RaceAll/internal/broadcast"
)

func TestValidateCarIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   uint16
		wantErr bool
	}{
		{"Valid index 0", 0, false},
		{"Valid index 50", 50, false},
		{"Valid index 999", 999, false},
		{"Invalid index 1000", 1000, true},
		{"Invalid index 2000", 2000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateCarIndex(tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCarIndex(%d) error = %v, wantErr %v", tt.index, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionType(t *testing.T) {
	tests := []struct {
		name        string
		sessionType broadcast.RaceSessionType
		wantErr     bool
	}{
		{"Practice session", broadcast.RaceSessionTypePractice, false},
		{"Qualifying session", broadcast.RaceSessionTypeQualifying, false},
		{"Race session", broadcast.RaceSessionTypeRace, false},
		{"Invalid session", broadcast.RaceSessionType(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateSessionType(tt.sessionType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionType(%v) error = %v, wantErr %v", tt.sessionType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionPhase(t *testing.T) {
	tests := []struct {
		name    string
		phase   broadcast.SessionPhase
		wantErr bool
	}{
		{"None phase", broadcast.SessionPhaseNone, false},
		{"Starting phase", broadcast.SessionPhaseStarting, false},
		{"Session phase", broadcast.SessionPhaseSession, false},
		{"Invalid phase", broadcast.SessionPhase(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateSessionPhase(tt.phase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionPhase(%v) error = %v, wantErr %v", tt.phase, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCarLocation(t *testing.T) {
	tests := []struct {
		name     string
		location broadcast.CarLocationEnum
		wantErr  bool
	}{
		{"None location", broadcast.CarLocationNone, false},
		{"Track location", broadcast.CarLocationTrack, false},
		{"Pitlane location", broadcast.CarLocationPitlane, false},
		{"Invalid location", broadcast.CarLocationEnum(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateCarLocation(tt.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCarLocation(%v) error = %v, wantErr %v", tt.location, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBroadcastingEventType(t *testing.T) {
	tests := []struct {
		name      string
		eventType broadcast.BroadcastingEventType
		wantErr   bool
	}{
		{"None event", broadcast.BroadcastingEventTypeNone, false},
		{"Green flag event", broadcast.BroadcastingEventTypeGreenFlag, false},
		{"Accident event", broadcast.BroadcastingEventTypeAccident, false},
		{"Invalid event", broadcast.BroadcastingEventType(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateBroadcastingEventType(tt.eventType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBroadcastingEventType(%v) error = %v, wantErr %v", tt.eventType, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDriverCategory(t *testing.T) {
	tests := []struct {
		name     string
		category broadcast.DriverCategory
		wantErr  bool
	}{
		{"Bronze category", broadcast.DriverCategoryBronze, false},
		{"Silver category", broadcast.DriverCategorySilver, false},
		{"Gold category", broadcast.DriverCategoryGold, false},
		{"Platinum category", broadcast.DriverCategoryPlatinum, false},
		{"Invalid category", broadcast.DriverCategory(99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := broadcast.ValidateDriverCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDriverCategory(%v) error = %v, wantErr %v", tt.category, err, tt.wantErr)
			}
		})
	}
}
