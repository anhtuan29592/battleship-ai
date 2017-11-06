package service

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/config"
	"strings"
	"log"
	"github.com/anhtuan29592/battleship-ai/lib"
)

type NotifyService struct {
	CacheService *CacheService `inject:""`
}

func (n *NotifyService) HandleNotification(request domain.NotifyRQ) (domain.NotifyRS, error) {
	if config.AI_NAME == request.ShotResult.PlayerId {
		gameMetric := NewGameMetric()

		cacheKey := strings.Join([]string{"gameState", request.SessionId}, "_")
		gameState := &lib.GameState{}
		err := n.CacheService.Get(cacheKey, gameState)
		if err != nil {
			log.Printf("couldn't get game state info of session %s, error = %s", request.SessionId, err)
			log.Print("retry...")
			n.CacheService.Get(cacheKey, gameState)
		}

		gameMetric.LoadGameState(*gameState)
		gameMetric.ShotResult(request.ShotResult)
		n.CacheService.Put(cacheKey, gameMetric.GetGameState())
	}
	return domain.NotifyRS{Success: true}, nil
}

