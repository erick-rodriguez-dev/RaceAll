package tyres

import (
	"math"
)

// TyrePosition representa la posición del neumático
type TyrePosition int

const (
	FrontLeft TyrePosition = iota
	FrontRight
	RearLeft
	RearRight
)

// TyreCompound representa el compuesto de neumático
type TyreCompound byte

const (
	CompoundDry TyreCompound = iota
	CompoundWet
)

// TyreData representa los datos de un neumático individual
type TyreData struct {
	Position         TyrePosition
	Pressure         float32 // PSI
	Temperature      float32 // °C (promedio)
	TempInner        float32 // °C
	TempMiddle       float32 // °C
	TempOuter        float32 // °C
	CoreTemperature  float32 // °C
	Wear             float32 // 0.0 - 1.0
	DirtyLevel       float32 // 0.0 - 1.0
	BrakeTemperature float32 // °C
	IsOptimalTemp    bool
	IsOverheating    bool
	IsUnderheating   bool
	IsCold           bool
}

// TyresTracker realiza seguimiento de los neumáticos
type TyresTracker struct {
	tyres          [4]TyreData
	compound       TyreCompound
	optimalTempMin float32
	optimalTempMax float32
	wearHistory    [][4]float32
}

// NewTyresTracker crea un nuevo tracker de neumáticos
func NewTyresTracker() *TyresTracker {
	return &TyresTracker{
		tyres:          [4]TyreData{},
		compound:       CompoundDry,
		optimalTempMin: 75.0, // Temperatura óptima mínima para slicks
		optimalTempMax: 95.0, // Temperatura óptima máxima para slicks
		wearHistory:    make([][4]float32, 0),
	}
}

// Update actualiza los datos de los neumáticos
func (tt *TyresTracker) Update(
	pressures [4]float32,
	temps [4]float32,
	tempInner [4]float32,
	tempMiddle [4]float32,
	tempOuter [4]float32,
	coreTemps [4]float32,
	wear [4]float32,
	dirtyLevel [4]float32,
	brakeTemps [4]float32,
) {
	positions := []TyrePosition{FrontLeft, FrontRight, RearLeft, RearRight}

	for i, pos := range positions {
		avgTemp := (tempInner[i] + tempMiddle[i] + tempOuter[i]) / 3.0

		tt.tyres[i] = TyreData{
			Position:         pos,
			Pressure:         pressures[i],
			Temperature:      avgTemp,
			TempInner:        tempInner[i],
			TempMiddle:       tempMiddle[i],
			TempOuter:        tempOuter[i],
			CoreTemperature:  coreTemps[i],
			Wear:             wear[i],
			DirtyLevel:       dirtyLevel[i],
			BrakeTemperature: brakeTemps[i],
			IsOptimalTemp:    avgTemp >= tt.optimalTempMin && avgTemp <= tt.optimalTempMax,
			IsOverheating:    avgTemp > tt.optimalTempMax+10.0,
			IsUnderheating:   avgTemp < tt.optimalTempMin-10.0,
			IsCold:           avgTemp < 50.0,
		}
	}

	// Guardar historial de desgaste
	currentWear := [4]float32{wear[0], wear[1], wear[2], wear[3]}
	tt.wearHistory = append(tt.wearHistory, currentWear)

	// Mantener solo últimas 100 muestras
	if len(tt.wearHistory) > 100 {
		tt.wearHistory = tt.wearHistory[1:]
	}
}

// GetTyre devuelve los datos de un neumático específico
func (tt *TyresTracker) GetTyre(position TyrePosition) TyreData {
	return tt.tyres[position]
}

// GetAllTyres devuelve los datos de todos los neumáticos
func (tt *TyresTracker) GetAllTyres() [4]TyreData {
	return tt.tyres
}

// GetAverageWear devuelve el desgaste promedio de todos los neumáticos
func (tt *TyresTracker) GetAverageWear() float32 {
	var total float32
	for _, tyre := range tt.tyres {
		total += tyre.Wear
	}
	return total / 4.0
}

// GetAverageTemperature devuelve la temperatura promedio de todos los neumáticos
func (tt *TyresTracker) GetAverageTemperature() float32 {
	var total float32
	for _, tyre := range tt.tyres {
		total += tyre.Temperature
	}
	return total / 4.0
}

// GetAveragePressure devuelve la presión promedio de todos los neumáticos
func (tt *TyresTracker) GetAveragePressure() float32 {
	var total float32
	for _, tyre := range tt.tyres {
		total += tyre.Pressure
	}
	return total / 4.0
}

// AreAllTyresOptimal verifica si todos los neumáticos están en temperatura óptima
func (tt *TyresTracker) AreAllTyresOptimal() bool {
	for _, tyre := range tt.tyres {
		if !tyre.IsOptimalTemp {
			return false
		}
	}
	return true
}

// HasOverheatingTyres verifica si algún neumático está sobrecalentando
func (tt *TyresTracker) HasOverheatingTyres() bool {
	for _, tyre := range tt.tyres {
		if tyre.IsOverheating {
			return true
		}
	}
	return false
}

