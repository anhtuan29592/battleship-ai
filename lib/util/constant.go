package constant

type Orientation int

const (
	HORIZONTAL Orientation = iota + 1
	VERTICAL
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

type Direction int

const (
	UP Direction = iota + 1
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
