package ship

import (
	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
	"github.com/anhtuan29592/paladin/lib/util"
)

var CARRIER_H_POINTS = []lib.Point{
	{X: 1, Y: 0},
	{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1},
}

var CARRIER_V_POINTS = []lib.Point{
	{X: 1, Y: 0},
	{X: 0, Y: 1}, {X: 1, Y: 1},
	{X: 1, Y: 2},
	{X: 1, Y: 3},
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

func (c *CarrierShip) Zoom(boardSize lib.Size, location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := c.GetPositions(location, orientation)
	for i := len(positions) - 1; i >= 0; i-- {
		pos := positions[i]
		// up
		point := lib.Point{X: pos.X, Y: pos.Y - 1}
		if point.ValidInBoard(boardSize) && !util.CheckPointInSlice(positions, point) {
			positions = append(positions, point)
		}

		// down
		point = lib.Point{X: pos.X, Y: pos.Y + 1}
		if point.ValidInBoard(boardSize) && !util.CheckPointInSlice(positions, point) {
			positions = append(positions, point)
		}

		// left
		point = lib.Point{X: pos.X - 1, Y: pos.Y}
		if point.ValidInBoard(boardSize) && !util.CheckPointInSlice(positions, point) {
			positions = append(positions, point)
		}

		// right
		point = lib.Point{X: pos.X + 1, Y: pos.Y}
		if point.ValidInBoard(boardSize) && !util.CheckPointInSlice(positions, point) {
			positions = append(positions, point)
		}
	}

	return positions
}
