package cars

// BrakeBias representa la distribución de frenos
type BrakeBias struct {
	FrontPercent float32 // 50.0 - 70.0 típicamente
}

// NewBrakeBias crea una nueva configuración de balance de frenos
func NewBrakeBias(frontPercent float32) BrakeBias {
	return BrakeBias{
		FrontPercent: clamp(frontPercent, 50.0, 70.0),
	}
}

// GetRearPercent devuelve el porcentaje de freno trasero
func (b *BrakeBias) GetRearPercent() float32 {
	return 100.0 - b.FrontPercent
}

// IsBalanced verifica si el balance está equilibrado (cerca de 55-60%)
func (b *BrakeBias) IsBalanced() bool {
	return b.FrontPercent >= 55.0 && b.FrontPercent <= 60.0
}

// IsFrontBiased verifica si está sesgado hacia adelante
func (b *BrakeBias) IsFrontBiased() bool {
	return b.FrontPercent > 60.0
}

// IsRearBiased verifica si está sesgado hacia atrás
func (b *BrakeBias) IsRearBiased() bool {
	return b.FrontPercent < 55.0
}

// clamp limita un valor entre min y max
func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
