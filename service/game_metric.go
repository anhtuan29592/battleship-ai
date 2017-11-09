package service

import (
	"github.com/anhtuan29592/paladin/domain"
	"github.com/anhtuan29592/paladin/lib"
	"github.com/anhtuan29592/paladin/lib/ship"
	"github.com/anhtuan29592/paladin/lib/constant"
	"github.com/anhtuan29592/paladin/lib/strategy"
	"fmt"
)

type GameMetric struct {
	GameAI strategy.GameAI
}

func NewGameMetric() GameMetric {
	asm := &strategy.SampleStrategy{}
	return GameMetric{GameAI: strategy.GameAI{Strategy: asm, Mixin: asm}}
}

func (g *GameMetric) StartGame(rule domain.GameRule, ships []ship.Ship) {
	g.GameAI.Strategy.StartGame(lib.Size{Width: rule.BoardWidth, Height: rule.BoardHeight}, ships)
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
	arrangedShips := strategy.ArrangeShips(boardSize, ships)
	PrintShip(boardSize, arrangedShips)
	return arrangedShips
}

func (g *GameMetric) GetShot() lib.Point {
	return g.GameAI.Strategy.GetShot()
}

func (g *GameMetric) ShotResult(result domain.ShotResult) {
	if constant.HIT == result.Status {
		g.GameAI.Strategy.ShotHit(result.Position, result.RecognizedWholeShip.Type, result.RecognizedWholeShip.Positions)
	} else {
		g.GameAI.Strategy.ShotMiss(result.Position)
	}
}

func PrintShip(boardSize lib.Size, ships []ship.Ship) {
	for r := 0; r < boardSize.Height; r++ {
		fmt.Print("|")
		for c := 0; c < boardSize.Width; c++ {
			printed := false
			for i := 0; i < len(ships); i++ {
				pos := ships[i].GetPositions()
				for j := 0; j < len(pos); j++ {
					if pos[j].X == c && pos[j].Y == r {
						printed = true
						fmt.Print(ships[i].GetType())
						fmt.Print("|")
					}
				}
			}
			if !printed {
				fmt.Print("--")
				fmt.Print("|")
			}
		}
		fmt.Print("\n")
	}

}