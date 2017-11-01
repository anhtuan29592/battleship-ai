package service

import (
	"github.com/anhtuan29592/battleship-ai/lib/ship"
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/lib/util"
	"math/rand"
	"github.com/anhtuan29592/battleship-ai/lib"
)

type GameMetric struct {
}

func (g *GameMetric) CreateShips(quantities []*domain.ShipQuantity) []*ship.Ship {
	ships := make([]*ship.Ship, 0);
	for i := 0; i < len(quantities); i++ {
		shipQuantity := quantities[i]
		for j := 0; j < shipQuantity.Quantity; j++ {
			switch constant.ShipType(shipQuantity.Type) {
			case constant.CARRIER:
				ships = append(ships, &ship.Ship{ShipAction: &ship.CarrierShip{}})
				break
			case constant.BATTLE_SHIP:
				ships = append(ships, &ship.Ship{ShipAction: &ship.BattleShip{}})
				break
			case constant.OIL_RIG:
				ships = append(ships, &ship.Ship{ShipAction: &ship.OilRigShip{}})
				break
			case constant.CRUISER:
				ships = append(ships, &ship.Ship{ShipAction: &ship.CruiserShip{}})
				break
			case constant.DESTROYER:
				ships = append(ships, &ship.Ship{ShipAction: &ship.DestroyerShip{}})
				break
			default:
				break
			}
		}
	}
	return ships
}

func (g *GameMetric) ArrangeShips(boardSize lib.Size, ships []*ship.Ship) []*ship.Ship {
	if g.Validate(boardSize, ships) {
		return ships
	}

	for i := 0; i < len(ships); i++ {
		x := rand.Intn(boardSize.Witdh - 1)
		y := rand.Intn(boardSize.Height - 1)
		orientation := constant.Orientation(rand.Intn(2))
		ships[i].ShipAction.UpdateLocation(orientation, &lib.Point{X: x, Y: y})
	}

	return g.ArrangeShips(boardSize, ships)
}

func (g *GameMetric) Validate(boardSize lib.Size, ships []*ship.Ship) bool {
	for i := 0; i < len(ships); i++ {
		if ships[i].ShipAction.IsValid(boardSize) {
			return false
		}
	}

	for i := 0; i < len(ships); i++ {
		for j := i + 1; j < len(ships); j++ {
			if ships[i].ShipAction.ConflictWith(ships[j]) {
				return false
			}
		}
	}

	return true
}
