package trackposition

import (
	"RaceAll/internal/broadcast"
	"sync"
)

// Car representa un auto en el grafo de posiciones
type Car struct {
	CarIndex         uint16
	LapIndex         int
	SplinePosition   float32
	Location         broadcast.CarLocationEnum
	PreviousLocation broadcast.CarLocationEnum
}

// PositionGraph mantiene el grafo de posiciones de todos los autos
type PositionGraph struct {
	cars map[uint16]*Car
	mu   sync.RWMutex
}

// NewPositionGraph crea un nuevo grafo de posiciones
func NewPositionGraph() *PositionGraph {
	return &PositionGraph{
		cars: make(map[uint16]*Car),
	}
}

// AddCar añade un auto al grafo
func (pg *PositionGraph) AddCar(carIndex uint16) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	if _, exists := pg.cars[carIndex]; !exists {
		pg.cars[carIndex] = &Car{
			CarIndex:         carIndex,
			Location:         broadcast.CarLocationNone,
			LapIndex:         0,
			SplinePosition:   0,
			PreviousLocation: broadcast.CarLocationNone,
		}
	}
}

// GetCar obtiene un auto del grafo
func (pg *PositionGraph) GetCar(carIndex uint16) *Car {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	return pg.cars[carIndex]
}

// RemoveCar elimina un auto del grafo
func (pg *PositionGraph) RemoveCar(carIndex uint16) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	delete(pg.cars, carIndex)
}

// UpdateLocation actualiza la ubicación de un auto
func (pg *PositionGraph) UpdateLocation(carIndex uint16, newSplinePosition float32, newLocation broadcast.CarLocationEnum) {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	car, exists := pg.cars[carIndex]
	if !exists {
		return
	}

	// Detectar cambios de ubicación
	if newLocation != car.Location {
		car.PreviousLocation = car.Location
		car.Location = newLocation
	}

	// Detectar cruce de línea de meta (vuelta completada)
	if car.SplinePosition > newSplinePosition && car.SplinePosition > 0.99 {
		// Vuelta completada en pista
		if newLocation == broadcast.CarLocationTrack && car.Location == broadcast.CarLocationTrack {
			car.LapIndex++
		}

		// Vuelta completada en pitlane
		if newLocation == broadcast.CarLocationPitlane &&
			car.Location == broadcast.CarLocationPitlane &&
			car.PreviousLocation == broadcast.CarLocationPitEntry {
			if car.LapIndex > 0 {
				car.LapIndex++
			}
		}

		// Vuelta completada desde NONE a Track (inicio de sesión)
		if newLocation == broadcast.CarLocationTrack && car.Location == broadcast.CarLocationNone {
			car.PreviousLocation = car.Location
			car.Location = newLocation
		}
	}

	car.SplinePosition = newSplinePosition
}

// Reset reinicia el grafo de posiciones
func (pg *PositionGraph) Reset() {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	for _, car := range pg.cars {
		car.Location = broadcast.CarLocationNone
		car.LapIndex = 0
		car.SplinePosition = 0
		car.PreviousLocation = broadcast.CarLocationNone
	}
}

// Clear limpia completamente el grafo
func (pg *PositionGraph) Clear() {
	pg.mu.Lock()
	defer pg.mu.Unlock()

	pg.cars = make(map[uint16]*Car)
}

// GetAllCars devuelve todos los autos en el grafo
func (pg *PositionGraph) GetAllCars() map[uint16]*Car {
	pg.mu.RLock()
	defer pg.mu.RUnlock()

	result := make(map[uint16]*Car, len(pg.cars))
	for k, v := range pg.cars {
		result[k] = v
	}
	return result
}
