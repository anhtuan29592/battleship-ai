package service

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/lib"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

type GameService struct {
	CacheService *CacheService `inject:""`
	GameMetric   *GameMetric
}

func (g *GameService) HandleInvitation(context *gin.Context, request domain.GameInvitationRQ) (domain.GameInvitationRS, error) {
	cacheKey := strings.Join([]string{"invitationRQ", request.SessionId}, "_")
	err := g.CacheService.Put(cacheKey, request, 0)
	if err != nil {
		log.Fatalf("set invitation cache error %s", err)
	}
	return domain.GameInvitationRS{Success: err == nil}, err
}

func (g *GameService) HandleGameStart(context *gin.Context, request domain.GameStartRQ) (domain.GameStartRS, error) {
	cacheKey := strings.Join([]string{"startRQ", request.SessionId}, "_")
	err := g.CacheService.Put(cacheKey, request, 0)
	if err != nil {
		log.Fatalf("set start game cache error %s", err)
	}

	invitationCacheKey := strings.Join([]string{"invitationRQ", request.SessionId}, "_")
	invitation := &domain.GameInvitationRQ{}
	err = g.CacheService.Get(invitationCacheKey, invitation)
	if err != nil || &invitation == nil {
		log.Fatalf("couldn't get invitation info of %s cache, error = %s", request.SessionId, err)
		return domain.GameStartRS{}, err
	}

	// Init ship
	gameRule := invitation.GameRule
	ships := g.GameMetric.CreateShips(gameRule.Ships)
	ships = g.GameMetric.ArrangeShips(lib.Size{Witdh: gameRule.BoardWidth, Height: gameRule.BoardHeight}, ships)

	shipPositions := make([]*domain.ShipPosition, 0)
	for i := 0; i < len(ships); i++ {
		shipPositions = append(shipPositions, &domain.ShipPosition{Type: ships[i].ShipAction.GetType().String(), Positions: ships[i].ShipAction.GetPositions()})
	}
	return domain.GameStartRS{Ships: shipPositions}, nil
}
