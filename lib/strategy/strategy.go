package strategy

import (
	"math/rand"

	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
	"github.com/anhtuan29592/paladin/lib/ship"
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

func SetUpShotPattern(boardSize lib.Size) []lib.PriorityPoint {
	shotPatterns := make([]lib.PriorityPoint, 0)

	for r := 0; r < boardSize.Height; r++ {
		for c := 0; c < boardSize.Width; c++ {
			if (r+c)%2 == 0 {
				shotPatterns = append(shotPatterns, lib.PriorityPoint{Location: lib.Point{X: c, Y: r}, Score: 0})
			}
		}
	}

	return shotPatterns
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
