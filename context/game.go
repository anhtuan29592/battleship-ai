package context

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GameContext struct {
	GameService *service.GameService `inject:""`
}

func (g *GameContext) Invite(context *gin.Context) {
	var request domain.GameInvitationRQ
	if context.Bind(&request) == nil {
		response, _ := g.GameService.HandleInvitation(context, request)
		context.JSON(http.StatusOK, response)
	}
}

func (g *GameContext) Start(context *gin.Context) {
	var request domain.GameStartRQ
	if context.Bind(&request) == nil {
		response, _ := g.GameService.HandleGameStart(context, request)
		context.JSON(http.StatusOK, response)
	}
}

func (*GameContext) Turn(c *gin.Context) {
	var request domain.TurnRQ
	if c.Bind(&request) == nil {
		c.JSON(http.StatusOK, gin.H{"status": true})
	}
}
