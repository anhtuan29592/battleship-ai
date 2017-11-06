package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/ship"
	"math/rand"
	"github.com/anhtuan29592/battleship-ai/lib/constant"
	"sort"
)

type Strategy interface {
	StartGame(boardSize lib.Size)
	GetShot() (point lib.Point)
	ShotHit(point lib.Point, shipPositions []lib.Point)
	ShotMiss(point lib.Point)
}

type GameAI struct {
	Strategy Strategy
	Mixin    Mixin
}

func SetUpShotPattern(boardSize lib.Size) []lib.Point {
	shotPatterns := make([]lib.Point, 0)

	for r := 0; r < boardSize.Height; r++ {
		for c := 0; c < boardSize.Width; c++ {
			if (r+c)%2 == 0 {
				shotPatterns = append(shotPatterns, lib.Point{X: c, Y: r})
			}
		}
	}

	return shotPatterns
}

func ArrangeShips(boardSize lib.Size, ships []ship.Ship, touchDistance int) []ship.Ship {

	hasConflict := func() bool {
		for i := 0; i < len(ships)-1; i++ {
			for j := i + 1; j < len(ships); j++ {
				if ships[i].ConflictWith(ships[j]) {
					return true
				}

				if !ships[i].Touch(ships[j], 1) {
					return false
				}
			}
		}
		return false
	}

	validTouch := func() int {
		count := 0
		for i := 0; i < len(ships)-1; i++ {
			for j := i + 1; j < len(ships); j++ {
				if ships[i].Touch(ships[j], touchDistance) {
					count++
				}
			}
		}
		return count
	}

	validOnBoard := func() bool {
		for i := 0; i < len(ships); i++ {
			if !ships[i].IsValid(boardSize) {
				return false
			}
		}
		return true
	}

	for {
		if !hasConflict() && validOnBoard() {
			break
		}

		for i := 0; i < len(ships); i++ {
			x := rand.Intn(boardSize.Width - 1)
			y := rand.Intn(boardSize.Height - 1)
			orientation := constant.Orientation(rand.Intn(2))
			ships[i].UpdateLocation(orientation, lib.Point{X: x, Y: y})
		}
	}

	if validTouch() < len(ships) {
		return ArrangeShips(boardSize, ships, touchDistance-1)
	}

	return ships
}

func GetDirection(p1 lib.Point, p2 lib.Point) constant.Direction {
	if p1.Y == p2.Y {
		if p1.X >= p2.X {
			return constant.RIGHT
		} else {
			return constant.LEFT
		}
	}

	// Vertical
	if p1.X == p2.X {
		if p1.Y >= p2.Y {
			return constant.DOWN
		} else {
			return constant.UP
		}
	}

	return constant.UP
}

func SortPoints(s []lib.Point, orientation constant.Orientation, isAscending bool) []lib.Point {
	if orientation == constant.HORIZONTAL {
		sort.Slice(s, func(i, j int) bool {
			if isAscending {
				return s[i].X < s[j].X
			}
			return s[i].X >= s[j].X
		})
	}

	if orientation == constant.VERTICAL {
		sort.Slice(s, func(i, j int) bool {
			if isAscending {
				return s[i].Y < s[j].Y
			}
			return s[i].Y >= s[j].Y
		})
	}

	return s
}

func CheckPointInSlice(slice []lib.Point, point lib.Point) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].X == point.X && slice[i].Y == point.Y {
			return true
		}
	}

	return false
}
