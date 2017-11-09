package util

import (
	"github.com/anhtuan29592/paladin/lib"
	"sort"
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

func CheckCoordinateInSlice(slice []lib.Point, x int, y int) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].X == x && slice[i].Y == y {
			return true
		}
	}

	return false
}

func FindIndex(slice []lib.Point, point lib.Point) int {
	return sort.Search(len(slice), func(i int) bool {
		return slice[i].X == point.X && slice[i].Y == point.Y
	})
}
