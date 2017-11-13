package ship

import (
	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
)

type OilRigShip struct {
}

func (c *OilRigShip) GetSize(orientation constant.Orientation) lib.Size {
	switch orientation {
	case constant.HORIZONTAL:
		return lib.Size{Width: 2, Height: 2}
	case constant.VERTICAL:
		return lib.Size{Width: 2, Height: 2}
	default:
		return lib.Size{Width: 2, Height: 2}
	}
}

func (c *OilRigShip) GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := make([]lib.Point, 0)
	size := c.GetSize(orientation)
	for r := 0; r < size.Height; r++ {
		for c := 0; c < size.Width; c++ {
			positions = append(positions, lib.Point{X: location.X + c, Y: location.Y + r})
		}
	}
	return positions
}

func (c *OilRigShip) GetType() constant.ShipType {
	return constant.OIL_RIG
}
