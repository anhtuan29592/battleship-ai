package strategy

import (
	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/ship"
	"math/rand"
	"github.com/anhtuan29592/paladin/lib/constant"
	"sort"
	"fmt"
)

type Strategy interface {
	StartGame(boardSize lib.Size, ships []ship.Ship)
	GetShot() (point lib.Point)
	ShotHit(point lib.Point, shipType string, shipPositions []lib.Point)
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

func SetUpPriorityShots(boardSize lib.Size, shotPatterns []lib.Point) []lib.Point {
	priorityShots := make([]lib.Point, 0)
	for i := len(shotPatterns) - 1; i >= 0; i-- {
		tmp := shotPatterns[i]
		//if (0 <= tmp.X && tmp.X < 1) || (boardSize.Width - 1 <= tmp.X && tmp.X < boardSize.Width) || (0 <= tmp.Y && tmp.Y < 1) || (boardSize.Height - 1 <= tmp.Y && tmp.Y < boardSize.Height) {
		//	priorityShots = append(priorityShots, tmp)
		//	shotPatterns = append(shotPatterns[:i], shotPatterns[i+1:]...)
		//}
		// center
		if (2 <= tmp.X && tmp.X < boardSize.Width-2) && (2 <= tmp.Y && tmp.Y < boardSize.Height-2) {
			priorityShots = append(priorityShots, tmp)
			continue
		}

		// corners
		if (0 <= tmp.X && tmp.X < 2 && 0 <= tmp.Y && tmp.Y < 2) || (boardSize.Width-2 <= tmp.X && tmp.X < boardSize.Width && 0 <= tmp.Y && tmp.Y < 2) || (0 <= tmp.X && tmp.X < 2 && boardSize.Height-2 <= tmp.Y && tmp.Y < boardSize.Height) || (boardSize.Width-2 <= tmp.X && tmp.X < boardSize.Width && boardSize.Height-2 <= tmp.Y && tmp.Y < boardSize.Height) {
			priorityShots = append(priorityShots, tmp)
			continue
		}

		// center vertexes
		if (boardSize.Width/2-2 <= tmp.X && tmp.X < boardSize.Width/2+2 && (0 <= tmp.Y && tmp.Y < 2 || boardSize.Height-2 <= tmp.Y && tmp.Y < boardSize.Height)) || ((0 <= tmp.X && tmp.X < 2 || boardSize.Width-2 <= tmp.X && tmp.X < boardSize.Width) && boardSize.Height/2-1 <= tmp.Y && tmp.Y < boardSize.Height/2+1) {
			priorityShots = append(priorityShots, tmp)
			continue
		}

		// quarter vertexes
		if (boardSize.Width/4-2 <= tmp.X && tmp.X < boardSize.Width/4+2 && (0 <= tmp.Y && tmp.Y < 2 || boardSize.Height-2 <= tmp.Y && tmp.Y < boardSize.Height)) || (boardSize.Width*3/4-2 <= tmp.X && tmp.X < boardSize.Width*3/4+2 && (0 <= tmp.Y && tmp.Y < 2 || boardSize.Height-2 <= tmp.Y && tmp.Y < boardSize.Height)) {
			priorityShots = append(priorityShots, tmp)
			continue
		}
	}
	return priorityShots
}

func ArrangeShips(boardSize lib.Size, ships []ship.Ship) []ship.Ship {

	for i := 0; i < len(ships); {
		retryCount := 0
		for {
			for {
				ships[i].UpdateLocation(constant.Orientation(rand.Intn(2)), lib.Point{X: rand.Intn(boardSize.Width), Y: rand.Intn(boardSize.Height)})
				if ships[i].IsValid(boardSize) {
					break
				}
			}
			hasConflict := false
			for j := 0; j < i; j++ {
				if ships[i].ConflictWith(ships[j], boardSize) {
					hasConflict = true
					break
				}
			}

			if !hasConflict {
				retryCount = 0
				break
			}

			if retryCount >= 100 {
				break
			}

			retryCount++
		}

		if retryCount != 0 {
			i = 0
		} else {
			i++
		}
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

func PrintPoints(boardSize lib.Size, points []lib.Point) {
	for r := 0; r < boardSize.Height; r++ {
		fmt.Print("|")
		for c := 0; c < boardSize.Width; c++ {
			printed := false
			for i := 0; i < len(points); i++ {
				if points[i].X == c && points[i].Y == r {
					fmt.Print("x")
					printed = true
					break
				}
			}
			if !printed {
				fmt.Print(" ")
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println()
}
