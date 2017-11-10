package strategy

import (
	"encoding/json"
	"math/rand"
	"sort"

	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
	"github.com/anhtuan29592/paladin/lib/ship"
	"github.com/anhtuan29592/paladin/lib/util"
	"log"
)

type SampleStrategy struct {
	BoardSize          lib.Size
	Shots              []lib.Point
	HitShots           []lib.Point
	ShotPatterns       []lib.PriorityPoint
	ComboShots         []lib.Point
	InvalidShots       []lib.Point
	CurrentTarget      *Target
	ShipTypeCount      map[constant.ShipType]int
	MissCrossLineCount int
}

type Target struct {
	Location   lib.Point
	BoardSize  lib.Size
	Neighbors  []lib.Point
	PrevTarget *Target
}

type PriorityPoint struct {
	Location lib.Point
	Score    int
}

func (s *SampleStrategy) StartGame(boardSize lib.Size, ships []ship.Ship) {
	s.MissCrossLineCount = 0
	s.BoardSize = boardSize
	s.Shots = make([]lib.Point, 0)
	s.HitShots = make([]lib.Point, 0)
	s.ComboShots = make([]lib.Point, 0)
	s.InvalidShots = make([]lib.Point, 0)
	s.ShipTypeCount = make(map[constant.ShipType]int)

	s.ShotPatterns = SetUpShotPattern(boardSize)
	util.PrintPriorityPoints(boardSize, s.ShotPatterns)

	for i := 0; i < len(ships); i++ {
		s.ShipTypeCount[ships[i].GetType()]++
	}
}

func (s *SampleStrategy) GetShot() (point lib.Point) {
	var shot lib.Point
	if s.CurrentTarget != nil {
		retryCount := 0
		for {

			shot = s.CurrentTarget.EvaluateNextShot(s.Shots, s.ShipTypeCount)
			if !util.CheckPointInSlice(s.Shots, shot) && point.ValidInBoard(s.BoardSize) {
				return shot
			}

			if retryCount > 100 {
				break
			}
			retryCount++
		}
	}

	log.Printf("Random")
	return s.FireRandom()
}

func (s *SampleStrategy) ShotHit(point lib.Point, shipType string, shipPositions []lib.Point) {
	s.Shots = append(s.Shots, point)
	s.HitShots = append(s.HitShots, point)
	s.InvalidShots = append(s.InvalidShots, point)
	if (0 <= point.X && point.X < 2) || (s.BoardSize.Width-2 <= point.X && point.X < s.BoardSize.Width) || (0 <= point.Y && point.Y < 2) || (s.BoardSize.Height-2 <= point.Y && point.Y < s.BoardSize.Height) {
		s.MissCrossLineCount--
		log.Printf("Miss cross line count %d", s.MissCrossLineCount)
		if s.MissCrossLineCount < 0 {
			s.MissCrossLineCount = 0
		}
	}

	if len(shipPositions) > 0 {
		s.ShipTypeCount[constant.ShipType(shipType)]--

		for i := len(s.ComboShots) - 1; i >= 0; i-- {
			if util.CheckPointInSlice(shipPositions, s.ComboShots[i]) {
				s.ComboShots = append(s.ComboShots[:i], s.ComboShots[i+1:]...)
			}
		}

		if len(s.ComboShots) > 0 {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.Shots, s.ComboShots[0])
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
			for j := 0; j < len(s.ShotPatterns); j++ {
				if s.ShotPatterns[j].Location.X == testPoint.X && s.ShotPatterns[j].Location.Y == testPoint.Y {
					s.ShotPatterns = append(s.ShotPatterns[:j], s.ShotPatterns[j+1:]...)
					break
				}
			}

			// down
			testPoint = lib.Point{X: tmp.X, Y: tmp.Y + 1}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.ShotPatterns); j++ {
				if s.ShotPatterns[j].Location.X == testPoint.X && s.ShotPatterns[j].Location.Y == testPoint.Y {
					s.ShotPatterns = append(s.ShotPatterns[:j], s.ShotPatterns[j+1:]...)
					break
				}
			}

			// left
			testPoint = lib.Point{X: tmp.X - 1, Y: tmp.Y}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.ShotPatterns); j++ {
				if s.ShotPatterns[j].Location.X == testPoint.X && s.ShotPatterns[j].Location.Y == testPoint.Y {
					s.ShotPatterns = append(s.ShotPatterns[:j], s.ShotPatterns[j+1:]...)
					break
				}
			}

			// right
			testPoint = lib.Point{X: tmp.X + 1, Y: tmp.Y}
			if !util.CheckPointInSlice(s.InvalidShots, testPoint) {
				s.InvalidShots = append(s.InvalidShots, testPoint)
			}
			for j := 0; j < len(s.ShotPatterns); j++ {
				if s.ShotPatterns[j].Location.X == testPoint.X && s.ShotPatterns[j].Location.Y == testPoint.Y {
					s.ShotPatterns = append(s.ShotPatterns[:j], s.ShotPatterns[j+1:]...)
					break
				}
			}
		}
	} else {
		if s.CurrentTarget == nil {
			s.CurrentTarget = NewTarget(s.Shots, point, s.BoardSize)
		} else {
			s.CurrentTarget = s.CurrentTarget.Tracking(s.Shots, point)
		}
		s.ComboShots = append(s.ComboShots, point)
	}

}

