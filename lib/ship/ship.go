package ship

import "github.com/anhtuan29592/battleship-ai/lib"

type Ship interface {
	GetPositions() ([]*lib.Point)
	ConflictWith(other *Ship) (bool)
	IsValid(boardW int, boardH int) (bool)
}
