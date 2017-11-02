package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
	"math/rand"
	"github.com/anhtuan29592/battleship-ai/lib/ship"
	"encoding/json"
)

type AgentSmith struct {
	Shots        []lib.Point
	HitShots     []lib.Point
	ShotPatterns []ShotPattern
	BoardSize    lib.Size
	ShotFront    *lib.Point
	ShotRear     *lib.Point
	MissCount    int
}

type ShotPattern struct {
	Location lib.Point
	Score    int
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
	a.SetUpShotPattern(boardSize)
}

func (a *AgentSmith) ArrangeShips(ships []ship.Ship) []ship.Ship {
	if a.Validate(ships) {
		return ships
	}

	for i := 0; i < len(ships); i++ {
		x := rand.Intn(a.BoardSize.Width - 1)
		y := rand.Intn(a.BoardSize.Height - 1)
		orientation := constant.Orientation(rand.Intn(2))
		ships[i].UpdateLocation(orientation, lib.Point{X: x, Y: y})
	}

	return a.ArrangeShips(ships)
}

func (a *AgentSmith) GetShot() (point lib.Point) {
	var shot lib.Point
	if a.ShotFront != nil {
		shot = a.FireTargeted()
	} else {
		shot = a.FireRandom()
	}
	a.Shots = append(a.Shots, shot)
	return shot
}

func (a *AgentSmith) ShotHit(point lib.Point, sunk bool) {
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
	if a.MissCount == 6 {
		a.ShotFront = nil
		a.ShotRear = nil
		a.MissCount = 0
	}
}

func (a *AgentSmith) Validate(ships []ship.Ship) bool {
	for i := 0; i < len(ships); i++ {
		if !ships[i].IsValid(a.BoardSize) {
			return false
		}
	}

	for i := 0; i < len(ships); i++ {
		for j := i + 1; j < len(ships); j++ {
			if ships[i].ConflictWith(ships[j]) {
				return false
			}
		}
	}

	return true
}

func (a *AgentSmith) SetUpShotPattern(boardSize lib.Size) {
	a.ShotPatterns = make([]ShotPattern, 0)
	for y := 0; y < boardSize.Height; y++ {
		for x := 0; x < boardSize.Width; x++ {
			if (x+y)%2 == 0 {
				a.ShotPatterns = append(a.ShotPatterns, ShotPattern{lib.Point{X: x, Y: y}, a.GetScore(x, y)})
			}
		}
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
	testShot := a.FireDirected(constant.UP, p)
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
			point = a.ShotPatterns[i].Location
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

func (a *AgentSmith) GetScore(x int, y int) int {
	return 0
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
