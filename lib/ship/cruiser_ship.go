package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type CruiserShip struct {
}

func (c *CruiserShip) GetSize(orientation constant.Orientation) lib.Size {
	switch orientation {
	case constant.HORIZONTAL:
		return lib.Size{Width: 2, Height: 1}
	case constant.VERTICAL:
		return lib.Size{Width: 1, Height: 2}
	default:
		return lib.Size{Width: 0, Height: 0}
	}
}

func (c *CruiserShip) GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := make([]lib.Point, 0)
	size := c.GetSize(orientation)
	if orientation == constant.HORIZONTAL {
		for i := 0; i < size.Width; i++ {
			positions = append(positions, lib.Point{X: location.X + i, Y: location.Y})
		}
	} else {
		for i := 0; i < size.Height; i++ {
			positions = append(positions, lib.Point{X: location.X, Y: location.Y + i})
		}
	}

	return positions
}

func (c *CruiserShip) GetType() constant.ShipType {
	return constant.CRUISER
}
