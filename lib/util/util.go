package util

import (
	"fmt"
	"sort"

	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/constant"
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

func PrintPoints(boardSize lib.Size, points []lib.Point) {
	for r := 0; r < boardSize.Height; r++ {
		fmt.Print("|")
		for c := 0; c < boardSize.Width; c++ {
			printed := false
			for i := 0; i < len(points); i++ {
				if points[i].X == c && points[i].Y == r {
					fmt.Print("x")
					printed = true
					break
				}
			}
			if !printed {
				fmt.Print(" ")
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintPriorityPoints(boardSize lib.Size, points []lib.PriorityPoint) {
	for r := 0; r < boardSize.Height; r++ {
		fmt.Print("|")
		for c := 0; c < boardSize.Width; c++ {
			printed := false
			for i := 0; i < len(points); i++ {
				if points[i].Location.X == c && points[i].Location.Y == r {
					fmt.Printf("%2d", points[i].Score)
					printed = true
					break
				}
			}
			if !printed {
				fmt.Print("  ")
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println()
}

func SortPoints(points []lib.Point, orientation constant.Orientation, isAscending bool) []lib.Point {
	if orientation == constant.HORIZONTAL {
		sort.Slice(points, func(i, j int) bool {
			if isAscending {
				return points[i].X < points[j].X
			}
			return points[i].X >= points[j].X
		})
	}

	if orientation == constant.VERTICAL {
		sort.Slice(points, func(i, j int) bool {
			if isAscending {
				return points[i].Y < points[j].Y
			}
			return points[i].Y >= points[j].Y
		})
	}

	return points
}

func GetDirection(p1 lib.Point, p2 lib.Point) constant.Direction {
	if p1.Y == p2.Y {
		if p1.X >= p2.X {
			return constant.RIGHT
		} else {
			return constant.LEFT
		}
	}

	// Vertical
	if p1.X == p2.X {
		if p1.Y >= p2.Y {
			return constant.DOWN
		} else {
			return constant.UP
		}
	}

	return constant.UP
}
