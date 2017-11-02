package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/ship"
)

type Strategy interface {
	StartGame(boardSize lib.Size)
	ArrangeShips(ships []ship.Ship) []ship.Ship
	GetShot() (point lib.Point)
	ShotHit(point lib.Point, sunk bool)
	ShotMiss(point lib.Point)
}

type GameAI struct {
	Strategy Strategy
	Mixin    Mixin
}
