package acc

import (
	"RaceAll/internal/acc/cars"
	"RaceAll/internal/acc/entrylist"
	"RaceAll/internal/acc/fuel"
	"RaceAll/internal/acc/gaps"
	"RaceAll/internal/acc/incidents"
	"RaceAll/internal/acc/laps"
	"RaceAll/internal/acc/leaderboard"
	"RaceAll/internal/acc/session"
	"RaceAll/internal/acc/sessiontime"
	"RaceAll/internal/acc/telemetry"
	"RaceAll/internal/acc/trackposition"
	"RaceAll/internal/acc/tracks"
	"RaceAll/internal/acc/tyres"
	"RaceAll/internal/broadcast"
	"RaceAll/internal/sharedmemory"
)

// DataManager gestiona todos los módulos de análisis de ACC
type DataManager struct {
	// Trackers y procesadores
	sessionTracker   *session.SessionTracker
	lapTracker       *laps.LapTracker
	fuelCalculator   *fuel.FuelCalculator
	tyresTracker     *tyres.TyresTracker
	telemetryProc    *telemetry.TelemetryProcessor
	leaderboard      *leaderboard.LeaderboardTracker
	gapTracker       *gaps.GapTracker
	positionGraph    *trackposition.PositionGraph
	incidentTracker  *incidents.IncidentTracker
	sessionTimer     *sessiontime.SessionTimeTracker
	entryListTracker *entrylist.EntryListTracker

	// Información del auto
	carModel cars.CarModel
	carIndex uint16
	trackID  tracks.TrackID

	// Estado
	initialized bool
}

// NewDataManager crea un nuevo gestor de datos ACC
func NewDataManager() *DataManager {
	return &DataManager{
		sessionTracker:   session.NewSessionTracker(),
		leaderboard:      leaderboard.NewLeaderboardTracker(),
		telemetryProc:    telemetry.NewTelemetryProcessor(),
		gapTracker:       gaps.NewGapTracker(),
		positionGraph:    trackposition.NewPositionGraph(),
		incidentTracker:  incidents.NewIncidentTracker(),
		sessionTimer:     sessiontime.NewSessionTimeTracker(),
		entryListTracker: entrylist.NewEntryListTracker(),
		initialized:      false,
	}
}

// Initialize inicializa el data manager con información del auto
func (dm *DataManager) Initialize(carModel cars.CarModel, carIndex uint16, trackID tracks.TrackID) {
	dm.carModel = carModel
	dm.carIndex = carIndex
	dm.trackID = trackID

	dm.lapTracker = laps.NewLapTracker(carIndex)
	dm.fuelCalculator = fuel.NewFuelCalculator(carModel)
	dm.tyresTracker = tyres.NewTyresTracker()

	// Inicializar gap tracker con la distancia del circuito
	trackInfo := tracks.GetTrackInfo(trackID)
	dm.gapTracker.Initialize(float32(trackInfo.LengthMeters))

	dm.initialized = true
}

// UpdateFromBroadcast actualiza con datos del broadcast
func (dm *DataManager) UpdateFromBroadcast(
	realtimeUpdate *broadcast.RealtimeUpdate,
	carUpdate *broadcast.RealtimeCarUpdate,
	allCars map[uint16]*broadcast.CarInfo,
	allUpdates map[uint16]*broadcast.RealtimeCarUpdate,
) {
	if !dm.initialized {
		return
	}

	// Actualizar sesión
	dm.sessionTracker.Update(realtimeUpdate)

	// Actualizar session timer
	dm.sessionTimer.Update(realtimeUpdate.TimeOfDay)

	// Actualizar entry list
	if allCars != nil {
		for _, carInfo := range allCars {
			dm.entryListTracker.UpdateCarInfo(carInfo)
		}
	}
	if allUpdates != nil {
		for _, update := range allUpdates {
			dm.entryListTracker.UpdateRealtimeCarUpdate(update)

			// Actualizar gap tracker
			dm.gapTracker.UpdateCarPosition(update.CarIndex, update.SplinePosition)

			// Actualizar position graph
			dm.positionGraph.UpdateLocation(update.CarIndex, update.SplinePosition, update.CarLocation)

			// Actualizar incident tracker con datos en tiempo real
			dm.incidentTracker.UpdateRealtimeCarUpdate(update, realtimeUpdate.SessionTime)
		}
	}

	// Actualizar vueltas si se completó una
	if carUpdate != nil && carUpdate.LastLap.HasValidLapTime() {
		dm.lapTracker.UpdateFromBroadcast(&carUpdate.LastLap)
	}

	// Actualizar leaderboard
	if allCars != nil && allUpdates != nil {
		dm.leaderboard.Update(
			allCars,
			allUpdates,
			realtimeUpdate.SessionType,
			dm.carIndex,
		)
	}
}

// HandleBroadcastEvent maneja eventos de broadcast (como accidentes)
func (dm *DataManager) HandleBroadcastEvent(event *broadcast.BroadcastingEvent, carInfo *broadcast.CarInfo) {
	if !dm.initialized {
		return
	}

	// Procesar eventos según tipo
	switch event.Type {
	case broadcast.BroadcastingEventTypeAccident:
		dm.incidentTracker.HandleAccidentEvent(event, carInfo)
	}
}

