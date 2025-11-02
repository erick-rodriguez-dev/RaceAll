package cars

import "RaceAll/internal/broadcast"

// CarModel representa los modelos de autos disponibles en ACC
type CarModel = byte

// CarInfo contiene información detallada de cada modelo de auto
type CarInfo struct {
	Model        CarModel
	Name         string
	Manufacturer string
	Category     CarCategory
	Year         int
	MaxFuel      float32 // Litros
}

// CarCategory representa la categoría del auto
type CarCategory byte

const (
	GT3 CarCategory = iota
	GT4
	GT2
	Cup
	SuperTrofeo
	ChallengeEvo
)

// GetCarInfo devuelve la información de un modelo de auto
func GetCarInfo(model CarModel) CarInfo {
	// Usar el nombre desde broadcast.CarModels
	name := broadcast.GetCarModelName(model)

	// Determinar categoría, fabricante y specs según el modelo
	var category CarCategory
	var manufacturer string
	var year int
	var maxFuel float32

	switch {
	// GT3 Cars
	case model <= 25, model >= 30 && model <= 36:
		category = GT3
		maxFuel = 120
		manufacturer = extractManufacturer(name)
		year = extractYear(name)

	// Cup/Challenge
	case model >= 26 && model <= 29:
		category = Cup
		maxFuel = 100
		manufacturer = extractManufacturer(name)
		year = extractYear(name)

	// GT4 Cars
	case model >= 50 && model <= 61:
		category = GT4
		maxFuel = 100
		manufacturer = extractManufacturer(name)
		year = extractYear(name)

	// GT2 Cars
	case model >= 80 && model <= 86:
		category = GT2
		maxFuel = 110
		manufacturer = extractManufacturer(name)
		year = extractYear(name)

	default:
		category = GT3
		maxFuel = 120
		manufacturer = "Unknown"
		year = 2024
	}

	return CarInfo{
		Model:        model,
		Name:         name,
		Manufacturer: manufacturer,
		Category:     category,
		Year:         year,
		MaxFuel:      maxFuel,
	}
}

// extractManufacturer extrae el fabricante del nombre del auto
func extractManufacturer(name string) string {
	// Extraer la primera palabra del nombre (generalmente el fabricante)
	for i, c := range name {
		if c == ' ' {
			return name[:i]
		}
	}
	return name
}

// extractYear extrae el año del nombre del auto
func extractYear(name string) int {
	// Buscar años comunes en los nombres (2012-2024)
	years := []string{"2024", "2023", "2022", "2021", "2020", "2019", "2018", "2017", "2016", "2015", "2014", "2013", "2012"}

	for _, yearStr := range years {
		if contains(name, yearStr) {
			// Convertir string a int
			year := 0
			for _, c := range yearStr {
				year = year*10 + int(c-'0')
			}
			return year
		}
	}

	return 2024 // Por defecto
}

// contains verifica si un string contiene otro
func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringContains(s, substr)
}

// stringContains implementación simple de contains
func stringContains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// GetCarName devuelve el nombre del modelo de auto
func GetCarName(model CarModel) string {
	return broadcast.GetCarModelName(model)
}

// GetCarCategory devuelve la categoría del auto
func GetCarCategory(model CarModel) CarCategory {
	return GetCarInfo(model).Category
}

// GetMaxFuel devuelve la capacidad máxima de combustible en litros
func GetMaxFuel(model CarModel) float32 {
	return GetCarInfo(model).MaxFuel
}

// IsGT3 verifica si el auto es GT3
func IsGT3(model CarModel) bool {
	return GetCarCategory(model) == GT3
}

// IsGT4 verifica si el auto es GT4
func IsGT4(model CarModel) bool {
	return GetCarCategory(model) == GT4
}

// IsGT2 verifica si el auto es GT2
func IsGT2(model CarModel) bool {
	return GetCarCategory(model) == GT2
}

// IsCup verifica si el auto es de categoría Cup
func IsCup(model CarModel) bool {
	category := GetCarCategory(model)
	return category == Cup || category == ChallengeEvo || category == SuperTrofeo
}