// HasColdTyres verifica si algún neumático está frío
func (tt *TyresTracker) HasColdTyres() bool {
	for _, tyre := range tt.tyres {
		if tyre.IsCold {
			return true
		}
	}
	return false
}

// GetWearRate calcula la tasa de desgaste por vuelta
func (tt *TyresTracker) GetWearRate() float32 {
	if len(tt.wearHistory) < 2 {
		return 0
	}

	// Comparar desgaste actual vs inicial
	first := tt.wearHistory[0]
	last := tt.wearHistory[len(tt.wearHistory)-1]

	var totalWearDiff float32
	for i := 0; i < 4; i++ {
		totalWearDiff += last[i] - first[i]
	}

	avgWearDiff := totalWearDiff / 4.0
	laps := float32(len(tt.wearHistory))

	return avgWearDiff / laps
}

// EstimateLapsRemaining estima las vueltas restantes antes de cambio de neumáticos
func (tt *TyresTracker) EstimateLapsRemaining() float32 {
	wearRate := tt.GetWearRate()
	if wearRate <= 0 {
		return 999 // Infinito
	}

	currentWear := tt.GetAverageWear()
	remainingWear := 1.0 - currentWear

	return remainingWear / wearRate
}

// ShouldChangeTyres recomienda si cambiar neumáticos
func (tt *TyresTracker) ShouldChangeTyres() bool {
	avgWear := tt.GetAverageWear()

	// Cambiar si el desgaste promedio supera 80%
	if avgWear > 0.8 {
		return true
	}

	// Cambiar si algún neumático tiene desgaste crítico (>90%)
	for _, tyre := range tt.tyres {
		if tyre.Wear > 0.9 {
			return true
		}
	}

	return false
}

// GetTempBalance calcula el balance de temperatura entre neumáticos
func (tt *TyresTracker) GetTempBalance() (frontBalance, rearBalance float32) {
	frontLeft := tt.tyres[FrontLeft].Temperature
	frontRight := tt.tyres[FrontRight].Temperature
	rearLeft := tt.tyres[RearLeft].Temperature
	rearRight := tt.tyres[RearRight].Temperature

	// Balance lateral (diferencia izquierda-derecha)
	frontBalance = frontLeft - frontRight
	rearBalance = rearLeft - rearRight

	return
}

// GetTempSpread calcula la diferencia de temperatura en un neumático
func (tt *TyresTracker) GetTempSpread(position TyrePosition) float32 {
	tyre := tt.tyres[position]

	temps := []float32{tyre.TempInner, tyre.TempMiddle, tyre.TempOuter}

	minTemp := temps[0]
	maxTemp := temps[0]

	for _, temp := range temps {
		if temp < minTemp {
			minTemp = temp
		}
		if temp > maxTemp {
			maxTemp = temp
		}
	}

	return maxTemp - minTemp
}

// IsCamberOptimal verifica si el camber es óptimo basado en temperaturas
func (tt *TyresTracker) IsCamberOptimal(position TyrePosition) bool {
	tyre := tt.tyres[position]

	// Camber óptimo: temperatura interior ligeramente mayor que exterior
	// Diferencia ideal: 5-10°C
	tempDiff := tyre.TempInner - tyre.TempOuter

	return tempDiff >= 5.0 && tempDiff <= 10.0
}

// GetPressureLoss calcula la pérdida de presión desde el inicio
func (tt *TyresTracker) GetPressureLoss() [4]float32 {
	// Esta función requeriría almacenar presiones iniciales
	// Por ahora retornamos ceros
	return [4]float32{0, 0, 0, 0}
}

// SetCompound establece el compuesto de neumáticos
func (tt *TyresTracker) SetCompound(compound TyreCompound) {
	tt.compound = compound

	// Ajustar temperaturas óptimas según compuesto
	if compound == CompoundWet {
		tt.optimalTempMin = 50.0
		tt.optimalTempMax = 80.0
	} else {
		tt.optimalTempMin = 75.0
		tt.optimalTempMax = 95.0
	}
}

// GetCompound devuelve el compuesto actual
func (tt *TyresTracker) GetCompound() TyreCompound {
	return tt.compound
}

// Reset reinicia el tracker
func (tt *TyresTracker) Reset() {
	tt.tyres = [4]TyreData{}
	tt.wearHistory = make([][4]float32, 0)
}

// GetWearDifference calcula la diferencia de desgaste entre neumáticos
func (tt *TyresTracker) GetWearDifference() (frontDiff, rearDiff, leftDiff, rightDiff float32) {
	frontDiff = float32(math.Abs(float64(tt.tyres[FrontLeft].Wear - tt.tyres[FrontRight].Wear)))
	rearDiff = float32(math.Abs(float64(tt.tyres[RearLeft].Wear - tt.tyres[RearRight].Wear)))
	leftDiff = float32(math.Abs(float64(tt.tyres[FrontLeft].Wear - tt.tyres[RearLeft].Wear)))
	rightDiff = float32(math.Abs(float64(tt.tyres[FrontRight].Wear - tt.tyres[RearRight].Wear)))
	return
}