// UpdateTrackData actualiza información del circuito
func (dm *DataManager) UpdateTrackData(trackData *broadcast.TrackData) {
	if !dm.initialized {
		return
	}

	dm.incidentTracker.UpdateTrackData(float32(trackData.TrackMeters))
}

// UpdateFromSharedMemory actualiza con datos de shared memory
func (dm *DataManager) UpdateFromSharedMemory(
	physics *sharedmemory.Physics,
	graphics *sharedmemory.Graphics,
	static *sharedmemory.Static,
) {
	if !dm.initialized {
		return
	}

	// Procesar telemetría
	_ = dm.telemetryProc.ProcessPhysics(physics)

	// Actualizar combustible (detectar vuelta completada)
	lapCompleted := false // Se puede determinar con lógica adicional
	_ = dm.fuelCalculator.Update(physics.Fuel, lapCompleted)

	// Actualizar neumáticos
	dm.tyresTracker.Update(
		physics.WheelsPressure,
		physics.TyreTemp,
		physics.TyreTempI,
		physics.TyreTempM,
		physics.TyreTempO,
		physics.TyreCoreTemperature,
		physics.TyreWear,
		physics.TyreDirtyLevel,
		physics.BrakeTemp,
	)
}

// GetSessionData devuelve datos de la sesión
func (dm *DataManager) GetSessionData() *session.SessionState {
	return dm.sessionTracker.GetCurrentState()
}

// GetLapData devuelve datos de vueltas
func (dm *DataManager) GetLapData() *laps.LapData {
	return dm.lapTracker.GetBestLap()
}

// GetFuelData devuelve datos de combustible
func (dm *DataManager) GetFuelData(currentFuel float32) fuel.FuelData {
	return dm.fuelCalculator.Update(currentFuel, false)
}

// GetTyresData devuelve datos de neumáticos
func (dm *DataManager) GetTyresData() [4]tyres.TyreData {
	return dm.tyresTracker.GetAllTyres()
}

// GetLeaderboardData devuelve datos del leaderboard
func (dm *DataManager) GetLeaderboardData() []leaderboard.DriverPosition {
	return dm.leaderboard.GetPositions()
}

// GetCarInfo devuelve información del auto actual
func (dm *DataManager) GetCarInfo() cars.CarInfo {
	return cars.GetCarInfo(dm.carModel)
}

// GetTrackInfo devuelve información del circuito actual
func (dm *DataManager) GetTrackInfo() tracks.TrackInfo {
	return tracks.GetTrackInfo(dm.trackID)
}

// Reset reinicia todos los trackers
func (dm *DataManager) Reset() {
	if dm.lapTracker != nil {
		dm.lapTracker = laps.NewLapTracker(dm.carIndex)
	}
	if dm.fuelCalculator != nil {
		dm.fuelCalculator.Reset()
	}
	if dm.tyresTracker != nil {
		dm.tyresTracker.Reset()
	}
	dm.gapTracker.Reset()
	dm.positionGraph.Reset()
	dm.incidentTracker.Clear()
	dm.sessionTimer.Reset()
	dm.entryListTracker.Clear()
	dm.initialized = false
}

// IsInitialized verifica si el manager está inicializado
func (dm *DataManager) IsInitialized() bool {
	return dm.initialized
}

// GetSessionTracker devuelve el tracker de sesión
func (dm *DataManager) GetSessionTracker() *session.SessionTracker {
	return dm.sessionTracker
}

// GetLapTracker devuelve el tracker de vueltas
func (dm *DataManager) GetLapTracker() *laps.LapTracker {
	return dm.lapTracker
}

// GetFuelCalculator devuelve el calculador de combustible
func (dm *DataManager) GetFuelCalculator() *fuel.FuelCalculator {
	return dm.fuelCalculator
}

// GetTyresTracker devuelve el tracker de neumáticos
func (dm *DataManager) GetTyresTracker() *tyres.TyresTracker {
	return dm.tyresTracker
}

// GetTelemetryProcessor devuelve el procesador de telemetría
func (dm *DataManager) GetTelemetryProcessor() *telemetry.TelemetryProcessor {
	return dm.telemetryProc
}

// GetLeaderboardTracker devuelve el tracker de leaderboard
func (dm *DataManager) GetLeaderboardTracker() *leaderboard.LeaderboardTracker {
	return dm.leaderboard
}

// GetGapTracker devuelve el tracker de gaps
func (dm *DataManager) GetGapTracker() *gaps.GapTracker {
	return dm.gapTracker
}

// GetPositionGraph devuelve el grafo de posiciones
func (dm *DataManager) GetPositionGraph() *trackposition.PositionGraph {
	return dm.positionGraph
}

// GetIncidentTracker devuelve el tracker de incidentes
func (dm *DataManager) GetIncidentTracker() *incidents.IncidentTracker {
	return dm.incidentTracker
}

// GetSessionTimeTracker devuelve el tracker de tiempo de sesión
func (dm *DataManager) GetSessionTimeTracker() *sessiontime.SessionTimeTracker {
	return dm.sessionTimer
}

// GetEntryListTracker devuelve el tracker de lista de entrada
func (dm *DataManager) GetEntryListTracker() *entrylist.EntryListTracker {
	return dm.entryListTracker
}
