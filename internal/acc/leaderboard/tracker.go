package leaderboard

import (
	"RaceAll/internal/broadcast"
	"fmt"
	"sort"
	"time"
)

// DriverPosition representa la posición de un piloto
type DriverPosition struct {
	Position       int
	CarIndex       uint16
	CarNumber      int32
	DriverName     string
	TeamName       string
	CarModel       byte
	BestLap        int32 // ms
	LastLap        int32 // ms
	CurrentLap     int32 // ms
	Laps           uint16
	Gap            string // Gap al líder
	Interval       string // Gap al anterior
	GapMS          int32  // Gap en ms
	IntervalMS     int32  // Interval en ms
	SplinePosition float32
	IsPlayer       bool
	IsInPit        bool
	CupPosition    int
	Category       byte
	Location       broadcast.CarLocationEnum
	Speed          uint16 // km/h
}

// LeaderboardTracker gestiona el leaderboard
type LeaderboardTracker struct {
	positions      []DriverPosition
	playerCarIndex uint16
	sessionType    broadcast.RaceSessionType
	leaderBestLap  int32
	lastUpdateTime time.Time
}

// NewLeaderboardTracker crea un nuevo tracker de leaderboard
func NewLeaderboardTracker() *LeaderboardTracker {
	return &LeaderboardTracker{
		positions:      make([]DriverPosition, 0),
		playerCarIndex: 0,
		leaderBestLap:  0,
		lastUpdateTime: time.Now(),
	}
}

// Update actualiza el leaderboard con datos del broadcast
func (lt *LeaderboardTracker) Update(
	cars map[uint16]*broadcast.CarInfo,
	updates map[uint16]*broadcast.RealtimeCarUpdate,
	sessionType broadcast.RaceSessionType,
	playerCarIndex uint16,
) {
	lt.sessionType = sessionType
	lt.playerCarIndex = playerCarIndex
	lt.positions = make([]DriverPosition, 0)

	// Construir lista de posiciones
	for carIndex, carInfo := range cars {
		update, hasUpdate := updates[carIndex]
		if !hasUpdate {
			continue
		}

		position := DriverPosition{
			Position:       int(update.Position),
			CarIndex:       carIndex,
			CarNumber:      carInfo.RaceNumber,
			DriverName:     carInfo.GetCurrentDriverName(),
			TeamName:       carInfo.TeamName,
			CarModel:       carInfo.CarModelType,
			BestLap:        update.BestSessionLap.GetLapTimeMS(),
			LastLap:        update.LastLap.GetLapTimeMS(),
			CurrentLap:     0, // Se calculará con tiempo actual
			Laps:           update.Laps,
			SplinePosition: update.SplinePosition,
			IsPlayer:       carIndex == playerCarIndex,
			IsInPit:        update.CarLocation == broadcast.CarLocationPitlane,
			CupPosition:    int(update.CupPosition),
			Category:       carInfo.CupCategory,
			Location:       update.CarLocation,
			Speed:          update.Kmh,
		}

		lt.positions = append(lt.positions, position)
	}

	// Ordenar por posición
	sort.Slice(lt.positions, func(i, j int) bool {
		return lt.positions[i].Position < lt.positions[j].Position
	})

	// Calcular gaps e intervals
	lt.calculateGaps()

	lt.lastUpdateTime = time.Now()
}

// calculateGaps calcula gaps e intervalos
func (lt *LeaderboardTracker) calculateGaps() {
	if len(lt.positions) == 0 {
		return
	}

	// El líder tiene mejor vuelta de referencia
	if lt.positions[0].BestLap > 0 {
		lt.leaderBestLap = lt.positions[0].BestLap
	}

	for i := range lt.positions {
		if i == 0 {
			// Líder
			lt.positions[i].Gap = "Leader"
			lt.positions[i].Interval = "---"
			lt.positions[i].GapMS = 0
			lt.positions[i].IntervalMS = 0
		} else {
			// Calcular gap basado en vueltas y posición
			leaderPos := lt.positions[0]
			currentPos := lt.positions[i]
			previousPos := lt.positions[i-1]

			// Gap en vueltas
			lapDiff := int(leaderPos.Laps) - int(currentPos.Laps)

			if lapDiff > 0 {
				// Diferencia de vueltas
				lt.positions[i].Gap = fmt.Sprintf("+%d Lap", lapDiff)
				if lapDiff > 1 {
					lt.positions[i].Gap = fmt.Sprintf("+%d Laps", lapDiff)
				}
				lt.positions[i].GapMS = int32(lapDiff) * lt.leaderBestLap
			} else {
				// Misma vuelta - calcular por spline position
				splineDiff := leaderPos.SplinePosition - currentPos.SplinePosition
				if splineDiff < 0 {
					splineDiff += 1.0
				}

				// Estimar tiempo basado en mejor vuelta
				estimatedGapMS := int32(splineDiff * float32(lt.leaderBestLap))
				lt.positions[i].GapMS = estimatedGapMS
				lt.positions[i].Gap = formatTimeGap(estimatedGapMS)
			}

			// Interval al anterior
			prevLapDiff := int(previousPos.Laps) - int(currentPos.Laps)

			if prevLapDiff > 0 {
				lt.positions[i].Interval = fmt.Sprintf("+%d L", prevLapDiff)
				lt.positions[i].IntervalMS = int32(prevLapDiff) * lt.leaderBestLap
			} else {
				splineDiff := previousPos.SplinePosition - currentPos.SplinePosition
				if splineDiff < 0 {
					splineDiff += 1.0
				}

				estimatedIntervalMS := int32(splineDiff * float32(lt.leaderBestLap))
				lt.positions[i].IntervalMS = estimatedIntervalMS
				lt.positions[i].Interval = formatTimeGap(estimatedIntervalMS)
			}
		}
	}
}

