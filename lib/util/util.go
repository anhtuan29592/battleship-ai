package util

import (
	"github.com/anhtuan29592/paladin/lib"
)

func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}


func CheckPointInSlice(slice []lib.Point, point lib.Point) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].X == point.X && slice[i].Y == point.Y {
			return true
		}
	}

	return false
}
