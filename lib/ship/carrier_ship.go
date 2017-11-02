package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

var CARRIER_H_POINTS = []lib.Point{
								{X: 2, Y: 0},
	{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1},
}

var CARRIER_V_POINTS = []lib.Point{
	{X: 0, Y: 0},
	{X: 0, Y: 1}, {X: 1, Y: 1},
	{X: 0, Y: 2},
	{X: 0, Y: 3},
}

type CarrierShip struct {
}

func (c *CarrierShip) GetSize(orientation constant.Orientation) lib.Size {
	switch orientation {
	case constant.HORIZONTAL:
		return lib.Size{Width: 4, Height: 2}
	case constant.VERTICAL:
		return lib.Size{Width: 2, Height: 4}
	default:
		return lib.Size{Width: 0, Height: 0}
	}
}

func (c *CarrierShip) GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := make([]lib.Point, 0)
	size := c.GetSize(orientation)
	if orientation == constant.HORIZONTAL {
		for r := 0; r < size.Height; r++ {
			for c := 0; c < size.Width; c++ {
				for i := 0; i < len(CARRIER_H_POINTS); i++ {
					if CARRIER_H_POINTS[i].X == c && CARRIER_H_POINTS[i].Y == r {
						positions = append(positions, lib.Point{X: location.X + c, Y: location.Y + r})
					}
				}
			}
		}
	} else {
		for r := 0; r < size.Height; r++ {
			for c := 0; c < size.Width; c++ {
				for i := 0; i < len(CARRIER_V_POINTS); i++ {
					if CARRIER_V_POINTS[i].X == c && CARRIER_V_POINTS[i].Y == r {
						positions = append(positions, lib.Point{X: location.X + c, Y: location.Y + r})
					}
				}
			}
		}
	}

	return positions
}

func (c *CarrierShip) GetType() constant.ShipType {
	return constant.CARRIER
}
