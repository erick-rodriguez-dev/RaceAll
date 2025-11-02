package session

import (
	"RaceAll/internal/broadcast"
	"fmt"
	"time"
)

// SessionState representa el estado actual de la sesión
type SessionState struct {
	Type           broadcast.RaceSessionType
	Phase          broadcast.SessionPhase
	TimeElapsed    time.Duration
	TimeRemaining  time.Duration
	IsActive       bool
	IsRace         bool
	IsQualifying   bool
	IsPractice     bool
	IsPreSession   bool
	IsPostSession  bool
	SessionNumber  int
	Weather        WeatherConditions
	BestLapTime    int32 // ms
	LastUpdateTime time.Time
}

// WeatherConditions representa las condiciones climáticas
type WeatherConditions struct {
	AmbientTemp byte    // °C
	TrackTemp   byte    // °C
	Clouds      float32 // 0.0 - 1.0
	RainLevel   float32 // 0.0 - 1.0
	Wetness     float32 // 0.0 - 1.0
	IsDry       bool
	IsWet       bool
	IsRaining   bool
}

// SessionTracker realiza seguimiento del estado de la sesión
type SessionTracker struct {
	currentState  *SessionState
	previousState *SessionState
	stateHistory  []SessionState
}

// NewSessionTracker crea un nuevo tracker de sesión
func NewSessionTracker() *SessionTracker {
	return &SessionTracker{
		currentState:  nil,
		previousState: nil,
		stateHistory:  make([]SessionState, 0),
	}
}

// Update actualiza el estado de la sesión con datos del broadcast
func (st *SessionTracker) Update(update *broadcast.RealtimeUpdate) {
	if st.currentState != nil {
		st.previousState = st.currentState
	}

	weather := WeatherConditions{
		AmbientTemp: update.AmbientTemp,
		TrackTemp:   update.TrackTemp,
		Clouds:      update.Clouds,
		RainLevel:   update.RainLevel,
		Wetness:     update.Wetness,
		IsDry:       update.Wetness < 0.1,
		IsWet:       update.Wetness >= 0.1,
		IsRaining:   update.RainLevel > 0.1,
	}

	st.currentState = &SessionState{
		Type:           update.SessionType,
		Phase:          update.Phase,
		TimeElapsed:    update.SessionTime,
		TimeRemaining:  update.SessionEndTime - update.SessionTime,
		IsActive:       update.Phase == broadcast.SessionPhaseSession,
		IsRace:         update.SessionType == broadcast.RaceSessionTypeRace,
		IsQualifying:   update.SessionType == broadcast.RaceSessionTypeQualifying || update.SessionType == broadcast.RaceSessionTypeSuperpole,
		IsPractice:     update.SessionType == broadcast.RaceSessionTypePractice,
		IsPreSession:   update.Phase == broadcast.SessionPhaseNone || update.Phase == broadcast.SessionPhaseStarting,
		IsPostSession:  update.Phase == broadcast.SessionPhaseResultUI,
		SessionNumber:  int(update.SessionIndex),
		Weather:        weather,
		BestLapTime:    update.BestSessionLap.GetLapTimeMS(),
		LastUpdateTime: time.Now(),
	}

	// Guardar en historial si cambió el tipo de sesión
	if st.HasSessionChanged() {
		st.stateHistory = append(st.stateHistory, *st.currentState)
	}
}

// GetCurrentState devuelve el estado actual de la sesión
func (st *SessionTracker) GetCurrentState() *SessionState {
	return st.currentState
}

// HasSessionChanged verifica si cambió el tipo de sesión
func (st *SessionTracker) HasSessionChanged() bool {
	if st.previousState == nil {
		return true
	}
	return st.currentState.Type != st.previousState.Type ||
		st.currentState.SessionNumber != st.previousState.SessionNumber
}

// HasPhaseChanged verifica si cambió la fase de la sesión
func (st *SessionTracker) HasPhaseChanged() bool {
	if st.previousState == nil {
		return true
	}
	return st.currentState.Phase != st.previousState.Phase
}

// IsSessionActive verifica si la sesión está activa
func (st *SessionTracker) IsSessionActive() bool {
	if st.currentState == nil {
		return false
	}
	return st.currentState.IsActive
}

// GetSessionProgress devuelve el progreso de la sesión (0.0 - 1.0)
func (st *SessionTracker) GetSessionProgress() float32 {
	if st.currentState == nil {
		return 0.0
	}

	totalTime := st.currentState.TimeElapsed + st.currentState.TimeRemaining
	if totalTime == 0 {
		return 0.0
	}

	return float32(st.currentState.TimeElapsed) / float32(totalTime)
}

// GetTimeRemainingString devuelve el tiempo restante formateado
func (st *SessionTracker) GetTimeRemainingString() string {
	if st.currentState == nil {
		return "--:--"
	}

	minutes := int(st.currentState.TimeRemaining.Minutes())
	seconds := int(st.currentState.TimeRemaining.Seconds()) % 60

	return formatTime(minutes, seconds)
}

// GetTimeElapsedString devuelve el tiempo transcurrido formateado
func (st *SessionTracker) GetTimeElapsedString() string {
	if st.currentState == nil {
		return "--:--"
	}

	minutes := int(st.currentState.TimeElapsed.Minutes())
	seconds := int(st.currentState.TimeElapsed.Seconds()) % 60

	return formatTime(minutes, seconds)
}

// formatTime formatea minutos y segundos como MM:SS
func formatTime(minutes, seconds int) string {
	if minutes < 0 {
		minutes = 0
	}
	if seconds < 0 {
		seconds = 0
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// GetWeatherDescription devuelve una descripción del clima
func (wc *WeatherConditions) GetDescription() string {
	if wc.IsRaining {
		return "Raining"
	}
	if wc.IsWet {
		return "Wet"
	}
	if wc.Clouds > 0.7 {
		return "Cloudy"
	}
	if wc.Clouds > 0.3 {
		return "Partly Cloudy"
	}
	return "Clear"
}

// IsOptimalConditions verifica si las condiciones son óptimas
func (wc *WeatherConditions) IsOptimalConditions() bool {
	return wc.IsDry &&
		wc.TrackTemp >= 25 && wc.TrackTemp <= 35 &&
		wc.AmbientTemp >= 20 && wc.AmbientTemp <= 30
}

// NeedWetTyres verifica si se necesitan neumáticos de lluvia
func (wc *WeatherConditions) NeedWetTyres() bool {
	return wc.Wetness > 0.3 || wc.IsRaining
}
