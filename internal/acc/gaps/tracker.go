package gaps

import (
	"math"
	"sync"
	"time"
)

const (
	GapDistanceMeter  = 50
	MeasuringInterval = 10 // milliseconds
)

// GapPointData representa un punto de medición de gap en el circuito
type GapPointData struct {
	PassedAt time.Time
}

// GapTracker rastrea los gaps entre autos basándose en posiciones en el circuito
type GapTracker struct {
	gapData   map[uint16][]GapPointData
	totalGaps int
	mu        sync.RWMutex
}

// NewGapTracker crea un nuevo rastreador de gaps
func NewGapTracker() *GapTracker {
	return &GapTracker{
		gapData: make(map[uint16][]GapPointData),
	}
}

// Initialize inicializa el tracker con la distancia del circuito
func (gt *GapTracker) Initialize(trackMeters float32) {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gt.totalGaps = int(math.Floor(float64(trackMeters) / GapDistanceMeter))
	gt.gapData = make(map[uint16][]GapPointData)
}

// UpdateCarPosition actualiza la posición de un auto en el circuito
func (gt *GapTracker) UpdateCarPosition(carIndex uint16, splinePosition float32) {
	if gt.totalGaps <= 0 {
		return
	}

	gt.mu.Lock()
	defer gt.mu.Unlock()

	// Crear array de gaps si no existe para este auto
	if _, exists := gt.gapData[carIndex]; !exists {
		gt.gapData[carIndex] = make([]GapPointData, gt.totalGaps)
	}

	data := gt.gapData[carIndex]
	gapStep := 1.0 / float32(gt.totalGaps)

	// Actualizar puntos de gap cercanos a la posición actual
	for i := 0; i < gt.totalGaps; i++ {
		estimatedGapSpline := float32(i) * gapStep
		if splinePosition > estimatedGapSpline && splinePosition < estimatedGapSpline+gapStep*1.5 {
			// Solo actualizar si es la primera vez o ha pasado más de 1 minuto
			if data[i].PassedAt.IsZero() || time.Since(data[i].PassedAt) > time.Minute {
				data[i] = GapPointData{PassedAt: time.Now()}
			}
		}
	}
}

// TimeGapBetween calcula el gap en segundos entre dos autos
func (gt *GapTracker) TimeGapBetween(currentCarIndex uint16, splineCurrent float32, carAheadIndex uint16) float32 {
	return gt.timeGapBetweenInternal(currentCarIndex, splineCurrent, carAheadIndex, false)
}

func (gt *GapTracker) timeGapBetweenInternal(currentCarIndex uint16, splineCurrent float32, carAheadIndex uint16, retry bool) float32 {
	if gt.totalGaps <= 0 {
		return -1
	}

	gt.mu.RLock()
	defer gt.mu.RUnlock()

	gapsA, existsA := gt.gapData[currentCarIndex]
	gapsB, existsB := gt.gapData[carAheadIndex]

	if !existsA || !existsB {
		return -1
	}

	// Estimar índice basado en posición spline
	estimatedIndex := int(float32(gt.totalGaps) * splineCurrent)
	if estimatedIndex < 0 {
		estimatedIndex = 0
	}
	if estimatedIndex >= gt.totalGaps {
		estimatedIndex = gt.totalGaps - 1
	}

	passedAtA := gapsA[estimatedIndex].PassedAt
	passedAtB := gapsB[estimatedIndex].PassedAt

	if passedAtA.IsZero() || passedAtB.IsZero() {
		return -1
	}

	gap := passedAtA.Sub(passedAtB)

	if gap >= 0 {
		return float32(gap.Seconds())
	}

	// Si el gap es negativo y no es un retry, intentar con posición anterior
	if !retry {
		splineGap := 1.0 / float32(gt.totalGaps)
		return gt.timeGapBetweenInternal(currentCarIndex, splineCurrent-splineGap, carAheadIndex, true)
	}

	return -1
}

// Clear limpia todos los datos de gaps
func (gt *GapTracker) Clear() {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gt.gapData = make(map[uint16][]GapPointData)
}

// Reset reinicia el tracker
func (gt *GapTracker) Reset() {
	gt.Clear()
	gt.totalGaps = 0
}

// GetTotalGaps devuelve el número total de puntos de medición
func (gt *GapTracker) GetTotalGaps() int {
	gt.mu.RLock()
	defer gt.mu.RUnlock()
	return gt.totalGaps
}
