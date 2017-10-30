package lib


type BattleShipOpponent interface {
	arrange_ship()
	get_shot() (point *Point)
	shot_hit(point *Point, sunk bool)
	shot_miss(point *Point)
}