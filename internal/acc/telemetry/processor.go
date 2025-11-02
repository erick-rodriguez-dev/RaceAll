package telemetry

import (
	"RaceAll/internal/sharedmemory"
	"math"
)

// TelemetryData representa datos de telemetría procesados
type TelemetryData struct {
	// Inputs
	Throttle     float32 // 0.0 - 1.0
	Brake        float32 // 0.0 - 1.0
	Clutch       float32 // 0.0 - 1.0
	Steering     float32 // Ángulo en radianes
	SteeringDeg  float32 // Ángulo en grados
	SteeringLock float32 // Porcentaje de lock usado

	// Engine
	Gear       int32
	RPM        int32
	MaxRPM     int32
	RPMPercent float32 // 0.0 - 100.0
	Speed      float32 // km/h
	SpeedMPH   float32 // mph

	// G-Forces
	GForceLateral      float32
	GForceLongitudinal float32
	GForceVertical     float32
	GForceTotal        float32

	// Wheels
	WheelSlip    [4]float32
	WheelLoad    [4]float32
	WheelLock    [4]bool
	IsLockingAny bool

	// Electronics
	TC           float32 // 0.0 - 1.0
	TCActive     bool
	ABS          float32 // 0.0 - 1.0
	ABSActive    bool
	DRSAvailable bool
	DRSEnabled   bool

	// Position & Orientation
	Heading float32
	Pitch   float32
	Roll    float32

	// Flags
	PitLimiter     bool
	AutoShifter    bool
	IsAIControlled bool
	EngineRunning  bool
}

// TelemetryProcessor procesa datos de telemetría
type TelemetryProcessor struct {
	lastThrottle     float32
	lastBrake        float32
	throttleSmoothed float32
	brakeSmoothed    float32
	smoothingFactor  float32
}

// NewTelemetryProcessor crea un nuevo procesador de telemetría
func NewTelemetryProcessor() *TelemetryProcessor {
	return &TelemetryProcessor{
		lastThrottle:     0,
		lastBrake:        0,
		throttleSmoothed: 0,
		brakeSmoothed:    0,
		smoothingFactor:  0.2, // Factor de suavizado exponencial
	}
}

// ProcessPhysics procesa datos de física de shared memory
func (tp *TelemetryProcessor) ProcessPhysics(physics *sharedmemory.Physics) TelemetryData {
	// Suavizar throttle y brake
	tp.throttleSmoothed = tp.smooth(tp.throttleSmoothed, physics.Gas, tp.smoothingFactor)
	tp.brakeSmoothed = tp.smooth(tp.brakeSmoothed, physics.Brake, tp.smoothingFactor)

	// Calcular ángulo de dirección en grados
	steeringDeg := physics.SteerAngle * (180.0 / math.Pi)

	// Calcular porcentaje de steering lock (asumiendo ~540° de lock to lock)
	maxSteeringAngle := float32(270.0)
	steeringLock := (float32(math.Abs(float64(steeringDeg))) / maxSteeringAngle) * 100.0
	if steeringLock > 100.0 {
		steeringLock = 100.0
	}

	// Calcular G-Force total
	gForceTotal := float32(math.Sqrt(
		float64(physics.AccG[0]*physics.AccG[0] +
			physics.AccG[1]*physics.AccG[1] +
			physics.AccG[2]*physics.AccG[2])))

	// Detectar wheel lock (slip ratio muy negativo)
	wheelLock := [4]bool{
		physics.SlipRatio[0] < -0.1,
		physics.SlipRatio[1] < -0.1,
		physics.SlipRatio[2] < -0.1,
		physics.SlipRatio[3] < -0.1,
	}

	isLockingAny := wheelLock[0] || wheelLock[1] || wheelLock[2] || wheelLock[3]

	// Calcular porcentaje de RPM
	rpmPercent := float32(0)
	if physics.CurrentMaxRpm > 0 {
		rpmPercent = (float32(physics.Rpms) / float32(physics.CurrentMaxRpm)) * 100.0
	}

	return TelemetryData{
		// Inputs
		Throttle:     tp.throttleSmoothed,
		Brake:        tp.brakeSmoothed,
		Clutch:       physics.Clutch,
		Steering:     physics.SteerAngle,
		SteeringDeg:  steeringDeg,
		SteeringLock: steeringLock,

		// Engine
		Gear:       physics.Gear,
		RPM:        physics.Rpms,
		MaxRPM:     physics.CurrentMaxRpm,
		RPMPercent: rpmPercent,
		Speed:      physics.SpeedKmh,
		SpeedMPH:   physics.SpeedKmh * 0.621371,

		// G-Forces
		GForceLateral:      physics.AccG[0],
		GForceLongitudinal: physics.AccG[1],
		GForceVertical:     physics.AccG[2],
		GForceTotal:        gForceTotal,

		// Wheels
		WheelSlip:    [4]float32{physics.WheelSlip[0], physics.WheelSlip[1], physics.WheelSlip[2], physics.WheelSlip[3]},
		WheelLoad:    [4]float32{physics.WheelLoad[0], physics.WheelLoad[1], physics.WheelLoad[2], physics.WheelLoad[3]},
		WheelLock:    wheelLock,
		IsLockingAny: isLockingAny,

		// Electronics
		TC:           physics.TC,
		TCActive:     physics.TcinAction == 1,
		ABS:          physics.Abs,
		ABSActive:    physics.AbsInAction == 1,
		DRSAvailable: physics.DrsAvailable == 1,
		DRSEnabled:   physics.DrsEnabled == 1,

		// Position & Orientation
		Heading: physics.Heading,
		Pitch:   physics.Pitch,
		Roll:    physics.Roll,

		// Flags
		PitLimiter:     physics.PitLimiterOn == 1,
		AutoShifter:    physics.AutoShifterOn == 1,
		IsAIControlled: physics.IsAIControlled == 1,
		EngineRunning:  physics.IsEngineRunning == 1,
	}
}

