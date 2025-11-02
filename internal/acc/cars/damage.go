package cars

// DamageLevel representa el nivel de daño en una parte del auto
type DamageLevel byte

const (
	NoDamage     DamageLevel = 0
	MinorDamage  DamageLevel = 1
	MediumDamage DamageLevel = 2
	HeavyDamage  DamageLevel = 3
	SevereDamage DamageLevel = 4
)

// CarDamage representa el estado de daño del auto
type CarDamage struct {
	Front        float32 // 0.0 - 1.0
	Rear         float32 // 0.0 - 1.0
	Left         float32 // 0.0 - 1.0
	Right        float32 // 0.0 - 1.0
	Center       float32 // 0.0 - 1.0
	SuspensionFL float32 // 0.0 - 1.0
	SuspensionFR float32 // 0.0 - 1.0
	SuspensionRL float32 // 0.0 - 1.0
	SuspensionRR float32 // 0.0 - 1.0
}

// GetDamageLevel convierte un valor float de daño a un nivel
func GetDamageLevel(damageValue float32) DamageLevel {
	if damageValue <= 0.0 {
		return NoDamage
	} else if damageValue <= 0.25 {
		return MinorDamage
	} else if damageValue <= 0.50 {
		return MediumDamage
	} else if damageValue <= 0.75 {
		return HeavyDamage
	}
	return SevereDamage
}

// HasDamage verifica si hay algún daño en el auto
func (d *CarDamage) HasDamage() bool {
	return d.Front > 0 || d.Rear > 0 || d.Left > 0 || d.Right > 0 || d.Center > 0
}

// HasSuspensionDamage verifica si hay daño en la suspensión
func (d *CarDamage) HasSuspensionDamage() bool {
	return d.SuspensionFL > 0 || d.SuspensionFR > 0 ||
		d.SuspensionRL > 0 || d.SuspensionRR > 0
}

// GetTotalDamage devuelve el daño total del auto (0.0 - 1.0)
func (d *CarDamage) GetTotalDamage() float32 {
	total := d.Front + d.Rear + d.Left + d.Right + d.Center
	total += d.SuspensionFL + d.SuspensionFR + d.SuspensionRL + d.SuspensionRR
	return total / 9.0 // Promedio de 9 componentes
}

// GetBodyDamage devuelve el daño del chasis (0.0 - 1.0)
func (d *CarDamage) GetBodyDamage() float32 {
	return (d.Front + d.Rear + d.Left + d.Right + d.Center) / 5.0
}

// GetSuspensionDamage devuelve el daño promedio de la suspensión
func (d *CarDamage) GetSuspensionDamage() float32 {
	return (d.SuspensionFL + d.SuspensionFR + d.SuspensionRL + d.SuspensionRR) / 4.0
}

// IsCritical verifica si el daño es crítico (>75%)
func (d *CarDamage) IsCritical() bool {
	return d.GetTotalDamage() > 0.75
}
