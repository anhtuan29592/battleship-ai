package strategy

import (
	"github.com/anhtuan29592/paladin/lib"
	"encoding/json"
	"github.com/anhtuan29592/paladin/lib/constant"
	"math/rand"
	"log"
	"github.com/anhtuan29592/paladin/lib/util"
	"github.com/anhtuan29592/paladin/lib/ship"
)

type SampleStrategy struct {
	BoardSize     lib.Size
	Shots         []lib.Point
	HitShots      []lib.Point
	ShotPatterns  []lib.Point
	PriorityShots []lib.Point
	ComboShots    []lib.Point
	InvalidShots  []lib.Point
	CurrentTarget *Target
	ShipTypeCount map[constant.ShipType]int
}

type Target struct {
	Location   lib.Point
	BoardSize  lib.Size
	Neighbors  []lib.Point
	PrevTarget *Target
}

func (s *SampleStrategy) StartGame(boardSize lib.Size, ships []ship.Ship) {
	s.BoardSize = boardSize
	s.Shots = make([]lib.Point, 0)
	s.HitShots = make([]lib.Point, 0)
	s.ComboShots = make([]lib.Point, 0)
	s.InvalidShots = make([]lib.Point, 0)
	s.ShipTypeCount = make(map[constant.ShipType]int)

	s.ShotPatterns = SetUpShotPattern(boardSize)
	s.PriorityShots = SetUpPriorityShots(boardSize, s.ShotPatterns)

	PrintPoints(boardSize, s.PriorityShots)

	for i := 0; i < len(ships); i++ {
		s.ShipTypeCount[ships[i].GetType()]++
	}

	log.Print(s.ShipTypeCount)
}

func (s *SampleStrategy) GetShot() (point lib.Point) {
	var shot lib.Point
	found := false
	if s.CurrentTarget != nil {
		retryCount := 0
		for {
			shot = s.CurrentTarget.EvaluateNextShot(s.InvalidShots, s.ShipTypeCount)
			if !util.CheckPointInSlice(s.InvalidShots, shot) && point.ValidInBoard(s.BoardSize) {
				found = true
				break
			}

			shot = s.CurrentTarget.EvaluateNextShot(s.Shots, s.ShipTypeCount)
			if !util.CheckPointInSlice(s.Shots, shot) && point.ValidInBoard(s.BoardSize) {
				found = true
				break
			}

			if retryCount > 100 {
				break
			}
			retryCount++
		}
	}

	if !found {
		shot = s.FireRandom()
	}

	return shot
}

func (s *SampleStrategy) ShotHit(point lib.Point, shipType string, shipPositions []lib.Point) {
	s.Shots = append(s.Shots, point)
	s.HitShots = append(s.HitShots, point)
	s.InvalidShots = append(s.InvalidShots, point)

	if len(shipPositions) > 0 {
		s.ShipTypeCount[constant.ShipType(shipType)]--

		for i := len(s.ComboShots) - 1; i >= 0; i-- {
			if util.CheckPointInSlice(shipPositions, s.ComboShots[i]) {
				s.ComboShots = append(s.ComboShots[:i], s.ComboShots[i+1:]...)
			}
		}

		if len(s.ComboShots) > 0 {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.InvalidShots, s.ComboShots[0])
		} else {
			s.CurrentTarget = nil
		}

		// clear around priority
		for i := 0; i < len(shipPositions); i++ {
			tmp := shipPositions[i]
			// up
			testPoint := lib.Point{X: tmp.X, Y: tmp.Y - 1}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.PriorityShots); j++ {
				if s.PriorityShots[j].X == testPoint.X && s.PriorityShots[j].Y == testPoint.Y {
					s.PriorityShots = append(s.PriorityShots[:j], s.PriorityShots[j+1:]...)
					break
				}
			}

			// down
			testPoint = lib.Point{X: tmp.X, Y: tmp.Y + 1}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.PriorityShots); j++ {
				if s.PriorityShots[j].X == testPoint.X && s.PriorityShots[j].Y == testPoint.Y {
					s.PriorityShots = append(s.PriorityShots[:j], s.PriorityShots[j+1:]...)
					break
				}
			}

			// left
			testPoint = lib.Point{X: tmp.X - 1, Y: tmp.Y}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.PriorityShots); j++ {
				if s.PriorityShots[j].X == testPoint.X && s.PriorityShots[j].Y == testPoint.Y {
					s.PriorityShots = append(s.PriorityShots[:j], s.PriorityShots[j+1:]...)
					break
				}
			}

			// right
			testPoint = lib.Point{X: tmp.X + 1, Y: tmp.Y}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.PriorityShots); j++ {
				if s.PriorityShots[j].X == testPoint.X && s.PriorityShots[j].Y == testPoint.Y {
					s.PriorityShots = append(s.PriorityShots[:j], s.PriorityShots[j+1:]...)
					break
				}
			}
		}
	} else {
		if s.CurrentTarget == nil {
			s.CurrentTarget = NewTarget(s.InvalidShots, point, s.BoardSize)
		} else {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.InvalidShots, point)
		}
		s.ComboShots = append(s.ComboShots, point)
	}

}

