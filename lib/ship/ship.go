package ship

import (
	"github.com/anhtuan29592/battleship-ai/lib"
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
	otherPositions := other.GetPositions()
	myPositions := s.GetPositions()

	for i := 0; i < len(otherPositions); i++ {
		for j := 0; j < len(myPositions); j++ {
			if otherPositions[i].X == myPositions[j].X && otherPositions[i].Y == myPositions[j].Y {
				return true
			}
		}
	}

	return false
}

func (s *Ship) Near(other Ship) bool {
	otherSize := other.GetSize()
	mySize := s.GetSize()

	// UPPER
	if other.Location.Y < s.Location.Y {
		return !(s.Location.Y - (other.Location.Y + otherSize.Height) <= 3)
	}

	// LOWER
	if other.Location.Y > s.Location.Y {
		return !(other.Location.Y - (s.Location.Y + mySize.Height) <= 3)
	}

	// LEFT
	if other.Location.X < s.Location.X {
		return !(s.Location.X - (other.Location.X + otherSize.Width) <= 3)
	}

	// RIGHT
	if other.Location.X > s.Location.X {
		return !(other.Location.X - (s.Location.X + mySize.Width) <= 3)
	}

	return false
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
