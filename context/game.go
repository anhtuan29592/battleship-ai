package context

import (
	"github.com/anhtuan29592/paladin/domain"
	"github.com/anhtuan29592/paladin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type GameContext struct {
	GameService *service.GameService `inject:""`
}

func (g *GameContext) Invite(context *gin.Context) {
	var request domain.GameInvitationRQ
	if context.Bind(&request) == nil {
		log.Printf("invite request %s" ,request)
		response, _ := g.GameService.HandleInvitation(context, request)
		context.JSON(http.StatusOK, response)
	}
}

func (g *GameContext) Start(context *gin.Context) {
	var request domain.GameStartRQ
	if context.Bind(&request) == nil {
		log.Printf("start request %s" ,request)
		response, _ := g.GameService.HandleGameStart(context, request)
		context.JSON(http.StatusOK, response)
	}
}

func (g *GameContext) Turn(context *gin.Context) {
	var request domain.TurnRQ
	if context.Bind(&request) == nil {
		log.Printf("turn request %s" ,request)
		response, _ := g.GameService.HandleTurn(context, request)
		context.JSON(http.StatusOK, response)
	}
}