func (s *SampleStrategy) ShotMiss(point lib.Point) {
	s.InvalidShots = append(s.InvalidShots, point)
	s.Shots = append(s.Shots, point)
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
		if len(s.PriorityShots) > 0 {
			i := rand.Intn(len(s.PriorityShots))
			point = s.PriorityShots[i]
			s.PriorityShots = append(s.PriorityShots[:i], s.PriorityShots[i+1:]...)
			PrintPoints(s.BoardSize, s.PriorityShots)
		} else if len(s.ShotPatterns) > 0 {
			i := rand.Intn(len(s.ShotPatterns))
			point = s.ShotPatterns[i]
			s.ShotPatterns = append(s.ShotPatterns[:i], s.ShotPatterns[i+1:]...)
			PrintPoints(s.BoardSize, s.ShotPatterns)
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
	if testPoint.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// up
	testPoint = lib.Point{X: t.Location.X, Y: t.Location.Y - 1}
	if testPoint.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// left
	testPoint = lib.Point{X: t.Location.X - 1, Y: t.Location.Y}
	if testPoint.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// right
	testPoint = lib.Point{X: t.Location.X + 1, Y: t.Location.Y}
	if testPoint.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}
}

func (t *Target) EvaluateNextShot(shots []lib.Point, shipTypeCount map[constant.ShipType]int) lib.Point {
	node := t
	nodes := make([]lib.Point, 0)
	allNeighbors := make([]lib.Point, 0)
	for {
		nodes = append(nodes, node.Location)
		for i := 0; i < len(node.Neighbors); i++ {
			if !util.CheckPointInSlice(shots, node.Neighbors[i]) {
				allNeighbors = append(allNeighbors, node.Neighbors[i])
			}
		}
		if node.PrevTarget == nil {
			break
		}
		node = node.PrevTarget
	}
	PrintPoints(t.BoardSize, allNeighbors)

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
	if nodeCount >= 4 && shipTypeCount[constant.CARRIER] > 0 {
		if wholeLineHorizontal {
			nodes = SortPoints(nodes, constant.HORIZONTAL, true)
			testNode := lib.Point{X: nodes[0].X + 1, Y: nodes[0].Y - 1}
			if util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		if wholeLineVertical {
			nodes = SortPoints(nodes, constant.VERTICAL, true)
			testNode := lib.Point{X: nodes[0].X - 1, Y: nodes[0].Y + 1}
			if util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}
	}

	// battle ship or oil rig
	if nodeCount >= 3 {
		// battle ship: next is neighbor of last no first
		if shipTypeCount[constant.BATTLE_SHIP] > 0 {
			if wholeLineHorizontal {
				nodes = SortPoints(nodes, constant.HORIZONTAL, true)
				lastNode := nodes[nodeCount-1]
				testNode := lib.Point{X: lastNode.X + 1, Y: lastNode.Y}
				if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				firstNode := nodes[0]
				testNode = lib.Point{X: firstNode.X - 1, Y: lastNode.Y}
				if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}

			if wholeLineVertical {
				nodes = SortPoints(nodes, constant.VERTICAL, true)
				lastNode := nodes[nodeCount-1]
				testNode := lib.Point{X: lastNode.X, Y: lastNode.Y + 1}
				if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				firstNode := nodes[0]
				testNode = lib.Point{X: firstNode.X, Y: firstNode.Y - 1}
				if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}
		} else if shipTypeCount[constant.OIL_RIG] > 0 {
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
	}

	// battle ship or oil rig
	if nodeCount >= 2 && shipTypeCount[constant.CRUISER] > 0 {
		if wholeLineHorizontal {
			nodes = SortPoints(nodes, constant.HORIZONTAL, true)
			lastNode := nodes[nodeCount-1]
			testNode := lib.Point{X: lastNode.X + 1, Y: lastNode.Y}
			if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			firstNode := nodes[0]
			testNode = lib.Point{X: firstNode.X - 1, Y: lastNode.Y}
			if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		if wholeLineVertical {
			nodes = SortPoints(nodes, constant.VERTICAL, true)
			lastNode := nodes[nodeCount-1]
			testNode := lib.Point{X: lastNode.X, Y: lastNode.Y + 1}
			if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}

			firstNode := nodes[0]
			testNode = lib.Point{X: firstNode.X, Y: firstNode.Y - 1}
			if testNode.InlineWith(nodes) && util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}
	}

	return allNeighbors[0]
}
