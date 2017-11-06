package strategy

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"encoding/json"
	"github.com/anhtuan29592/battleship-ai/lib/constant"
	"math/rand"
)

type SampleStrategy struct {
	Shots         []lib.Point
	HitShots      []lib.Point
	ShotPatterns  []lib.Point
	BoardSize     lib.Size
	ComboShots    []lib.Point
	CurrentTarget *Target
}

type Target struct {
	Location   lib.Point
	BoardSize  lib.Size
	Neighbors  []lib.Point
	PrevTarget *Target
}

func (s *SampleStrategy) StartGame(boardSize lib.Size) {
	s.Shots = make([]lib.Point, 0)
	s.HitShots = make([]lib.Point, 0)
	s.ComboShots = make([]lib.Point, 0)
	s.BoardSize = boardSize
	s.ShotPatterns = SetUpShotPattern(boardSize)
}

func (s *SampleStrategy) GetShot() (point lib.Point) {
	var shot lib.Point
	if s.CurrentTarget != nil {
		shot = s.CurrentTarget.EvaluateNextShot(s.Shots)
	} else {
		shot = s.FireRandom()
	}
	s.Shots = append(s.Shots, shot)
	return shot
}

func (s *SampleStrategy) ShotHit(point lib.Point, shipPositions []lib.Point) {
	if len(shipPositions) > 0 {
		for i := len(s.ComboShots) - 1; i >= 0; i-- {
			if CheckPointInSlice(shipPositions, s.ComboShots[i]) {
				s.ComboShots = append(s.ComboShots[:i], s.ComboShots[i+1:]...)
			}
		}

		if len(s.ComboShots) > 0 {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.Shots, s.ComboShots[0])
		} else {
			s.CurrentTarget = nil
		}
	} else {
		if s.CurrentTarget == nil {
			s.CurrentTarget = NewTarget(s.Shots, point, s.BoardSize)
		} else {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.Shots, point)
		}
		s.ComboShots = append(s.ComboShots, point)
		s.HitShots = append(s.HitShots, point)
	}
}

func (s *SampleStrategy) ShotMiss(point lib.Point) {

}

func (s *SampleStrategy) GetGameState() lib.GameState {
	gameData, _ := json.Marshal(s)
	return lib.GameState{
		GameStrategy: constant.SAMPLE_STRATEGY,
		GameData:     string(gameData),
	}
}

func (s *SampleStrategy) LoadGameState(state lib.GameState) {
	if state.GameStrategy != constant.SAMPLE_STRATEGY {
		return
	}

	data := state.GameData
	json.Unmarshal([]byte(data), s)
}

func (s *SampleStrategy) FireRandom() lib.Point {
	var point lib.Point
	for {
		if len(s.ShotPatterns) > 0 {
			i := rand.Intn(len(s.ShotPatterns))
			point = s.ShotPatterns[i]
			s.ShotPatterns = append(s.ShotPatterns[:i], s.ShotPatterns[i+1:]...)
		} else {
			for {
				point = s.FireAroundPoint(s.HitShots[0])
				if s.ValidShot(point) {
					break
				}

				s.HitShots = append(s.HitShots[:0], s.HitShots[1:]...)
				if len(s.HitShots) <= 0 {
					break
				}
			}
		}

		if s.ValidShot(point) {
			break
		}
	}
	return point
}

func (s *SampleStrategy) FireAroundPoint(p lib.Point) lib.Point {

	testShot := s.FireDirected(constant.DOWN, p)

	if !s.ValidShot(testShot) {
		testShot = s.FireDirected(constant.UP, p)
	}

	if !s.ValidShot(testShot) {
		testShot = s.FireDirected(constant.RIGHT, p)
	}

	if !s.ValidShot(testShot) {
		testShot = s.FireDirected(constant.LEFT, p)
	}

	return testShot
}

func (s *SampleStrategy) ValidShot(p lib.Point) bool {
	for i := 0; i < len(s.Shots); i++ {
		s := s.Shots[i]
		if s.X == p.X && s.Y == p.Y {
			return false
		}
	}

	return p.ValidInBoard(s.BoardSize)
}

