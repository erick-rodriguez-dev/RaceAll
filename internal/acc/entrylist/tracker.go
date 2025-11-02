package entrylist

import (
	"RaceAll/internal/broadcast"
	"sync"
	"time"
)

// DriverInfo representa información de un piloto
type DriverInfo struct {
	FirstName   string
	LastName    string
	ShortName   string
	Category    broadcast.DriverCategory
	Nationality broadcast.NationalityEnum
}

// CarData representa toda la información de un auto en la entrada
type CarData struct {
	// Información básica
	CarIndex           uint16
	CarModelType       byte
	TeamName           string
	RaceNumber         int32
	CupCategory        byte
	CurrentDriverIndex uint8
	Nationality        broadcast.NationalityEnum

	// Drivers
	Drivers []DriverInfo

	// Datos en tiempo real
	RealtimeUpdate *broadcast.RealtimeCarUpdate

	// Control
	LastUpdate time.Time
}

// EntryListTracker rastrea todos los autos participantes en la sesión
type EntryListTracker struct {
	cars map[uint16]*CarData
	mu   sync.RWMutex
}

// NewEntryListTracker crea un nuevo rastreador de lista de entrada
func NewEntryListTracker() *EntryListTracker {
	return &EntryListTracker{
		cars: make(map[uint16]*CarData),
	}
}

// UpdateCarInfo actualiza la información estática de un auto
func (elt *EntryListTracker) UpdateCarInfo(carInfo *broadcast.CarInfo) {
	elt.mu.Lock()
	defer elt.mu.Unlock()

	carIndex := carInfo.CarIndex

	// Crear o obtener CarData
	carData, exists := elt.cars[carIndex]
	if !exists {
		carData = &CarData{
			CarIndex: carIndex,
		}
		elt.cars[carIndex] = carData
	}

	// Actualizar información
	carData.CarModelType = carInfo.CarModelType
	carData.TeamName = carInfo.TeamName
	carData.RaceNumber = carInfo.RaceNumber
	carData.CupCategory = carInfo.CupCategory
	carData.CurrentDriverIndex = carInfo.CurrentDriverIndex
	carData.Nationality = carInfo.Nationality

	// Actualizar drivers
	carData.Drivers = make([]DriverInfo, len(carInfo.Drivers))
	for i, driver := range carInfo.Drivers {
		carData.Drivers[i] = DriverInfo{
			FirstName:   driver.FirstName,
			LastName:    driver.LastName,
			ShortName:   driver.ShortName,
			Category:    driver.Category,
			Nationality: driver.Nationality,
		}
	}

	carData.LastUpdate = time.Now()
}

// UpdateRealtimeCarUpdate actualiza los datos en tiempo real de un auto
func (elt *EntryListTracker) UpdateRealtimeCarUpdate(carUpdate *broadcast.RealtimeCarUpdate) {
	elt.mu.Lock()
	defer elt.mu.Unlock()

	carIndex := carUpdate.CarIndex

	// Crear o obtener CarData
	carData, exists := elt.cars[carIndex]
	if !exists {
		carData = &CarData{
			CarIndex: carIndex,
		}
		elt.cars[carIndex] = carData
	}

	carData.RealtimeUpdate = carUpdate
	carData.LastUpdate = time.Now()
}

// GetCarData devuelve los datos de un auto específico
func (elt *EntryListTracker) GetCarData(carIndex uint16) *CarData {
	elt.mu.RLock()
	defer elt.mu.RUnlock()

	return elt.cars[carIndex]
}

// GetAllCars devuelve todos los autos en la lista de entrada
func (elt *EntryListTracker) GetAllCars() map[uint16]*CarData {
	elt.mu.RLock()
	defer elt.mu.RUnlock()

	// Retornar copia del mapa para evitar problemas de concurrencia
	result := make(map[uint16]*CarData, len(elt.cars))
	for k, v := range elt.cars {
		result[k] = v
	}
	return result
}

// GetCarCount devuelve el número de autos en la lista
func (elt *EntryListTracker) GetCarCount() int {
	elt.mu.RLock()
	defer elt.mu.RUnlock()
	return len(elt.cars)
}

// Cleanup limpia autos que no han recibido actualizaciones en un tiempo
func (elt *EntryListTracker) Cleanup(maxAge time.Duration) {
	elt.mu.Lock()
	defer elt.mu.Unlock()

	now := time.Now()
	for carIndex, carData := range elt.cars {
		if now.Sub(carData.LastUpdate) > maxAge {
			delete(elt.cars, carIndex)
		}
	}
}

// Clear limpia toda la lista de entrada
func (elt *EntryListTracker) Clear() {
	elt.mu.Lock()
	defer elt.mu.Unlock()

	elt.cars = make(map[uint16]*CarData)
}

// GetCurrentDriverName devuelve el nombre del piloto actual de un auto
func (cd *CarData) GetCurrentDriverName() string {
	if cd == nil || len(cd.Drivers) == 0 {
		return ""
	}

	if int(cd.CurrentDriverIndex) >= len(cd.Drivers) {
		return ""
	}

	driver := cd.Drivers[cd.CurrentDriverIndex]
	if driver.ShortName != "" {
		return driver.ShortName
	}
	return driver.FirstName + " " + driver.LastName
}

// GetCurrentDriver devuelve el piloto actual
func (cd *CarData) GetCurrentDriver() *DriverInfo {
	if cd == nil || len(cd.Drivers) == 0 {
		return nil
	}

	if int(cd.CurrentDriverIndex) >= len(cd.Drivers) {
		return nil
	}

	return &cd.Drivers[cd.CurrentDriverIndex]
}
