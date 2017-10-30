package opponent


type Opponent interface {
	ArrangeShip()
	GetShot() (point *Point)
	ShotHit(point *Point, sunk bool)
	ShotMiss(point *Point)
}