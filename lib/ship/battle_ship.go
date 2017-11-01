package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type BattleShip struct {
	Location    *lib.Point
	Orientation constant.Orientation
}

func (c *BattleShip) GetPositions() []*lib.Point {
	positions := make([]*lib.Point, 0)
	positions = append(positions, &lib.Point{X: c.Location.X, Y: c.Location.Y})
	return positions
}

func (c *BattleShip) ConflictWith(other *Ship) bool {
	return false
}

func (c *BattleShip) IsValid(boardSize lib.Size) bool {
	if c.Location.X < 0 || c.Location.Y < 0 {
		return false
	}

	return true
}

func (c *BattleShip) UpdateLocation(orientation constant.Orientation, point *lib.Point) {
	c.Location = point
}

func (c *BattleShip) GetType() constant.ShipType {
	return constant.BATTLE_SHIP
}