func (*SampleStrategy) FireDirected(direction constant.Direction, target lib.Point) lib.Point {
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

func NewTarget(shots []lib.Point, location lib.Point, boardSize lib.Size) *Target {
	target := Target{Location: location, BoardSize: boardSize, PrevTarget: nil, Neighbors: make([]lib.Point, 0)}
	target.InitNeighbors(shots)
	return &target
}

func (t *Target) Tracking(shots []lib.Point, point lib.Point) *Target {
	newTarget := NewTarget(shots, point, t.BoardSize)
	newTarget.PrevTarget = t
	return newTarget
}

func (t *Target) InitNeighbors(shots []lib.Point) {

	// down
	testPoint := lib.Point{X: t.Location.X, Y: t.Location.Y + 1}
	if testPoint.ValidInBoard(t.BoardSize) && !CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// up
	testPoint = lib.Point{X: t.Location.X, Y: t.Location.Y - 1}
	if testPoint.ValidInBoard(t.BoardSize) && !CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// left
	testPoint = lib.Point{X: t.Location.X - 1, Y: t.Location.Y}
	if testPoint.ValidInBoard(t.BoardSize) && !CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// right
	testPoint = lib.Point{X: t.Location.X + 1, Y: t.Location.Y}
	if testPoint.ValidInBoard(t.BoardSize) && !CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}
}

func (t *Target) EvaluateNextShot(shots []lib.Point) lib.Point {
	node := t
	nodes := make([]lib.Point, 0)
	allNeighbors := make([]lib.Point, 0)
	for {
		nodes = append(nodes, node.Location)
		for i := 0; i < len(node.Neighbors); i++ {
			if !CheckPointInSlice(shots, node.Neighbors[i]) {
				allNeighbors = append(allNeighbors, node.Neighbors[i])
			}
		}
		if node.PrevTarget == nil {
			break
		}
		node = node.PrevTarget
	}

	wholeLineHorizontal := true
	wholeLineVertical := true
	for {
		if node.PrevTarget == nil {
			break
		}

		if node.Location.X != node.PrevTarget.Location.X {
			wholeLineVertical = false
		}

		if node.Location.Y != node.PrevTarget.Location.Y {
			wholeLineHorizontal = false
		}
	}

	nodeCount := len(nodes)

	// carrier ship: next is located
	if nodeCount == 4 {
		if wholeLineHorizontal {
			nodes = SortPoints(nodes, constant.HORIZONTAL, true)
			testNode := lib.Point{X: nodes[0].X + 1, Y: nodes[0].Y - 1}
			if CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			testNode = lib.Point{X: nodes[nodeCount-1].X + 1, Y: nodes[nodeCount-1].Y + 1}
			if CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		if wholeLineVertical {
			nodes = SortPoints(nodes, constant.VERTICAL, true)
			testNode := lib.Point{X: nodes[0].X - 1, Y: nodes[0].Y + 1}
			if CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			testNode = lib.Point{X: nodes[nodeCount-1].X + 1, Y: nodes[nodeCount-1].Y - 1}
			if CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		nodes = SortPoints(nodes, constant.VERTICAL, true)
		testNode := lib.Point{X: nodes[nodeCount-1].X, Y: nodes[nodeCount-1].Y + 1}
		if CheckPointInSlice(allNeighbors, testNode) {
			return testNode
		}

		nodes = SortPoints(nodes, constant.HORIZONTAL, true)
		testNode = lib.Point{X: nodes[nodeCount-1].X + 1, Y: nodes[nodeCount-1].Y}
		if CheckPointInSlice(allNeighbors, testNode) {
			return testNode
		}
	}

	// battle ship or oil rig
	if nodeCount == 3 || nodeCount == 2 {
		// battle ship: next is neighbor of last no first
		if wholeLineHorizontal {
			nodes = SortPoints(nodes, constant.HORIZONTAL, true)
			lastNode := nodes[nodeCount-1]
			testNode := lib.Point{X: lastNode.X + 1, Y: lastNode.Y}
			if testNode.InlineWith(nodes) && CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			firstNode := nodes[0]
			testNode = lib.Point{X: firstNode.X - 1, Y: lastNode.Y}
			if testNode.InlineWith(nodes) && CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		if wholeLineVertical {
			nodes = SortPoints(nodes, constant.VERTICAL, true)
			lastNode := nodes[nodeCount-1]
			testNode := lib.Point{X: lastNode.X, Y: lastNode.Y + 1}
			if testNode.InlineWith(nodes) && CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			firstNode := nodes[0]
			testNode = lib.Point{X: firstNode.X, Y: firstNode.Y - 1}
			if testNode.InlineWith(nodes) && CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		// oil rig
		neighborCount := len(allNeighbors)
		for i := 0; i < neighborCount; i++ {
			for j := i + 1; j < neighborCount; j++ {
				if allNeighbors[i].X == allNeighbors[j].X && allNeighbors[i].Y == allNeighbors[j].Y {
					return allNeighbors[i]
				}
			}
		}
	}

	return allNeighbors[0]
}
