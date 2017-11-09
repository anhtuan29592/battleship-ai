package ship

import (
	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
	"github.com/anhtuan29592/paladin/lib/util"
)

type CruiserShip struct {
}

func (c *CruiserShip) GetSize(orientation constant.Orientation) lib.Size {
	switch orientation {
	case constant.HORIZONTAL:
		return lib.Size{Width: 3, Height: 1}
	case constant.VERTICAL:
		return lib.Size{Width: 1, Height: 3}
	default:
		return lib.Size{Width: 0, Height: 0}
	}
}

func (c *CruiserShip) GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point {
	positions := make([]lib.Point, 0)
	size := c.GetSize(orientation)
	for r := 0; r < size.Height; r++ {
		for c := 0; c < size.Width; c++ {
			positions = append(positions, lib.Point{X: location.X + c, Y: location.Y + r})
		}
	}
	return positions
}

func (c *CruiserShip) GetType() constant.ShipType {
	return constant.CRUISER
}

func (c *CruiserShip) Zoom(boardSize lib.Size, location lib.Point, orientation constant.Orientation) []lib.Point {
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