package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
)

type CarrierShip struct {
	Location    *lib.Point
	Orientation int
}

func (self *CarrierShip) GetPositions() ([]*lib.Point) {
	positions := make([]*lib.Point, 0)
	positions = append(positions, &lib.Point{X: self.Location.X, Y: self.Location.Y})
	return positions
}

func (self *CarrierShip) ConflictWith(other *Ship) (bool) {
	return false
}

func (self *CarrierShip) IsValid(boardW int, boardH int) (bool) {
	if self.Location.X < 0 || self.Location.Y < 0 {
		return false
	}

	return true
}
