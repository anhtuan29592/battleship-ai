package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type IShipAction interface {
	GetPositions() []*lib.Point
	ConflictWith(other *Ship) bool
	IsValid(boardSize lib.Size) bool
	UpdateLocation(orientation constant.Orientation, point *lib.Point)
	GetType() constant.ShipType
}

type Ship struct {
	ShipAction IShipAction
}
