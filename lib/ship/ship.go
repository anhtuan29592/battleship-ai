package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/constant"
	"github.com/anhtuan29592/battleship-ai/lib/util"
)

type ShipAction interface {
	GetType() constant.ShipType
	GetSize(orientation constant.Orientation) lib.Size
	GetPositions(location lib.Point, orientation constant.Orientation) []lib.Point
}

type Ship struct {
	Location    lib.Point
	Orientation constant.Orientation
	Action      ShipAction
}

func (s *Ship) GetType() constant.ShipType {
	return s.Action.GetType()
}

func (s *Ship) GetSize() lib.Size {
	return s.Action.GetSize(s.Orientation)
}

func (s *Ship) GetPositions() []lib.Point {
	return s.Action.GetPositions(s.Location, s.Orientation)
}

func (s *Ship) ConflictWith(other Ship) bool {
	otherPositions := other.Zoom()
	myPositions := s.Zoom()

	for i := 0; i < len(otherPositions); i++ {
		for j := 0; j < len(myPositions); j++ {
			if otherPositions[i].X == myPositions[j].X && otherPositions[i].Y == myPositions[j].Y {
				return true
			}
		}
	}

	return false
}

func (s *Ship) Touch(other Ship, touchDistance int) bool {
	otherPositions := other.Zoom()
	myPositions := s.Zoom()

	for i := 0; i < len(otherPositions); i++ {
		for j := 0; j < len(myPositions); j++ {
			if util.Abs(otherPositions[i].X-myPositions[j].X) < touchDistance || util.Abs(otherPositions[i].Y-myPositions[j].Y) < touchDistance {
				return false
			}
		}
	}

	return true
}

func (s *Ship) UpdateLocation(orientation constant.Orientation, point lib.Point) {
	s.Location = point
	s.Orientation = orientation
}

func (s *Ship) IsValid(boardSize lib.Size) bool {
	if &s.Location == nil {
		return false
	}

	if s.Location.X < 0 || s.Location.Y < 0 {
		return false
	}

	if s.Location.X >= boardSize.Width || s.Location.Y >= boardSize.Height {
		return false
	}

	size := s.GetSize()

	if s.Location.X+size.Width > boardSize.Width || s.Location.Y+size.Height > boardSize.Height {
		return false
	}

	return true
}

func (s *Ship) Zoom(boardSize lib.Size) []lib.Point {
	sType := s.GetType()
	positions := s.GetPositions()

	if sType == constant.CARRIER || sType == constant.BATTLE_SHIP {
		return positions
	}

	// zoom size
	size := s.GetSize()
	size = lib.Size{Width: size.Width + 2, Height: size.Height + 2}

	positions := make([]lib.Point, 0)
	for r := 0; r < size.Height; r++ {
		for c := 0; c < size.Width; c++ {
			point := lib.Point{X: s.Location.X + c, Y: s.Location.Y + r}
			if point.ValidInBoard(boardSize) {
				positions = append(positions, point)
			}
		}
	}
	return positions
}