// smooth aplica suavizado exponencial
func (tp *TelemetryProcessor) smooth(current, target, factor float32) float32 {
	return current + (target-current)*factor
}

// IsThrottleIncreasing verifica si el throttle está aumentando
func (tp *TelemetryProcessor) IsThrottleIncreasing() bool {
	result := tp.throttleSmoothed > tp.lastThrottle
	tp.lastThrottle = tp.throttleSmoothed
	return result
}

// IsBraking verifica si está frenando
func (tp *TelemetryProcessor) IsBraking() bool {
	return tp.brakeSmoothed > 0.1
}

// IsFullThrottle verifica si está a fondo
func (td *TelemetryData) IsFullThrottle() bool {
	return td.Throttle > 0.95
}

// IsCoasting verifica si está en coast (sin throttle ni brake)
func (td *TelemetryData) IsCoasting() bool {
	return td.Throttle < 0.05 && td.Brake < 0.05
}

// IsTrailBraking verifica si está haciendo trail braking
func (td *TelemetryData) IsTrailBraking() bool {
	return td.Brake > 0.1 && td.Throttle > 0.1
}

// GetGearRatio calcula el ratio de marcha aproximado
func (td *TelemetryData) GetGearRatio() float32 {
	if td.Speed <= 0 || td.RPM <= 0 {
		return 0
	}
	return float32(td.RPM) / td.Speed
}

// ShouldUpshift verifica si debería cambiar a marcha superior
func (td *TelemetryData) ShouldUpshift() bool {
	return td.RPMPercent > 95.0 && td.Throttle > 0.9
}

// ShouldDownshift verifica si debería cambiar a marcha inferior
func (td *TelemetryData) ShouldDownshift() bool {
	return td.RPMPercent < 40.0 && td.Gear > 1
}

// GetCorneringForce calcula la fuerza lateral en curvas
func (td *TelemetryData) GetCorneringForce() float32 {
	return float32(math.Abs(float64(td.GForceLateral)))
}

// IsCornering verifica si está en curva
func (td *TelemetryData) IsCornering() bool {
	return math.Abs(float64(td.SteeringDeg)) > 5.0
}

// GetTotalWheelSlip calcula el slip total de las ruedas
func (td *TelemetryData) GetTotalWheelSlip() float32 {
	var total float32
	for _, slip := range td.WheelSlip {
		total += float32(math.Abs(float64(slip)))
	}
	return total / 4.0
}

// HasTractionIssue verifica si hay problemas de tracción
func (td *TelemetryData) HasTractionIssue() bool {
	return td.GetTotalWheelSlip() > 0.3
}

// IsUndersteering verifica si hay subviraje
func (td *TelemetryData) IsUndersteering() bool {
	// Subviraje: steering input alto pero poca fuerza lateral
	return math.Abs(float64(td.SteeringDeg)) > 30.0 && td.GetCorneringForce() < 1.0
}

// IsOversteering verifica si hay sobreviraje
func (td *TelemetryData) IsOversteering() bool {
	// Sobreviraje: alta fuerza lateral con las ruedas traseras deslizando
	rearSlip := (td.WheelSlip[2] + td.WheelSlip[3]) / 2.0
	return td.GetCorneringForce() > 1.5 && rearSlip > 0.4
}

// GetDrivingStyle analiza el estilo de conducción
func (td *TelemetryData) GetDrivingStyle() string {
	if td.IsFullThrottle() && !td.IsCornering() {
		return "Straight Line"
	}
	if td.IsTrailBraking() {
		return "Trail Braking"
	}
	if td.Brake > 0.1 && td.IsCornering() {
		return "Braking in Corner"
	}
	if td.IsCornering() && td.Throttle > 0.5 {
		return "Accelerating in Corner"
	}
	if td.IsCoasting() {
		return "Coasting"
	}
	return "Normal"
}

// GetInputSmoothness calcula la suavidad de inputs (0-100, mayor es mejor)
func (td *TelemetryData) GetInputSmoothness(previousThrottle, previousBrake float32) float32 {
	throttleDiff := float32(math.Abs(float64(td.Throttle - previousThrottle)))
	brakeDiff := float32(math.Abs(float64(td.Brake - previousBrake)))

	avgDiff := (throttleDiff + brakeDiff) / 2.0
	smoothness := (1.0 - avgDiff) * 100.0

	if smoothness < 0 {
		smoothness = 0
	}
	if smoothness > 100 {
		smoothness = 100
	}

	return smoothness
}
