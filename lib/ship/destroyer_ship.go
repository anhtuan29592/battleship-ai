package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type DestroyerShip struct {
	Location    *lib.Point
	Orientation constant.Orientation
}

func (c *DestroyerShip) GetPositions() []*lib.Point {
	positions := make([]*lib.Point, 0)
	if c.Orientation == constant.HORIZONTAL {
		for i := 0; i < 2; i++ {
			positions = append(positions, &lib.Point{X: c.Location.X + i, Y: c.Location.Y})
		}
	} else {
		for i := 0; i < 2; i++ {
			positions = append(positions, &lib.Point{X: c.Location.X, Y: c.Location.Y + i})
		}
	}

	return positions
}

func (c *DestroyerShip) ConflictWith(other *Ship) bool {
	otherPositions := other.ShipAction.GetPositions()
	myPositions := c.GetPositions()

	for i := 0; i < len(otherPositions); i++ {
		for j := 0; j < len(myPositions); i++ {
			if otherPositions[i].X == myPositions[j].X && otherPositions[i].Y == myPositions[j].Y {
				return true
			}
		}
	}

	return false
}

func (c *DestroyerShip) IsValid(boardSize lib.Size) bool {
	if c.Location == nil {
		return false
	}

	if c.Location.X < 0 || c.Location.Y < 0 {
		return false
	}

	if c.Orientation == constant.HORIZONTAL {
		if c.Location.X >= boardSize.Witdh || c.Location.X + 2 > boardSize.Witdh {
			return false
		}
	} else {
		if c.Location.Y >= boardSize.Height || c.Location.Y + 2 > boardSize.Height {
			return false
		}
	}

	return true
}

func (c *DestroyerShip) UpdateLocation(orientation constant.Orientation, point *lib.Point) {
	c.Location = point
}

func (c *DestroyerShip) GetType() constant.ShipType {
	return constant.DESTROYER
}
