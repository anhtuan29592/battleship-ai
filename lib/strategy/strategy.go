package strategy

import "github.com/anhtuan29592/battleship-ai/lib"

type Strategy interface {
	StartGame(boardSize lib.Size)
	ArrangeShip()
	GetShot() (point lib.Point)
	ShotHit(point lib.Point, sunk bool)
	ShotMiss(point lib.Point)
}
