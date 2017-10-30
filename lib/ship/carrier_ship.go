package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
)

type CarrierShip struct {
	Ship
	position *lib.Point
}

func (s CarrierShip) get_positions() ([]*lib.Point) {
	positions := make([]*lib.Point, 0)
	//positions = append(positions, lib.Point{s.position.x, s.position.y})
	return positions
}