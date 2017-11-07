package strategy

import "github.com/anhtuan29592/paladin/lib"

type Mixin interface {
	GetGameState() lib.GameState
	LoadGameState(state lib.GameState)
}
