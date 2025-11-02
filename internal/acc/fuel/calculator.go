package fuel

import (
	"RaceAll/internal/acc/cars"
	"math"
)

// FuelData representa los datos de combustible
type FuelData struct {
	CurrentFuel        float32 // Litros actuales
	MaxFuel            float32 // Capacidad máxima
	FuelPercent        float32 // Porcentaje 0.0-100.0
	LastLapConsumption float32 // Consumo en última vuelta
	AvgConsumption     float32 // Consumo promedio por vuelta
	EstimatedLaps      float32 // Vueltas estimadas con combustible actual
	IsLow              bool    // Combustible bajo (<10%)
	IsCritical         bool    // Combustible crítico (<5%)
}

// FuelCalculator realiza cálculos de combustible
type FuelCalculator struct {
	carModel           cars.CarModel
	maxFuel            float32
	consumptionHistory []float32
	lastFuelLevel      float32
	initialized        bool
}

// NewFuelCalculator crea un nuevo calculador de combustible
func NewFuelCalculator(carModel cars.CarModel) *FuelCalculator {
	return &FuelCalculator{
		carModel:           carModel,
		maxFuel:            cars.GetMaxFuel(carModel),
		consumptionHistory: make([]float32, 0),
		lastFuelLevel:      0,
		initialized:        false,
	}
}

// Update actualiza los cálculos de combustible
func (fc *FuelCalculator) Update(currentFuel float32, lapCompleted bool) FuelData {
	// Inicializar en la primera actualización
	if !fc.initialized {
		fc.lastFuelLevel = currentFuel
		fc.initialized = true
	}

	// Calcular consumo de última vuelta si se completó una vuelta
	var lastLapConsumption float32
	if lapCompleted && fc.lastFuelLevel > currentFuel {
		lastLapConsumption = fc.lastFuelLevel - currentFuel

		// Agregar al historial (máximo 10 vueltas)
		fc.consumptionHistory = append(fc.consumptionHistory, lastLapConsumption)
		if len(fc.consumptionHistory) > 10 {
			fc.consumptionHistory = fc.consumptionHistory[1:]
		}

		fc.lastFuelLevel = currentFuel
	}

	// Calcular consumo promedio
	avgConsumption := fc.calculateAverageConsumption()

	// Calcular vueltas estimadas
	estimatedLaps := float32(0)
	if avgConsumption > 0 {
		estimatedLaps = currentFuel / avgConsumption
	}

	// Calcular porcentaje
	fuelPercent := float32(0)
	if fc.maxFuel > 0 {
		fuelPercent = (currentFuel / fc.maxFuel) * 100.0
	}

	return FuelData{
		CurrentFuel:        currentFuel,
		MaxFuel:            fc.maxFuel,
		FuelPercent:        fuelPercent,
		LastLapConsumption: lastLapConsumption,
		AvgConsumption:     avgConsumption,
		EstimatedLaps:      estimatedLaps,
		IsLow:              fuelPercent < 10.0,
		IsCritical:         fuelPercent < 5.0,
	}
}

// calculateAverageConsumption calcula el consumo promedio basado en el historial
func (fc *FuelCalculator) calculateAverageConsumption() float32 {
	if len(fc.consumptionHistory) == 0 {
		return 0
	}

	var sum float32
	for _, consumption := range fc.consumptionHistory {
		sum += consumption
	}

	return sum / float32(len(fc.consumptionHistory))
}

// CalculateFuelForLaps calcula cuánto combustible se necesita para N vueltas
func (fc *FuelCalculator) CalculateFuelForLaps(laps int) float32 {
	avgConsumption := fc.calculateAverageConsumption()
	if avgConsumption <= 0 {
		// Estimación por defecto: 3 litros por vuelta
		avgConsumption = 3.0
	}

	fuelNeeded := avgConsumption * float32(laps)

	// Agregar 5% de margen de seguridad
	fuelNeeded *= 1.05

	// No exceder el máximo
	if fuelNeeded > fc.maxFuel {
		fuelNeeded = fc.maxFuel
	}

	return fuelNeeded
}

