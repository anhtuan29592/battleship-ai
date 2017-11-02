package service

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/anhtuan29592/battleship-ai/lib/ship"
	"github.com/anhtuan29592/battleship-ai/lib/util"
	"github.com/anhtuan29592/battleship-ai/lib/strategy"
)

type GameMetric struct {
	GameAI strategy.GameAI
}

func NewGameMetric() GameMetric {
	asm := &strategy.AgentSmith{}
	return GameMetric{GameAI: strategy.GameAI{Strategy: asm, Mixin: asm}}
}

func (g *GameMetric) StartGame(rule domain.GameRule) {
	g.GameAI.Strategy.StartGame(lib.Size{Width: rule.BoardWidth, Height: rule.BoardHeight})
}

func (g *GameMetric) GetGameState() lib.GameState {
	return g.GameAI.Mixin.GetGameState()
}

func (g *GameMetric) LoadGameState(state lib.GameState) {
	g.GameAI.Mixin.LoadGameState(state)
}

func (g *GameMetric) CreateShips(quantities []domain.ShipQuantity) []ship.Ship {
	ships := make([]ship.Ship, 0)
	for i := 0; i < len(quantities); i++ {
		shipQuantity := quantities[i]
		for j := 0; j < shipQuantity.Quantity; j++ {
			switch constant.ShipType(shipQuantity.Type) {
			case constant.CARRIER:
				ships = append(ships, ship.Ship{Action: &ship.CarrierShip{}})
				break
			case constant.BATTLE_SHIP:
				ships = append(ships, ship.Ship{Action: &ship.BattleShip{}})
				break
			case constant.OIL_RIG:
				ships = append(ships, ship.Ship{Action: &ship.OilRigShip{}})
				break
			case constant.CRUISER:
				ships = append(ships, ship.Ship{Action: &ship.CruiserShip{}})
				break
			case constant.DESTROYER:
				ships = append(ships, ship.Ship{Action: &ship.DestroyerShip{}})
				break
			default:
				break
			}
		}
	}
	return ships
}

func (g *GameMetric) ArrangeShips(boardSize lib.Size, ships []ship.Ship) []ship.Ship {
	return g.GameAI.Strategy.ArrangeShips(ships)
}

func (g *GameMetric) GetShot() lib.Point {
	return g.GameAI.Strategy.GetShot()
}

func (g *GameMetric) ShotResult(result domain.ShotResult) {
	if constant.HIT == result.Status {
		g.GameAI.Strategy.ShotHit(result.Position, &result.RecognizedWholeShip != nil && len(result.RecognizedWholeShip.Positions) > 0)
	} else {
		g.GameAI.Strategy.ShotMiss(result.Position)
	}
}

func (g *GameMetric) Validate(boardSize lib.Size, ships []ship.Ship) bool {
	for i := 0; i < len(ships); i++ {
		if !ships[i].IsValid(boardSize) {
			return false
		}
	}

	for i := 0; i < len(ships); i++ {
		for j := i + 1; j < len(ships); j++ {
			if ships[i].ConflictWith(ships[j]) {
				return false
			}
		}
	}

	return true
}