func (s *SampleStrategy) ShotMiss(point lib.Point) {
	s.InvalidShots = append(s.InvalidShots, point)
	s.Shots = append(s.Shots, point)

	if (0 <= point.X && point.X < 2) || (s.BoardSize.Width-2 <= point.X && point.X < s.BoardSize.Width) || (0 <= point.Y && point.Y < 2) || (s.BoardSize.Height-2 <= point.Y && point.Y < s.BoardSize.Height) {
		s.MissCrossLineCount++
		log.Printf("Miss cross line count %d", s.MissCrossLineCount)
		if s.MissCrossLineCount > 6 {
			s.MissCrossLineCount = 6
		}
	}
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
			s.UpdatePriority()
			maxScore := s.ShotPatterns[0].Score
			randMax := make([]lib.PriorityPoint, 0)
			for i := 0; i < len(s.ShotPatterns); i++ {
				if s.ShotPatterns[i].Score == maxScore {
					randMax = append(randMax, s.ShotPatterns[i])
				}
			}

			point = randMax[rand.Intn(len(randMax))].Location
			for i := 0; i < len(s.ShotPatterns); i++ {
				if s.ShotPatterns[i].Location.X == point.X && s.ShotPatterns[i].Location.Y == point.Y {
					s.ShotPatterns = append(s.ShotPatterns[:i], s.ShotPatterns[i+1:]...)
					break
				}
			}

			util.PrintPriorityPoints(s.BoardSize, s.ShotPatterns)
		} else {
			for {
				i := rand.Intn(len(s.HitShots))
				point = s.FireAroundPoint(s.HitShots[i])
				if s.ValidShot(point) {
					break
				}

				s.HitShots = append(s.HitShots[:i], s.HitShots[i+1:]...)
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

	// up
	testPoint := lib.Point{X: t.Location.X, Y: t.Location.Y - 1}
	if testPoint.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testPoint) {
		t.Neighbors = append(t.Neighbors, testPoint)
	}

	// down
	testPoint = lib.Point{X: t.Location.X, Y: t.Location.Y + 1}
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
	util.PrintPoints(t.BoardSize, allNeighbors)

	wholeLineHorizontal := true
	wholeLineVertical := true
	node = t
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

		node = node.PrevTarget
	}

	nodeCount := len(nodes)

	// carrier ship: next is located
	if nodeCount >= 4 && shipTypeCount[constant.CARRIER] > 0 {
		log.Printf("Check carrier size %d", nodeCount)
		if wholeLineHorizontal {
			nodes = util.SortPoints(nodes, constant.HORIZONTAL, true)
			testNode := lib.Point{X: nodes[0].X + 1, Y: nodes[0].Y - 1}
			if util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		if wholeLineVertical {
			nodes = util.SortPoints(nodes, constant.VERTICAL, true)
			testNode := lib.Point{X: nodes[0].X - 1, Y: nodes[0].Y + 1}
			if util.CheckPointInSlice(allNeighbors, testNode) {
				return testNode
			}
		}

		countNode := make(map[int][]lib.Point)
		idx := 0
		// Scan from left to right to get pattern
		for i := 0; i < t.BoardSize.Width; i++ {
			found := false
			for j := 0; j < nodeCount; j++ {
				if nodes[j].X == i {
					countNode[idx] = append(countNode[idx], nodes[j])
					found = true
				}
			}
			if found {
				idx++
			}
		}

		if len(countNode[0]) == 1 && len(countNode[1]) == 3 {
			sameVertical := util.SortPoints(countNode[1], constant.VERTICAL, true)
			//    CV
			// CV CV
			//    CV
			//    XX
			if sameVertical[1].Y == countNode[0][0].Y {
				testNode := lib.Point{X: sameVertical[0].X, Y: sameVertical[0].Y + 3}
				if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
					return testNode
				}
			}

			//    XX
			// CV CV
			//    CV
			//    CV
			if sameVertical[0].Y == countNode[0][0].Y {
				testNode := lib.Point{X: sameVertical[0].X, Y: sameVertical[0].Y - 1}
				if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
					return testNode
				}
			}
		}

		if len(countNode[0]) == 1 && len(countNode[1]) == 2 && len(countNode[2]) == 1 {
			testNode := lib.Point{X: countNode[2][0].X + 1, Y: countNode[2][0].Y}
			if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
				return testNode
			}
		}

		if len(countNode[0]) == 2 && len(countNode[1]) == 1 && len(countNode[2]) == 1 {
			testNode := lib.Point{X: countNode[1][0].X - 2, Y: countNode[1][0].Y}
			if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
				return testNode
			}
		}
	}

	// battle ship or oil rig
	if nodeCount >= 3 {
		if shipTypeCount[constant.OIL_RIG] > 0 {
			log.Printf("Check oil rig size %d", nodeCount)
			// oil rig
			neighborCount := len(allNeighbors)
			for i := 0; i < neighborCount-1; i++ {
				for j := i + 1; j < neighborCount; j++ {
					if allNeighbors[i].X == allNeighbors[j].X && allNeighbors[i].Y == allNeighbors[j].Y {
						return allNeighbors[i]
					}
				}
			}
		}

		// battle ship: next is neighbor of last no first
		if shipTypeCount[constant.BATTLE_SHIP] > 0 {
			log.Printf("Check battle size %d", nodeCount)
			if wholeLineHorizontal {
				nodes = util.SortPoints(nodes, constant.HORIZONTAL, true)
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
				nodes = util.SortPoints(nodes, constant.VERTICAL, true)
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

		if shipTypeCount[constant.CARRIER] > 0 {
			log.Printf("Check carrier size %d", nodeCount)
			countNode := make(map[int][]lib.Point)
			idx := 0
			// Scan from left to right to get pattern
			for i := 0; i < t.BoardSize.Width; i++ {
				found := false
				for j := 0; j < nodeCount; j++ {
					if nodes[j].X == i {
						countNode[idx] = append(countNode[idx], nodes[j])
						found = true
					}
				}
				if found {
					idx++
				}
			}

			if len(countNode[0]) == 1 && len(countNode[1]) == 2 {
				sameVertical := util.SortPoints(countNode[1], constant.VERTICAL, true)
				//    XX
				// CV CV
				//    CV
				//    XX
				if sameVertical[0].Y == countNode[0][0].Y {
					//    XX
					// CV CV
					//    CV
					testNode := lib.Point{X: sameVertical[0].X, Y: sameVertical[0].Y - 1}
					if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
						return testNode
					}

					// CV CV
					//    CV
					//    XX
					testNode = lib.Point{X: sameVertical[0].X, Y: sameVertical[0].Y + 2}
					if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
						return testNode
					}
				}

				//    CV
				// CV CV XX
				//    XX
				if sameVertical[1].Y == countNode[0][0].Y {
					//    CV
					// CV CV
					//    XX
					testNode := lib.Point{X: sameVertical[0].X, Y: sameVertical[0].Y + 2}
					if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
						return testNode
					}
					//    CV
					// CV CV XX
					testNode = lib.Point{X: sameVertical[0].X + 1, Y: sameVertical[0].Y + 1}
					if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
						return testNode
					}
				}
			}

			if len(countNode[0]) == 2 && len(countNode[1]) == 1 {
				sameVertical := util.SortPoints(countNode[0], constant.VERTICAL, true)
				//    CV
				// XX CV CV
				testNode := lib.Point{X: sameVertical[0].X - 1, Y: sameVertical[0].Y + 1}
				if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
					return testNode
				}

				// CV
				// CV CV XX
				testNode = lib.Point{X: sameVertical[0].X + 2, Y: sameVertical[0].Y + 1}
				if testNode.ValidInBoard(t.BoardSize) && !util.CheckPointInSlice(shots, testNode) {
					return testNode
				}
			}
		}
	}

	// cruiser
	if nodeCount >= 2 {
		if shipTypeCount[constant.CRUISER] > 0 || shipTypeCount[constant.BATTLE_SHIP] > 0 {
			log.Printf("Check cruiser or battle ship size %d", nodeCount)
			if wholeLineHorizontal {
				nodes = util.SortPoints(nodes, constant.HORIZONTAL, true)
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
				nodes = util.SortPoints(nodes, constant.VERTICAL, true)
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

		if shipTypeCount[constant.OIL_RIG] > 0 {
			log.Printf("Check oil rig size %d", nodeCount)
			if wholeLineHorizontal {
				nodes = util.SortPoints(nodes, constant.HORIZONTAL, true)
				testNode := lib.Point{X: nodes[0].X, Y: nodes[0].Y - 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[0].X, Y: nodes[0].Y + 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X, Y: nodes[1].Y - 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X, Y: nodes[1].Y + 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}

			if wholeLineVertical {
				nodes = util.SortPoints(nodes, constant.VERTICAL, true)
				testNode := lib.Point{X: nodes[0].X - 1, Y: nodes[0].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[0].X + 1, Y: nodes[0].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X - 1, Y: nodes[1].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X + 1, Y: nodes[1].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}
		}

		if shipTypeCount[constant.CARRIER] > 0 {
			log.Printf("Check carrier size %d", nodeCount)
			if wholeLineHorizontal {
				nodes = util.SortPoints(nodes, constant.HORIZONTAL, true)
				testNode := lib.Point{X: nodes[0].X, Y: nodes[0].Y - 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X, Y: nodes[1].Y - 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X + 1, Y: nodes[1].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}

			if wholeLineVertical {
				nodes = util.SortPoints(nodes, constant.VERTICAL, true)
				testNode := lib.Point{X: nodes[0].X + 1, Y: nodes[0].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[1].X + 1, Y: nodes[1].Y}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}

				testNode = lib.Point{X: nodes[0].X, Y: nodes[0].Y - 1}
				if util.CheckPointInSlice(allNeighbors, testNode) {
					return testNode
				}
			}
		}
	}

	log.Printf("Shot direction with size %d", nodeCount)

	if len(allNeighbors) > 0 {
		return allNeighbors[0]
	}
	return lib.Point{X: -1, Y: -1}
}

func (s *SampleStrategy) GetScore(point lib.PriorityPoint) int {
	score := 0
	// space to left
	x := point.Location.X - 1
	y := point.Location.Y
	for {
		if x < 0 {
			break
		}

		if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
			break
		}

		score++
		x--
	}

	// space to right
	x = point.Location.X + 1
	y = point.Location.Y
	for {
		if x >= s.BoardSize.Width {
			break
		}

		if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
			break
		}

		score++
		x++
	}

	// space to top
	x = point.Location.X
	y = point.Location.Y - 1
	for {
		if y < 0 {
			break
		}

		if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
			break
		}

		score++
		y--
	}

	// space to bottom
	x = point.Location.X
	y = point.Location.Y - 1
	for {
		if y >= s.BoardSize.Height {
			break
		}

		if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
			break
		}

		score++
		y++
	}

	if s.MissCrossLineCount <= 6 {
		// space to top - left
		x = point.Location.X - 1
		y = point.Location.Y - 1
		for {
			if x < 0 && y < 0 {
				break
			}
			if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
				break
			}

			score++
			if x >= 0 {
				x--
			}

			if y >= 0 {
				y--
			}
		}

		// space to top - right
		x = point.Location.X + 1
		y = point.Location.Y - 1
		for {
			if x >= s.BoardSize.Width && y < 0 {
				break
			}
			if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
				break
			}

			score++
			if x < s.BoardSize.Width {
				x++
			}
			if y >= 0 {
				y--
			}
		}

		// space to bottom - left
		x = point.Location.X - 1
		y = point.Location.Y + 1
		for {
			if x < 0 && y >= s.BoardSize.Height {
				break
			}
			if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
				break
			}

			score++
			if x >= 0 {
				x--
			}

			if y < s.BoardSize.Height {
				y++
			}
		}

		// space to bottom - right
		x = point.Location.X + 1
		y = point.Location.Y + 1
		for {
			if x >= s.BoardSize.Width && y >= s.BoardSize.Height {
				break
			}
			if util.CheckCoordinateInSlice(s.InvalidShots, x, y) {
				break
			}

			score++
			if x < s.BoardSize.Width {
				x++
			}

			if y < s.BoardSize.Height {
				y++
			}
		}
	}

	return score
}

func (s *SampleStrategy) UpdatePriority() {
	for i := 0; i < len(s.ShotPatterns); i++ {
		s.ShotPatterns[i].Score = s.GetScore(s.ShotPatterns[i])
	}

	sort.Slice(s.ShotPatterns, func(i, j int) bool {
		return s.ShotPatterns[i].Score > s.ShotPatterns[j].Score
	})
}
