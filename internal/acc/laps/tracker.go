package laps

import (
	"RaceAll/internal/broadcast"
	"math"
	"time"
)

// LapData representa los datos completos de una vuelta
type LapData struct {
	LapNumber      int
	LapTime        int32 // ms
	Sector1        int32 // ms
	Sector2        int32 // ms
	Sector3        int32 // ms
	IsValid        bool
	IsValidForBest bool
	IsOutlap       bool
	IsInlap        bool
	IsPersonalBest bool
	IsSessionBest  bool
	Timestamp      time.Time
	CarIndex       uint16
	DriverIndex    uint16
}

// LapTracker realiza seguimiento de vueltas y tiempos
type LapTracker struct {
	laps            []LapData
	currentLap      *LapData
	bestLap         *LapData
	lastLap         *LapData
	bestSector1     int32
	bestSector2     int32
	bestSector3     int32
	currentSector   int
	sectorStartTime int32
	carIndex        uint16
}

// NewLapTracker crea un nuevo tracker de vueltas
func NewLapTracker(carIndex uint16) *LapTracker {
	return &LapTracker{
		laps:            make([]LapData, 0),
		currentLap:      nil,
		bestLap:         nil,
		lastLap:         nil,
		bestSector1:     math.MaxInt32,
		bestSector2:     math.MaxInt32,
		bestSector3:     math.MaxInt32,
		currentSector:   0,
		sectorStartTime: 0,
		carIndex:        carIndex,
	}
}

// UpdateFromBroadcast actualiza el tracker con datos del broadcast
func (lt *LapTracker) UpdateFromBroadcast(lapInfo *broadcast.LapInfo) {
	if lapInfo == nil {
		return
	}

	// Crear nueva vuelta
	lap := LapData{
		LapNumber:      len(lt.laps) + 1,
		LapTime:        lapInfo.GetLapTimeMS(),
		IsValid:        !lapInfo.IsInvalid,
		IsValidForBest: lapInfo.IsValidForBest,
		IsOutlap:       lapInfo.Type == broadcast.LapTypeOutlap,
		IsInlap:        lapInfo.Type == broadcast.LapTypeInlap,
		Timestamp:      time.Now(),
		CarIndex:       lapInfo.CarIndex,
		DriverIndex:    lapInfo.DriverIndex,
	}

	// Extraer tiempos de sectores
	if len(lapInfo.Splits) >= 1 && lapInfo.Splits[0] != nil {
		lap.Sector1 = *lapInfo.Splits[0]
	}
	if len(lapInfo.Splits) >= 2 && lapInfo.Splits[1] != nil {
		lap.Sector2 = *lapInfo.Splits[1]
	}
	if len(lapInfo.Splits) >= 3 && lapInfo.Splits[2] != nil {
		lap.Sector3 = *lapInfo.Splits[2]
	}

	// Actualizar mejores sectores
	if lap.IsValidForBest {
		if lap.Sector1 > 0 && lap.Sector1 < lt.bestSector1 {
			lt.bestSector1 = lap.Sector1
		}
		if lap.Sector2 > 0 && lap.Sector2 < lt.bestSector2 {
			lt.bestSector2 = lap.Sector2
		}
		if lap.Sector3 > 0 && lap.Sector3 < lt.bestSector3 {
			lt.bestSector3 = lap.Sector3
		}
	}

	// Verificar si es mejor vuelta personal
	if lap.IsValidForBest && lap.LapTime > 0 {
		if lt.bestLap == nil || lap.LapTime < lt.bestLap.LapTime {
			lap.IsPersonalBest = true
			lt.bestLap = &lap
		}
	}

	// Guardar vuelta
	lt.lastLap = &lap
	lt.laps = append(lt.laps, lap)
}

// GetBestLap devuelve la mejor vuelta
func (lt *LapTracker) GetBestLap() *LapData {
	return lt.bestLap
}

// GetLastLap devuelve la última vuelta completada
func (lt *LapTracker) GetLastLap() *LapData {
	return lt.lastLap
}

// GetLapCount devuelve el número de vueltas completadas
func (lt *LapTracker) GetLapCount() int {
	return len(lt.laps)
}

