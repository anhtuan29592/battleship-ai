package lib

type Point struct {
	x int64 `json:"x"`
	y int64 `json:"y"`
}

func (p *Point) get_x() int64 {
	return p.x
}

func (p *Point) get_y() int64 {
	return p.y
}