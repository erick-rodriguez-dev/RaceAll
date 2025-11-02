package tracks

import "RaceAll/internal/broadcast"

// TrackID representa los identificadores de circuitos en ACC
type TrackID = int32

// TrackInfo contiene información detallada de un circuito
type TrackInfo struct {
	ID           TrackID
	Name         string
	Country      string
	LengthMeters int32
	Sectors      int
	Corners      int
	PitlaneSpeed int // km/h
}

// GetTrackInfoFromBroadcast convierte TrackData de broadcast a TrackInfo enriquecido
func GetTrackInfoFromBroadcast(trackData broadcast.TrackData) TrackInfo {
	// Obtener información adicional basada en el ID
	additionalInfo := getAdditionalTrackInfo(trackData.TrackId)

	return TrackInfo{
		ID:           trackData.TrackId,
		Name:         trackData.TrackName,
		Country:      additionalInfo.Country,
		LengthMeters: trackData.TrackMeters,
		Sectors:      additionalInfo.Sectors,
		Corners:      additionalInfo.Corners,
		PitlaneSpeed: additionalInfo.PitlaneSpeed,
	}
}

// GetTrackInfo devuelve información de un circuito solo por ID (sin datos de broadcast)
func GetTrackInfo(trackID TrackID) TrackInfo {
	additionalInfo := getAdditionalTrackInfo(trackID)

	return TrackInfo{
		ID:           trackID,
		Name:         additionalInfo.Name,
		Country:      additionalInfo.Country,
		LengthMeters: additionalInfo.LengthMeters,
		Sectors:      additionalInfo.Sectors,
		Corners:      additionalInfo.Corners,
		PitlaneSpeed: additionalInfo.PitlaneSpeed,
	}
}

// AdditionalTrackInfo contiene información complementaria del circuito
type AdditionalTrackInfo struct {
	Name         string
	Country      string
	LengthMeters int32
	Sectors      int
	Corners      int
	PitlaneSpeed int
}

// getAdditionalTrackInfo devuelve información adicional basada en el ID
func getAdditionalTrackInfo(trackID TrackID) AdditionalTrackInfo {
	// Base de datos de información adicional por ID de circuito
	trackDatabase := map[TrackID]AdditionalTrackInfo{
		0:  {Name: "Monza", Country: "Italy", LengthMeters: 5793, Sectors: 3, Corners: 11, PitlaneSpeed: 80},
		1:  {Name: "Zolder", Country: "Belgium", LengthMeters: 4011, Sectors: 3, Corners: 10, PitlaneSpeed: 60},
		2:  {Name: "Brands Hatch", Country: "United Kingdom", LengthMeters: 3908, Sectors: 3, Corners: 9, PitlaneSpeed: 60},
		3:  {Name: "Silverstone", Country: "United Kingdom", LengthMeters: 5891, Sectors: 3, Corners: 18, PitlaneSpeed: 80},
		4:  {Name: "Paul Ricard", Country: "France", LengthMeters: 5842, Sectors: 3, Corners: 15, PitlaneSpeed: 80},
		5:  {Name: "Misano", Country: "Italy", LengthMeters: 4226, Sectors: 3, Corners: 16, PitlaneSpeed: 60},
		6:  {Name: "Spa-Francorchamps", Country: "Belgium", LengthMeters: 7004, Sectors: 3, Corners: 19, PitlaneSpeed: 60},
		7:  {Name: "Nürburgring", Country: "Germany", LengthMeters: 5137, Sectors: 3, Corners: 15, PitlaneSpeed: 60},
		8:  {Name: "Barcelona", Country: "Spain", LengthMeters: 4655, Sectors: 3, Corners: 16, PitlaneSpeed: 80},
		9:  {Name: "Hungaroring", Country: "Hungary", LengthMeters: 4381, Sectors: 3, Corners: 14, PitlaneSpeed: 80},
		10: {Name: "Zandvoort", Country: "Netherlands", LengthMeters: 4259, Sectors: 3, Corners: 14, PitlaneSpeed: 80},
		11: {Name: "Kyalami", Country: "South Africa", LengthMeters: 4522, Sectors: 3, Corners: 16, PitlaneSpeed: 60},
		12: {Name: "Mount Panorama", Country: "Australia", LengthMeters: 6213, Sectors: 3, Corners: 23, PitlaneSpeed: 60},
		13: {Name: "Suzuka", Country: "Japan", LengthMeters: 5807, Sectors: 3, Corners: 18, PitlaneSpeed: 80},
		14: {Name: "Laguna Seca", Country: "USA", LengthMeters: 3602, Sectors: 3, Corners: 11, PitlaneSpeed: 55},
		15: {Name: "Imola", Country: "Italy", LengthMeters: 4909, Sectors: 3, Corners: 19, PitlaneSpeed: 80},
		16: {Name: "Oulton Park", Country: "United Kingdom", LengthMeters: 4332, Sectors: 3, Corners: 16, PitlaneSpeed: 60},
		17: {Name: "Donington", Country: "United Kingdom", LengthMeters: 4023, Sectors: 3, Corners: 12, PitlaneSpeed: 60},
		18: {Name: "Snetterton", Country: "United Kingdom", LengthMeters: 4778, Sectors: 3, Corners: 9, PitlaneSpeed: 60},
		19: {Name: "COTA", Country: "USA", LengthMeters: 5513, Sectors: 3, Corners: 20, PitlaneSpeed: 80},
		20: {Name: "Indianapolis", Country: "USA", LengthMeters: 4024, Sectors: 3, Corners: 14, PitlaneSpeed: 60},
		21: {Name: "Watkins Glen", Country: "USA", LengthMeters: 5472, Sectors: 3, Corners: 11, PitlaneSpeed: 55},
		22: {Name: "Valencia", Country: "Spain", LengthMeters: 4005, Sectors: 3, Corners: 14, PitlaneSpeed: 60},
		23: {Name: "Red Bull Ring", Country: "Austria", LengthMeters: 4318, Sectors: 3, Corners: 10, PitlaneSpeed: 80},
	}

	if info, exists := trackDatabase[trackID]; exists {
		return info
	}

	// Default unknown track
	return AdditionalTrackInfo{
		Name:         "Unknown Track",
		Country:      "Unknown",
		LengthMeters: 5000,
		Sectors:      3,
		Corners:      15,
		PitlaneSpeed: 80,
	}
}

// GetTrackName devuelve el nombre del circuito
func GetTrackName(trackID TrackID) string {
	return getAdditionalTrackInfo(trackID).Name
}

// GetTrackLength devuelve la longitud del circuito en metros
func GetTrackLength(trackID TrackID) int32 {
	return getAdditionalTrackInfo(trackID).LengthMeters
}

// GetTrackLengthKm devuelve la longitud del circuito en kilómetros
func GetTrackLengthKm(trackID TrackID) float32 {
	return float32(getAdditionalTrackInfo(trackID).LengthMeters) / 1000.0
}

// GetPitlaneSpeed devuelve el límite de velocidad del pitlane
func GetPitlaneSpeed(trackID TrackID) int {
	return getAdditionalTrackInfo(trackID).PitlaneSpeed
}

// GetSectorCount devuelve el número de sectores
func GetSectorCount(trackID TrackID) int {
	return getAdditionalTrackInfo(trackID).Sectors
}
