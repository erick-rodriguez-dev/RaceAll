package sharedmemory

import (
	"unicode/utf16"
)

// UTF16ToString converts a UTF16 array to a Go string
func UTF16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s))
}

// Helper methods for Graphics struct
func (g *Graphics) GetCurrentTime() string {
	return UTF16ToString(g.CurrentTime[:])
}

func (g *Graphics) GetLastTime() string {
	return UTF16ToString(g.LastTime[:])
}

func (g *Graphics) GetBestTime() string {
	return UTF16ToString(g.BestTime[:])
}

func (g *Graphics) GetTyreCompound() string {
	return UTF16ToString(g.TyreCompound[:])
}

func (g *Graphics) GetTrackStatus() string {
	return UTF16ToString(g.TrackStatus[:])
}

func (g *Graphics) GetDeltaLapTime() string {
	return UTF16ToString(g.DeltaLapTime[:])
}

func (g *Graphics) GetEstimatedLapTime() string {
	return UTF16ToString(g.EstimatedLapTime[:])
}

// Helper methods for Static struct
func (s *Static) GetSMVersion() string {
	return UTF16ToString(s.SMVersion[:])
}

func (s *Static) GetACVersion() string {
	return UTF16ToString(s.ACVersion[:])
}

func (s *Static) GetCarModel() string {
	return UTF16ToString(s.CarModel[:])
}

func (s *Static) GetTrack() string {
	return UTF16ToString(s.Track[:])
}

func (s *Static) GetPlayerName() string {
	return UTF16ToString(s.PlayerName[:])
}

func (s *Static) GetPlayerSurname() string {
	return UTF16ToString(s.PlayerSurname[:])
}

func (s *Static) GetPlayerNick() string {
	return UTF16ToString(s.PlayerNick[:])
}

func (s *Static) GetTrackConfiguration() string {
	return UTF16ToString(s.TrackConfiguration[:])
}

func (s *Static) GetCarSkin() string {
	return UTF16ToString(s.CarSkin[:])
}

func (s *Static) GetDryTyresName() string {
	return UTF16ToString(s.DryTyresName[:])
}

func (s *Static) GetWetTyresName() string {
	return UTF16ToString(s.WetTyresName[:])
}
