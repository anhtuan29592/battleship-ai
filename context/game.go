package context

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/anhtuan29592/battleship-ai/lib"
	stg "github.com/anhtuan29592/battleship-ai/lib/strategy"
	"github.com/gin-gonic/gin"
	"net/http"
)

var strategy *stg.AgentSmith

func Start(c *gin.Context) {
	var request domain.GameStartRQ
	if c.Bind(&request) == nil {
		strategy = new(stg.AgentSmith)
		strategy.StartGame(lib.Size{4, 4})
		c.JSON(http.StatusOK, domain.GameStartRS{})
	}
}

func Turn(c *gin.Context) {
	var request domain.TurnRQ
	if c.Bind(&request) == nil {
		c.JSON(http.StatusOK, strategy.GetShot())
	}
}
