package sessiontime

import (
	"sync"
	"time"
)

// SessionTimeTracker rastrea el multiplicador de tiempo de la sesión
type SessionTimeTracker struct {
	timeMultiplier      int
	lastServerMillis    int64
	serverChanges       int
	isTracking          bool
	mu                  sync.RWMutex
	multiplierCallbacks []func(int)
}

// NewSessionTimeTracker crea un nuevo rastreador de tiempo de sesión
func NewSessionTimeTracker() *SessionTimeTracker {
	return &SessionTimeTracker{
		timeMultiplier:      -1,
		lastServerMillis:    -1,
		serverChanges:       0,
		multiplierCallbacks: make([]func(int), 0),
	}
}

// OnMultiplierChanged registra un callback para cambios en el multiplicador
func (st *SessionTimeTracker) OnMultiplierChanged(callback func(int)) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.multiplierCallbacks = append(st.multiplierCallbacks, callback)
}

// Update actualiza el tracker con el tiempo del día del servidor
func (st *SessionTimeTracker) Update(timeOfDay time.Duration) {
	st.mu.Lock()
	defer st.mu.Unlock()

	newMilliseconds := int64(timeOfDay.Milliseconds())

	if newMilliseconds != st.lastServerMillis {
		serverTimeChange := float64(newMilliseconds - st.lastServerMillis)

		// Validar que el cambio sea razonable (menos de 4 minutos y al menos 1ms)
		if serverTimeChange < 240000 && serverTimeChange >= 1 {
			if st.lastServerMillis != -1 && st.serverChanges > 3 {
				// Redondear el cambio de tiempo
				serverTimeChange = float64(int64(serverTimeChange + 0.5))     // ceil
				serverTimeChange = float64((int64(serverTimeChange) / 5) * 5) // floor to nearest 5

				if serverTimeChange > 0 {
					possibleMultiplier := int(serverTimeChange/5 + 0.5) // ceil

					// Validar rango del multiplicador
					if possibleMultiplier > 0 && possibleMultiplier < 25 {
						if st.timeMultiplier != possibleMultiplier {
							st.timeMultiplier = possibleMultiplier

							// Notificar callbacks
							for _, callback := range st.multiplierCallbacks {
								go callback(possibleMultiplier)
							}
						}
					}
				}
			}
		}

		st.lastServerMillis = newMilliseconds
		st.serverChanges++
	}
}

// Reset reinicia los datos de seguimiento de tiempo
func (st *SessionTimeTracker) Reset() {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.serverChanges = 0
	st.lastServerMillis = -1
	st.timeMultiplier = -1
}

// GetTimeMultiplier devuelve el multiplicador de tiempo actual
func (st *SessionTimeTracker) GetTimeMultiplier() int {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.timeMultiplier
}

// GetFormattedTimeRemaining formatea el tiempo restante en formato MM:SS
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	if minutes >= 60 {
		hours := minutes / 60
		minutes = minutes % 60
		return formatTime(hours, minutes, seconds, true)
	}

	return formatTime(0, minutes, seconds, false)
}

func formatTime(hours, minutes, seconds int, includeHours bool) string {
	if includeHours {
		return formatHMS(hours, minutes, seconds)
	}
	return formatMS(minutes, seconds)
}

func formatHMS(hours, minutes, seconds int) string {
	return pad2(hours) + ":" + pad2(minutes) + ":" + pad2(seconds)
}

func formatMS(minutes, seconds int) string {
	return pad2(minutes) + ":" + pad2(seconds)
}

func pad2(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}
