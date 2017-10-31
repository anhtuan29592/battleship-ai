package opponent

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
	"math/rand"
)

type AgentSmith struct {
	Shots        []*lib.Point
	HitShots     []*lib.Point
	ShotPatterns []*ShotPattern
}

type ShotPattern struct {
	Location *lib.Point
	Score    int
}

func (a *AgentSmith) StartGame(boardSize lib.Size) {
	a.Shots = make([]*lib.Point, 0)
	a.HitShots = make([]*lib.Point, 0)
	a.SetUpShotPattern(boardSize)
}

func (*AgentSmith) ArrangeShip() {

}

func (a *AgentSmith) GetShot() (point lib.Point) {
	a.Shots = append(a.Shots, &a.FireRandom())
	return *a.Shots[len(a.Shots)-1]
}

func (*AgentSmith) ShotHit(point lib.Point, sunk bool) {

}

func (*AgentSmith) ShotMiss(point lib.Point) {

}

func (a *AgentSmith) SetUpShotPattern(boardSize lib.Size) {
	a.ShotPatterns = make([]*ShotPattern, 0)
	for y := 0; y < boardSize.Height; y++ {
		for x := 0; x < boardSize.Witdh; x++ {
			if (x+y)%2 == 0 {
				a.ShotPatterns = append(a.ShotPatterns, &ShotPattern{&lib.Point{X: x, Y: y}, a.GetScore(x, y)})
			}
		}
	}
}

func (a *AgentSmith) FireDirected(direction int, target lib.Point) (lib.Point) {
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

func (a *AgentSmith) FireAroundPoint(p lib.Point) (lib.Point) {
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

func (a *AgentSmith) FireRandom() (lib.Point) {
	var point lib.Point
	for {
		if len(a.ShotPatterns) > 0 {
			i := rand.Intn(len(a.ShotPatterns))
			point = *a.ShotPatterns[i].Location
			a.ShotPatterns = append(a.ShotPatterns[:i], a.ShotPatterns[i+1:]...)
		} else {
			for {
				point = a.FireAroundPoint(*a.HitShots[0])
				if a.ValidShot(point) {
					a.HitShots = append(a.HitShots[:0], a.HitShots[1:]...)
				}

				if a.ValidShot(point) || len(a.HitShots) <= 0 {
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

func (a *AgentSmith) GetScore(x int, y int) (int) {
	return 0
}

func (a *AgentSmith) ValidShot(p lib.Point) (bool) {
	for i := 0; i < len(a.Shots); i++ {
		s := a.Shots[i]
		if s.X == p.X && s.Y == p.Y {
			return false
		}
	}

	return true
}
