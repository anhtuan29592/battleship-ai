package constant

type Orientation int

const (
	HORIZONTAL Orientation = iota
	VERTICAL
	UNKNOWN
)

func (o Orientation) String() string {
	switch o {
	case HORIZONTAL:
		return "HORIZONTAL"
	case VERTICAL:
		return "VERTICAL"
	default:
		return ""
	}
}

type ShipType string

const (
	CARRIER     ShipType = "CV"
	BATTLE_SHIP          = "BB"
	OIL_RIG              = "OR"
	CRUISER              = "CA"
	DESTROYER            = "DD"
)

func (s ShipType) String() string {
	switch s {
	case CARRIER:
		return "CV"
	case BATTLE_SHIP:
		return "BB"
	case OIL_RIG:
		return "OR"
	case CRUISER:
		return "CA"
	case DESTROYER:
		return "DD"
	default:
		return ""
	}
}

type Direction int

const (
	UP    Direction = iota
	DOWN
	LEFT
	RIGHT
)

func (d Direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	default:
		return ""
	}
}

func (d Direction) Invert() Direction {
	switch d {
	case UP:
		return DOWN
	case DOWN:
		return UP
	case LEFT:
		return RIGHT
	case RIGHT:
		return LEFT
	default:
		return UP
	}
}

const (
	AGENT_SMITH_STRATEGY = "ASM"
	SAMPLE_STRATEGY      = "SSG"
)

const (
	HIT   = "HIT"
	MISS  = "MISS"
	CLEAR = "CLEAR"
	SUNK  = "SUNK"
)