// GetValidLapCount devuelve el número de vueltas válidas
func (lt *LapTracker) GetValidLapCount() int {
	count := 0
	for _, lap := range lt.laps {
		if lap.IsValid && !lap.IsOutlap && !lap.IsInlap {
			count++
		}
	}
	return count
}

// GetAverageLapTime calcula el tiempo promedio de vuelta (solo vueltas válidas)
func (lt *LapTracker) GetAverageLapTime() int32 {
	if len(lt.laps) == 0 {
		return 0
	}

	var total int64
	count := 0

	for _, lap := range lt.laps {
		if lap.IsValid && !lap.IsOutlap && !lap.IsInlap && lap.LapTime > 0 {
			total += int64(lap.LapTime)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return int32(total / int64(count))
}

// GetBestSector devuelve el mejor tiempo de un sector específico (1-3)
func (lt *LapTracker) GetBestSector(sector int) int32 {
	switch sector {
	case 1:
		return lt.bestSector1
	case 2:
		return lt.bestSector2
	case 3:
		return lt.bestSector3
	default:
		return 0
	}
}

// GetTheoreticalBest calcula la mejor vuelta teórica (suma de mejores sectores)
func (lt *LapTracker) GetTheoreticalBest() int32 {
	if lt.bestSector1 == math.MaxInt32 ||
		lt.bestSector2 == math.MaxInt32 ||
		lt.bestSector3 == math.MaxInt32 {
		return 0
	}
	return lt.bestSector1 + lt.bestSector2 + lt.bestSector3
}

// CanImprove verifica si se puede mejorar la mejor vuelta con los sectores actuales
func (lt *LapTracker) CanImprove() bool {
	if lt.bestLap == nil {
		return true
	}
	theoretical := lt.GetTheoreticalBest()
	return theoretical > 0 && theoretical < lt.bestLap.LapTime
}

// GetDeltaToBest calcula el delta respecto a la mejor vuelta
func (lt *LapTracker) GetDeltaToBest(currentLapTime int32) int32 {
	if lt.bestLap == nil || currentLapTime <= 0 {
		return 0
	}
	return currentLapTime - lt.bestLap.LapTime
}

// GetDeltaToLast calcula el delta respecto a la última vuelta
func (lt *LapTracker) GetDeltaToLast(currentLapTime int32) int32 {
	if lt.lastLap == nil || currentLapTime <= 0 {
		return 0
	}
	return currentLapTime - lt.lastLap.LapTime
}

// GetConsistency calcula la consistencia de vueltas (desviación estándar)
func (lt *LapTracker) GetConsistency() float64 {
	validLaps := make([]int32, 0)

	for _, lap := range lt.laps {
		if lap.IsValid && !lap.IsOutlap && !lap.IsInlap && lap.LapTime > 0 {
			validLaps = append(validLaps, lap.LapTime)
		}
	}

	if len(validLaps) < 2 {
		return 0.0
	}

	// Calcular media
	var sum int64
	for _, lapTime := range validLaps {
		sum += int64(lapTime)
	}
	mean := float64(sum) / float64(len(validLaps))

	// Calcular desviación estándar
	var variance float64
	for _, lapTime := range validLaps {
		diff := float64(lapTime) - mean
		variance += diff * diff
	}
	variance /= float64(len(validLaps))

	return math.Sqrt(variance)
}

// GetLapHistory devuelve todas las vueltas
func (lt *LapTracker) GetLapHistory() []LapData {
	return lt.laps
}

// GetRecentLaps devuelve las últimas N vueltas
func (lt *LapTracker) GetRecentLaps(count int) []LapData {
	if count <= 0 {
		return []LapData{}
	}

	totalLaps := len(lt.laps)
	if count > totalLaps {
		count = totalLaps
	}

	return lt.laps[totalLaps-count:]
}

// IsImproving verifica si el rendimiento está mejorando (últimas 3 vueltas)
func (lt *LapTracker) IsImproving() bool {
	recentLaps := lt.GetRecentLaps(3)
	if len(recentLaps) < 3 {
		return false
	}

	// Verificar si las últimas 3 vueltas están mejorando
	return recentLaps[2].LapTime < recentLaps[1].LapTime &&
		recentLaps[1].LapTime < recentLaps[0].LapTime
}