// formatTimeGap formatea un gap de tiempo en ms a string
func formatTimeGap(ms int32) string {
	if ms <= 0 {
		return "---"
	}

	seconds := float64(ms) / 1000.0

	if seconds < 60 {
		return fmt.Sprintf("+%.3f", seconds)
	}

	minutes := int(seconds / 60)
	secs := seconds - float64(minutes*60)
	return fmt.Sprintf("+%d:%06.3f", minutes, secs)
}

// GetPositions devuelve todas las posiciones
func (lt *LeaderboardTracker) GetPositions() []DriverPosition {
	return lt.positions
}

// GetPlayerPosition devuelve la posición del jugador
func (lt *LeaderboardTracker) GetPlayerPosition() *DriverPosition {
	for i := range lt.positions {
		if lt.positions[i].IsPlayer {
			return &lt.positions[i]
		}
	}
	return nil
}

// GetRelativePositions devuelve posiciones relativas al jugador (±5 posiciones)
func (lt *LeaderboardTracker) GetRelativePositions(range_ int) []DriverPosition {
	playerPos := lt.GetPlayerPosition()
	if playerPos == nil {
		return lt.positions
	}

	relative := make([]DriverPosition, 0)
	playerIdx := playerPos.Position - 1

	startIdx := playerIdx - range_
	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := playerIdx + range_ + 1
	if endIdx > len(lt.positions) {
		endIdx = len(lt.positions)
	}

	for i := startIdx; i < endIdx; i++ {
		relative = append(relative, lt.positions[i])
	}

	return relative
}

// GetTopN devuelve las primeras N posiciones
func (lt *LeaderboardTracker) GetTopN(n int) []DriverPosition {
	if n > len(lt.positions) {
		n = len(lt.positions)
	}
	return lt.positions[:n]
}

// GetPositionByCarIndex devuelve la posición de un auto específico
func (lt *LeaderboardTracker) GetPositionByCarIndex(carIndex uint16) *DriverPosition {
	for i := range lt.positions {
		if lt.positions[i].CarIndex == carIndex {
			return &lt.positions[i]
		}
	}
	return nil
}

// GetLeader devuelve el líder
func (lt *LeaderboardTracker) GetLeader() *DriverPosition {
	if len(lt.positions) > 0 {
		return &lt.positions[0]
	}
	return nil
}

// GetTotalCars devuelve el número total de autos
func (lt *LeaderboardTracker) GetTotalCars() int {
	return len(lt.positions)
}

// GetClassPositions devuelve posiciones filtradas por clase
func (lt *LeaderboardTracker) GetClassPositions(category byte) []DriverPosition {
	filtered := make([]DriverPosition, 0)

	for _, pos := range lt.positions {
		if pos.Category == category {
			filtered = append(filtered, pos)
		}
	}

	// Reordenar posiciones de clase
	for i := range filtered {
		filtered[i].CupPosition = i + 1
	}

	return filtered
}

// IsBattleFor verifica si hay batalla por una posición específica
func (lt *LeaderboardTracker) IsBattleFor(position int, gapThreshold float32) bool {
	if position <= 0 || position >= len(lt.positions) {
		return false
	}

	currentPos := lt.positions[position]

	// Batalla si el gap es menor al threshold (en segundos)
	gapSeconds := float32(currentPos.IntervalMS) / 1000.0
	return gapSeconds < gapThreshold && gapSeconds > 0
}

// GetBattles devuelve lista de batallas (gaps < 2 segundos)
func (lt *LeaderboardTracker) GetBattles() [][]DriverPosition {
	battles := make([][]DriverPosition, 0)
	currentBattle := make([]DriverPosition, 0)

	for i, pos := range lt.positions {
		if i == 0 {
			currentBattle = append(currentBattle, pos)
			continue
		}

		gapSeconds := float32(pos.IntervalMS) / 1000.0

		if gapSeconds < 2.0 && gapSeconds > 0 {
			if len(currentBattle) == 0 {
				currentBattle = append(currentBattle, lt.positions[i-1])
			}
			currentBattle = append(currentBattle, pos)
		} else {
			if len(currentBattle) >= 2 {
				battles = append(battles, currentBattle)
			}
			currentBattle = make([]DriverPosition, 0)
		}
	}

	// Agregar última batalla si existe
	if len(currentBattle) >= 2 {
		battles = append(battles, currentBattle)
	}

	return battles
}

// GetFastestLap devuelve la mejor vuelta de la sesión
func (lt *LeaderboardTracker) GetFastestLap() (int32, string) {
	fastestTime := int32(0)
	driverName := ""

	for _, pos := range lt.positions {
		if pos.BestLap > 0 {
			if fastestTime == 0 || pos.BestLap < fastestTime {
				fastestTime = pos.BestLap
				driverName = pos.DriverName
			}
		}
	}

	return fastestTime, driverName
}

// SortByBestLap ordena por mejor vuelta (para qualifying)
func (lt *LeaderboardTracker) SortByBestLap() {
	sort.Slice(lt.positions, func(i, j int) bool {
		// Autos sin tiempo al final
		if lt.positions[i].BestLap == 0 {
			return false
		}
		if lt.positions[j].BestLap == 0 {
			return true
		}
		return lt.positions[i].BestLap < lt.positions[j].BestLap
	})

	// Actualizar posiciones
	for i := range lt.positions {
		lt.positions[i].Position = i + 1
	}
}
