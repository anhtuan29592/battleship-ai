package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
)

type CarrierShip struct {
	Location *lib.Point
	Orientation int
}

func (s *CarrierShip) GetPositions() ([]*lib.Point) {
	positions := make([]*lib.Point, 0)
	positions = append(positions, &lib.Point{s.Location.X, s.Location.Y})
	return positions
}

func (s *CarrierShip) ConflictWith(other *Ship) (bool) {
	return false
}

func (s *CarrierShip) IsValid(boardW int, boardH int) (bool) {
	if (s.Location.X < 0 || s.Location.Y < 0) {
		return false
	}

	if (s.Orientation == lib.HORIZONTAL) {
		
	}

	return true
}