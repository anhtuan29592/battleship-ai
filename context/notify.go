package context

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotifyContext struct {
}

func (*NotifyContext) Result(context *gin.Context) {
	var request domain.NotifyRQ
	if context.Bind(&request) == nil {
		context.JSON(http.StatusOK, gin.H{"status": true})
	}
}
