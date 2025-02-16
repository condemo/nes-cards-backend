package types

type Poison struct {
	dmg uint8
}

func NewPosion(dmg uint8) Poison {
	return Poison{dmg: dmg}
}

func (p *Poison) GetDmg() uint8 {
	return p.dmg
}

func (p *Poison) AddDmg(dmg uint8) {
	p.dmg += dmg
}

type Confusion struct {
	turnsLeft uint8
}

// TODO:
func (c *Confusion) Apply() uint8 {
	// whatever...
	if c.turnsLeft > 0 {
		c.turnsLeft--
		return c.turnsLeft
	}

	return c.turnsLeft
}

type Intangible struct {
	turnsLeft uint8
}

// TODO:
func (i *Intangible) Apply() uint8 {
	// whatever...
	if i.turnsLeft > 0 {
		i.turnsLeft--
		return i.turnsLeft
	}
	return i.turnsLeft
}
