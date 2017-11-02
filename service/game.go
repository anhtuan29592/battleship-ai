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
		log.Printf("set start game cache error %s", err)
	}

	invitationCacheKey := strings.Join([]string{"invitationRQ", request.SessionId}, "_")
	invitation := &domain.GameInvitationRQ{}
	err = g.CacheService.Get(invitationCacheKey, invitation)
	if err != nil || &invitation == nil {
		log.Fatalf("couldn't get invitation info of session %s, error = %s", request.SessionId, err)
		return domain.GameStartRS{}, err
	}

	// Init game
	gameRule := invitation.GameRule
	gameMetric := NewGameMetric()
	gameMetric.StartGame(gameRule)
	ships := gameMetric.CreateShips(gameRule.Ships)
	ships = gameMetric.ArrangeShips(lib.Size{Width: gameRule.BoardWidth, Height: gameRule.BoardHeight}, ships)

	shipPositions := make([]domain.ShipPosition, 0)
	for i := 0; i < len(ships); i++ {
		shipPositions = append(shipPositions, domain.ShipPosition{Type: ships[i].GetType().String(), Positions: ships[i].GetPositions()})
	}

	gameState := gameMetric.GetGameState()
	gameStateCacheKey := strings.Join([]string{"gameState", request.SessionId}, "_")
	err = g.CacheService.Put(gameStateCacheKey, gameState, 0)
	if err != nil {
		log.Printf("set game state cache error %s", err)
	}

	return domain.GameStartRS{Ships: shipPositions}, nil
}

func (g *GameService) HandleTurn(context *gin.Context, request domain.TurnRQ) (domain.TurnRS, error) {
	gameMetric := NewGameMetric()

	cacheKey := strings.Join([]string{"gameState", request.SessionId}, "_")
	gameState := &lib.GameState{}
	err := g.CacheService.Get(cacheKey, gameState)
	if err != nil {
		log.Printf("couldn't get game state info of session %s, error = %s", request.SessionId, err)
		log.Println("start new game")

		invitationCacheKey := strings.Join([]string{"invitationRQ", request.SessionId}, "_")
		invitation := &domain.GameInvitationRQ{}
		err = g.CacheService.Get(invitationCacheKey, invitation)

		gameMetric.StartGame(invitation.GameRule)
	} else {
		gameMetric.LoadGameState(*gameState)
	}

	shot := gameMetric.GetShot()
	g.CacheService.Put(cacheKey, gameMetric.GetGameState(), 0)
	return domain.TurnRS{FirePosition: shot}, nil
}