// CalculateLapsWithFuel calcula cuántas vueltas se pueden hacer con cierta cantidad
func (fc *FuelCalculator) CalculateLapsWithFuel(fuelAmount float32) float32 {
	avgConsumption := fc.calculateAverageConsumption()
	if avgConsumption <= 0 {
		avgConsumption = 3.0 // Estimación por defecto
	}

	return fuelAmount / avgConsumption
}

// ShouldRefuel determina si se debe repostar basado en vueltas restantes
func (fc *FuelCalculator) ShouldRefuel(currentFuel float32, lapsRemaining int) bool {
	estimatedLaps := fc.CalculateLapsWithFuel(currentFuel)

	// Repostar si no alcanza para las vueltas restantes + 2 vueltas de margen
	return estimatedLaps < float32(lapsRemaining+2)
}

// CalculateRefuelAmount calcula cuánto combustible agregar
func (fc *FuelCalculator) CalculateRefuelAmount(currentFuel float32, lapsRemaining int) float32 {
	fuelNeeded := fc.CalculateFuelForLaps(lapsRemaining)
	refuelAmount := fuelNeeded - currentFuel

	if refuelAmount < 0 {
		refuelAmount = 0
	}

	// No exceder capacidad máxima
	if currentFuel+refuelAmount > fc.maxFuel {
		refuelAmount = fc.maxFuel - currentFuel
	}

	return refuelAmount
}

// GetConsumptionTrend devuelve la tendencia de consumo (positivo = aumentando)
func (fc *FuelCalculator) GetConsumptionTrend() float32 {
	histLen := len(fc.consumptionHistory)
	if histLen < 3 {
		return 0
	}

	// Comparar promedio de últimas 3 vueltas vs anteriores
	recentAvg := (fc.consumptionHistory[histLen-1] +
		fc.consumptionHistory[histLen-2] +
		fc.consumptionHistory[histLen-3]) / 3.0

	olderAvg := float32(0)
	count := 0
	for i := 0; i < histLen-3; i++ {
		olderAvg += fc.consumptionHistory[i]
		count++
	}
	if count > 0 {
		olderAvg /= float32(count)
	}

	return recentAvg - olderAvg
}

// IsConsumptionStable verifica si el consumo es estable
func (fc *FuelCalculator) IsConsumptionStable() bool {
	if len(fc.consumptionHistory) < 3 {
		return false
	}

	// Calcular desviación estándar
	avg := fc.calculateAverageConsumption()
	var variance float32

	for _, consumption := range fc.consumptionHistory {
		diff := consumption - avg
		variance += diff * diff
	}
	variance /= float32(len(fc.consumptionHistory))

	stdDev := float32(math.Sqrt(float64(variance)))

	// Considerar estable si la desviación es < 10% del promedio
	return stdDev < (avg * 0.1)
}

// GetFuelToAdd calcula combustible óptimo para agregar en pit stop
func (fc *FuelCalculator) GetFuelToAdd(currentFuel float32, lapsToFinish int, targetLaps int) float32 {
	// Calcular combustible necesario para las vueltas objetivo
	fuelForTarget := fc.CalculateFuelForLaps(targetLaps)

	// Calcular cuánto agregar
	toAdd := fuelForTarget - currentFuel

	// Asegurar que no se exceda el máximo
	if currentFuel+toAdd > fc.maxFuel {
		toAdd = fc.maxFuel - currentFuel
	}

	// No agregar valores negativos
	if toAdd < 0 {
		toAdd = 0
	}

	return toAdd
}

// Reset reinicia el calculador
func (fc *FuelCalculator) Reset() {
	fc.consumptionHistory = make([]float32, 0)
	fc.lastFuelLevel = 0
	fc.initialized = false
}
