package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type DestroyerShip struct {
}

func (c *DestroyerShip) GetSize(orientation constant.Orientation) lib.Size {
	switch orientation {
	case constant.HORIZONTAL:
		return lib.Size{Width: 2, Height: 1}
	case constant.VERTICAL:
		return lib.Size{Width: 1, Height: 2}
	default:
		return lib.Size{Width: 0, Height: 0}
	}
}

func (c *DestroyerShip) GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := make([]lib.Point, 0)
	size := c.GetSize(orientation)
	for r := 0; r < size.Height; r++ {
		for c := 0; c < size.Width; c++ {
			positions = append(positions, lib.Point{X: location.X + c, Y: location.Y + r})
		}
	}
	return positions
}

func (c *DestroyerShip) GetType() constant.ShipType {
	return constant.DESTROYER
}
