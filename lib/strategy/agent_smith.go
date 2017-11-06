package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/constant"
	"math/rand"
	"encoding/json"
	"log"
)

type AgentSmith struct {
	Shots             []lib.Point
	HitShots          []lib.Point
	ShotPatterns      []lib.Point
	BoardSize         lib.Size
	ShotFront         *lib.Point
	ShotRear          *lib.Point
	LastShotDirection constant.Direction
	MissCount         int
}

func (a *AgentSmith) GetGameState() lib.GameState {
	gameData, _ := json.Marshal(a)
	return lib.GameState{
		GameStrategy: constant.AGENT_SMITH_STRATEGY,
		GameData:     string(gameData),
	}
}

func (a *AgentSmith) LoadGameState(state lib.GameState) {
	if state.GameStrategy != constant.AGENT_SMITH_STRATEGY {
		return
	}

	data := state.GameData
	json.Unmarshal([]byte(data), a)
}

func (a *AgentSmith) StartGame(boardSize lib.Size) {
	a.MissCount = 0
	a.Shots = make([]lib.Point, 0)
	a.HitShots = make([]lib.Point, 0)
	a.BoardSize = boardSize
	a.ShotPatterns = SetUpShotPattern(boardSize)
}

func (a *AgentSmith) GetShot() (point lib.Point) {
	var shot lib.Point
	if a.ShotFront != nil {
		shot = a.FireTargeted()
	} else {
		shot = a.FireRandom()
	}
	a.Shots = append(a.Shots, shot)
	if a.ShotRear != nil {
		a.LastShotDirection = GetDirection(shot, *a.ShotRear)
	}
	return shot
}

func (a *AgentSmith) ShotHit(point lib.Point, shipPositions []lib.Point) {
	sunk := len(shipPositions) > 0
	log.Printf("hit location %s, sunk %s", point, sunk)
	a.HitShots = append(a.HitShots, point)
	a.MissCount = 0
	a.ShotFront = &point
	if a.ShotRear == nil {
		a.ShotRear = &point
	}

	if sunk {
		a.ShotFront = nil
		a.ShotRear = nil
		a.MissCount = 0
	}
}

func (a *AgentSmith) ShotMiss(point lib.Point) {
	a.MissCount++
	if a.MissCount == 4 {
		a.ShotFront = nil
		a.ShotRear = nil
		a.MissCount = 0
	} else {
		a.LastShotDirection = a.LastShotDirection.Invert()
	}
}

func (*AgentSmith) FireDirected(direction constant.Direction, target lib.Point) lib.Point {
	switch direction {
	case constant.UP:
		return lib.Point{X: target.X, Y: target.Y - 1}
	case constant.DOWN:
		return lib.Point{X: target.X, Y: target.Y + 1}
	case constant.LEFT:
		return lib.Point{X: target.X - 1, Y: target.Y}
	case constant.RIGHT:
		return lib.Point{X: target.X + 1, Y: target.Y}
	default:
		return target
	}
}

func (a *AgentSmith) FireAroundPoint(p lib.Point) lib.Point {

	testShot := a.FireDirected(a.LastShotDirection, p)
	if !a.ValidShot(testShot) {
		testShot = a.FireDirected(a.LastShotDirection.Invert(), p)
	}

	if !a.ValidShot(testShot) {
		testShot = a.FireDirected(constant.UP, p)
	}

	if !a.ValidShot(testShot) {
		testShot = a.FireDirected(constant.LEFT, p)
	}

	if !a.ValidShot(testShot) {
		testShot = a.FireDirected(constant.RIGHT, p)
	}

	if !a.ValidShot(testShot) {
		testShot = a.FireDirected(constant.DOWN, p)
	}

	return testShot
}

func (a *AgentSmith) FireTargeted() lib.Point {
	for {
		point := a.FireAroundPoint(lib.Point{X: a.ShotFront.X, Y: a.ShotFront.Y})
		if a.ValidShot(point) {
			return point
		}

		if a.ShotFront != a.ShotRear {
			a.ShotFront = a.ShotRear
		} else {
			a.ShotFront = nil
			a.ShotRear = nil
			a.MissCount = 0
			break
		}
	}

	return a.FireRandom()
}

func (a *AgentSmith) FireRandom() lib.Point {
	var point lib.Point
	for {
		if len(a.ShotPatterns) > 0 {
			i := rand.Intn(len(a.ShotPatterns))
			point = a.ShotPatterns[i]
			a.ShotPatterns = append(a.ShotPatterns[:i], a.ShotPatterns[i+1:]...)
		} else {
			for {
				point = a.FireAroundPoint(a.HitShots[0])
				if a.ValidShot(point) {
					break
				}

				a.HitShots = append(a.HitShots[:0], a.HitShots[1:]...)
				if len(a.HitShots) <= 0 {
					break
				}
			}
		}

		if a.ValidShot(point) {
			break
		}
	}
	return point
}

func (a *AgentSmith) ValidShot(p lib.Point) bool {
	for i := 0; i < len(a.Shots); i++ {
		s := a.Shots[i]
		if s.X == p.X && s.Y == p.Y {
			return false
		}
	}

	if p.X < 0 || p.Y < 0 {
		return false
	}

	if p.X >= a.BoardSize.Width || p.Y >= a.BoardSize.Height {
		return false
	}

	return true
}