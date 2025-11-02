package broadcast_test

import (
	"testing"

	"RaceAll/internal/broadcast"
)

func TestDriverInfo_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		driver   broadcast.DriverInfo
		expected string
	}{
		{
			name: "Normal name",
			driver: broadcast.DriverInfo{
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: "John Doe",
		},
		{
			name: "Empty first name",
			driver: broadcast.DriverInfo{
				FirstName: "",
				LastName:  "Doe",
			},
			expected: " Doe",
		},
		{
			name: "Empty last name",
			driver: broadcast.DriverInfo{
				FirstName: "John",
				LastName:  "",
			},
			expected: "John ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.driver.GetFullName()
			if got != tt.expected {
				t.Errorf("GetFullName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestLapInfo_GetLapTimeMS(t *testing.T) {
	split1 := int32(30000)
	split2 := int32(32000)
	split3 := int32(31000)

	tests := []struct {
		name     string
		lapInfo  broadcast.LapInfo
		expected int32
	}{
		{
			name: "All splits valid",
			lapInfo: broadcast.LapInfo{
				Splits: [3]*int32{&split1, &split2, &split3},
			},
			expected: 93000,
		},
		{
			name: "One nil split",
			lapInfo: broadcast.LapInfo{
				Splits: [3]*int32{&split1, nil, &split3},
			},
			expected: 61000,
		},
		{
			name: "All nil splits",
			lapInfo: broadcast.LapInfo{
				Splits: [3]*int32{nil, nil, nil},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lapInfo.GetLapTimeMS()
			if got != tt.expected {
				t.Errorf("GetLapTimeMS() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestLapInfo_HasValidLapTime(t *testing.T) {
	validTime := int32(93000)
	invalidTime := int32(2147483647) // math.MaxInt32

	tests := []struct {
		name     string
		lapInfo  broadcast.LapInfo
		expected bool
	}{
		{
			name: "Valid lap time",
			lapInfo: broadcast.LapInfo{
				LaptimeMS: &validTime,
			},
			expected: true,
		},
		{
			name: "Invalid lap time (MaxInt32)",
			lapInfo: broadcast.LapInfo{
				LaptimeMS: &invalidTime,
			},
			expected: false,
		},
		{
			name: "Nil lap time",
			lapInfo: broadcast.LapInfo{
				LaptimeMS: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lapInfo.HasValidLapTime()
			if got != tt.expected {
				t.Errorf("HasValidLapTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCarInfo_GetCurrentDriver(t *testing.T) {
	drivers := []broadcast.DriverInfo{
		{FirstName: "Driver", LastName: "One"},
		{FirstName: "Driver", LastName: "Two"},
		{FirstName: "Driver", LastName: "Three"},
	}

	tests := []struct {
		name           string
		carInfo        broadcast.CarInfo
		expectedDriver *broadcast.DriverInfo
	}{
		{
			name: "First driver",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 0,
				Drivers:            drivers,
			},
			expectedDriver: &drivers[0],
		},
		{
			name: "Second driver",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 1,
				Drivers:            drivers,
			},
			expectedDriver: &drivers[1],
		},
		{
			name: "Out of range index",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 10,
				Drivers:            drivers,
			},
			expectedDriver: nil,
		},
		{
			name: "Empty drivers list",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 0,
				Drivers:            []broadcast.DriverInfo{},
			},
			expectedDriver: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.carInfo.GetCurrentDriver()
			if (got == nil) != (tt.expectedDriver == nil) {
				t.Errorf("GetCurrentDriver() returned nil = %v, want nil = %v", got == nil, tt.expectedDriver == nil)
				return
			}
			if got != nil && tt.expectedDriver != nil {
				if got.FirstName != tt.expectedDriver.FirstName || got.LastName != tt.expectedDriver.LastName {
					t.Errorf("GetCurrentDriver() = %v %v, want %v %v",
						got.FirstName, got.LastName,
						tt.expectedDriver.FirstName, tt.expectedDriver.LastName)
				}
			}
		})
	}
}

func TestCarInfo_GetCurrentDriverName(t *testing.T) {
	drivers := []broadcast.DriverInfo{
		{FirstName: "John", LastName: "Doe"},
		{FirstName: "Jane", LastName: "Smith"},
	}

	tests := []struct {
		name     string
		carInfo  broadcast.CarInfo
		expected string
	}{
		{
			name: "Valid driver",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 0,
				Drivers:            drivers,
			},
			expected: "John Doe",
		},
		{
			name: "Invalid index",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 10,
				Drivers:            drivers,
			},
			expected: "Unknown",
		},
		{
			name: "Empty drivers",
			carInfo: broadcast.CarInfo{
				CurrentDriverIndex: 0,
				Drivers:            []broadcast.DriverInfo{},
			},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.carInfo.GetCurrentDriverName()
			if got != tt.expected {
				t.Errorf("GetCurrentDriverName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetCarModelName(t *testing.T) {
	tests := []struct {
		name     string
		model    byte
		expected string
	}{
		{"Porsche 911 GT3 R", 0, "Porsche 911 GT3 R 2018"},
		{"Mercedes-AMG GT3", 1, "Mercedes-AMG GT3 2015"},
		{"Ferrari 488 GT3", 2, "Ferrari 488 GT3 2018"},
		{"Audi R8 LMS", 3, "Audi R8 LMS 2015"},
		{"Unknown model", 255, "Unknown Car"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := broadcast.GetCarModelName(tt.model)
			if got != tt.expected {
				t.Errorf("GetCarModelName(%d) = %q, want %q", tt.model, got, tt.expected)
			}
		})
	}
}
