package lib

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Point) ValidInBoard(boardSize Size) bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}

	if p.X >= boardSize.Width || p.Y >= boardSize.Height {
		return false
	}

	return true
}

func (p *Point) InlineWith(points []Point) bool {
	horizontal := true
	vertical := true
	for i := 0; i < len(points); i++ {
		if points[i].X != p.X {
			vertical = false
			break
		}
	}

	for i := 0; i < len(points); i++ {
		if points[i].Y != p.Y {
			horizontal = false
			break
		}
	}

	return horizontal || vertical
}
