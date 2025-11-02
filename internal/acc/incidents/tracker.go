package incidents

import (
	"RaceAll/internal/broadcast"
	"sync"
	"time"
)

// IncidentType representa el tipo de incidente
type IncidentType int

const (
	IncidentTypeAccident IncidentType = iota
	IncidentTypeCollision
	IncidentTypeOfftrack
	IncidentTypeCutting
)

// Incident representa un incidente en la pista
type Incident struct {
	Type         IncidentType
	Timestamp    time.Time
	SessionTime  time.Duration
	CarIndex     uint16
	DriverName   string
	RaceNumber   int32
	Location     string
	Message      string
	Severity     int
	InvolvedCars []uint16
}

// IncidentTracker rastrea incidentes durante la sesión
type IncidentTracker struct {
	incidents          []Incident
	realtimeCarHistory map[float64]map[uint16]*broadcast.RealtimeCarUpdate
	lastAccidentTime   time.Time
	trackDistance      float32
	currentSessionTime time.Duration
	mu                 sync.RWMutex
	callbacks          []func(Incident)
}

const (
	// Tiempo máximo para agrupar accidentes relacionados (ms)
	AccidentGroupingTime = 1000 * time.Millisecond
	// Tiempo de retención del historial de datos (ms)
	HistoryRetentionTime = 20000 * time.Millisecond
)

// NewIncidentTracker crea un nuevo rastreador de incidentes
func NewIncidentTracker() *IncidentTracker {
	return &IncidentTracker{
		incidents:          make([]Incident, 0),
		realtimeCarHistory: make(map[float64]map[uint16]*broadcast.RealtimeCarUpdate),
		callbacks:          make([]func(Incident), 0),
	}
}

// OnIncident registra un callback para cuando ocurra un incidente
func (it *IncidentTracker) OnIncident(callback func(Incident)) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.callbacks = append(it.callbacks, callback)
}

// UpdateTrackData actualiza información del circuito
func (it *IncidentTracker) UpdateTrackData(trackMeters float32) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.trackDistance = trackMeters
}

// UpdateRealtimeCarUpdate añade datos de actualización en tiempo real al historial
func (it *IncidentTracker) UpdateRealtimeCarUpdate(carUpdate *broadcast.RealtimeCarUpdate, sessionTime time.Duration) {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.currentSessionTime = sessionTime

	key := sessionTime.Milliseconds()
	if key == 0 {
		return // Sesión no iniciada
	}

	// Crear entrada en el historial si no existe
	if _, exists := it.realtimeCarHistory[float64(key)]; !exists {
		it.realtimeCarHistory[float64(key)] = make(map[uint16]*broadcast.RealtimeCarUpdate)
	}

	// Actualizar o añadir datos del auto
	it.realtimeCarHistory[float64(key)][carUpdate.CarIndex] = carUpdate

	// Limpiar historial antiguo
	it.cleanOldHistory(float64(key))
}

// HandleAccidentEvent procesa un evento de accidente del broadcast
func (it *IncidentTracker) HandleAccidentEvent(event *broadcast.BroadcastingEvent, carInfo *broadcast.CarInfo) {
	it.mu.Lock()
	defer it.mu.Unlock()

	if carInfo == nil {
		return
	}

	// Los eventos de accidente suelen llegar con 5 segundos de retraso
	correctedSessionTime := it.currentSessionTime - (5000 * time.Millisecond)
	key := it.getValidHistoryKey(float64(correctedSessionTime.Milliseconds()))

	historyData, exists := it.realtimeCarHistory[key]
	if !exists {
		return
	}

	carUpdate, hasCarUpdate := historyData[uint16(event.CarId)]
	if !hasCarUpdate {
		return
	}

	// Crear incidente
	incident := Incident{
		Type:        IncidentTypeAccident,
		Timestamp:   time.Now(),
		SessionTime: time.Duration(key) * time.Millisecond,
		CarIndex:    uint16(event.CarId),
		DriverName:  carInfo.GetCurrentDriverName(),
		RaceNumber:  carInfo.RaceNumber,
		Message:     event.Msg,
		Severity:    1,
	}

	// Determinar ubicación aproximada
	if carUpdate != nil {
		splinePercent := int(carUpdate.SplinePosition * 100)
		incident.Location = formatLocation(splinePercent)
	}

	it.incidents = append(it.incidents, incident)

	// Marcar tiempo de accidente para agrupar
	if it.lastAccidentTime.IsZero() {
		it.lastAccidentTime = time.Now()
	}

	// Notificar callbacks
	for _, callback := range it.callbacks {
		go callback(incident)
	}
}

// GetIncidents devuelve todos los incidentes registrados
func (it *IncidentTracker) GetIncidents() []Incident {
	it.mu.RLock()
	defer it.mu.RUnlock()

	result := make([]Incident, len(it.incidents))
	copy(result, it.incidents)
	return result
}

// GetRecentIncidents devuelve incidentes de los últimos N segundos
func (it *IncidentTracker) GetRecentIncidents(duration time.Duration) []Incident {
	it.mu.RLock()
	defer it.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	result := make([]Incident, 0)

	for _, incident := range it.incidents {
		if incident.Timestamp.After(cutoff) {
			result = append(result, incident)
		}
	}

	return result
}

// Clear limpia todos los incidentes
func (it *IncidentTracker) Clear() {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.incidents = make([]Incident, 0)
	it.realtimeCarHistory = make(map[float64]map[uint16]*broadcast.RealtimeCarUpdate)
	it.lastAccidentTime = time.Time{}
}

// cleanOldHistory limpia datos del historial más antiguos que el tiempo de retención
func (it *IncidentTracker) cleanOldHistory(currentKey float64) {
	for key := range it.realtimeCarHistory {
		if currentKey-key > float64(HistoryRetentionTime.Milliseconds()) {
			delete(it.realtimeCarHistory, key)
		}
	}
}

// getValidHistoryKey encuentra la clave de historial más cercana a la solicitada
func (it *IncidentTracker) getValidHistoryKey(requestedKey float64) float64 {
	if _, exists := it.realtimeCarHistory[requestedKey]; exists {
		return requestedKey
	}

	// Encontrar la clave más cercana
	nearest := 999999.0
	nearestKey := 0.0

	for key := range it.realtimeCarHistory {
		diff := abs(key - requestedKey)
		if diff < nearest {
			nearest = diff
			nearestKey = key
		}
	}

	return nearestKey
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func formatLocation(splinePercent int) string {
	if splinePercent < 10 {
		return "Start/Finish"
	} else if splinePercent < 30 {
		return "Sector 1"
	} else if splinePercent < 60 {
		return "Sector 2"
	} else if splinePercent < 90 {
		return "Sector 3"
	}
	return "Final Sector"
}

// GetIncidentCount devuelve el número total de incidentes
func (it *IncidentTracker) GetIncidentCount() int {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return len(it.incidents)
}
