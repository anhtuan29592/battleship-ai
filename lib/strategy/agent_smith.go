package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/util"
	"math/rand"
)

type AgentSmith struct {
	Shots        []*lib.Point
	HitShots     []*lib.Point
	ShotPatterns []*ShotPattern
	Endpoints    []*lib.Point
	BoardSize    lib.Size
}

type ShotPattern struct {
	Location *lib.Point
	Score    int
}

func (self *AgentSmith) StartGame(boardSize lib.Size) {
	self.Shots = make([]*lib.Point, 0)
	self.HitShots = make([]*lib.Point, 0)
	self.Endpoints = make([]*lib.Point, 2)
	self.BoardSize = boardSize
	self.SetUpShotPattern(boardSize)
}

func (*AgentSmith) ArrangeShip() {

}

func (self *AgentSmith) GetShot() (point lib.Point) {
	var shot lib.Point
	if self.Endpoints[1] != nil {
		shot = self.FireTargetted()
	} else {
		shot = self.FireRandom()
	}
	self.Shots = append(self.Shots, &shot)
	return shot
}

func (self *AgentSmith) ShotHit(point lib.Point, sunk bool) {
	self.HitShots = append(self.HitShots, &point)
	self.Endpoints[1] = &point
	if self.Endpoints[0] == nil {
		self.Endpoints[0] = &point
	}
}

func (self *AgentSmith) ShotMiss(point lib.Point) {

}

func (self *AgentSmith) SetUpShotPattern(boardSize lib.Size) {
	self.ShotPatterns = make([]*ShotPattern, 0)
	for y := 0; y < boardSize.Height; y++ {
		for x := 0; x < boardSize.Witdh; x++ {
			if (x+y)%2 == 0 {
				self.ShotPatterns = append(self.ShotPatterns, &ShotPattern{&lib.Point{X: x, Y: y}, self.GetScore(x, y)})
			}
		}
	}
}

func (self *AgentSmith) FireDirected(direction int, target lib.Point) (lib.Point) {
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

func (self *AgentSmith) FireAroundPoint(p lib.Point) (lib.Point) {
	testShot := self.FireDirected(constant.UP, p)
	if !self.ValidShot(testShot) {
		testShot = self.FireDirected(constant.LEFT, p)
	}

	if !self.ValidShot(testShot) {
		testShot = self.FireDirected(constant.RIGHT, p)
	}

	if !self.ValidShot(testShot) {
		testShot = self.FireDirected(constant.DOWN, p)
	}

	return testShot
}

func (self *AgentSmith) FireTargetted() (lib.Point) {
	for {
		point := self.FireAroundPoint(lib.Point{X: self.Endpoints[1].X, Y: self.Endpoints[1].Y})
		if self.ValidShot(point) {
			return point
		}

		if self.Endpoints[1] != self.Endpoints[0] {
			self.Endpoints[1] = self.Endpoints[0]
		} else {
			self.Endpoints = make([]*lib.Point, 2)
			break
		}
	}

	return self.FireRandom()
}

func (self *AgentSmith) FireRandom() (lib.Point) {
	var point lib.Point
	for {
		if len(self.ShotPatterns) > 0 {
			i := rand.Intn(len(self.ShotPatterns))
			point = *self.ShotPatterns[i].Location
			self.ShotPatterns = append(self.ShotPatterns[:i], self.ShotPatterns[i+1:]...)
		} else {
			for {
				point = self.FireAroundPoint(*self.HitShots[0])
				if self.ValidShot(point) {
					break
				}

				self.HitShots = append(self.HitShots[:0], self.HitShots[1:]...)
				if len(self.HitShots) <= 0 {
					break
				}
			}
		}

		if self.ValidShot(point) {
			break
		}
	}
	return point
}

func (self *AgentSmith) GetScore(x int, y int) (int) {
	return 0
}

func (self *AgentSmith) ValidShot(p lib.Point) (bool) {
	for i := 0; i < len(self.Shots); i++ {
		s := self.Shots[i]
		if s.X == p.X && s.Y == p.Y {
			return false
		}
	}

	if p.X < 0 || p.Y < 0 {
		return false
	}

	if p.X >= self.BoardSize.Witdh ||  p.Y >= self.BoardSize.Height {
		return false
	}

	return true
}
