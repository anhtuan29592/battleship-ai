package context

import (
	"github.com/anhtuan29592/battleship-ai/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/anhtuan29592/battleship-ai/service"
)

type NotifyContext struct {
	NotifyService *service.NotifyService `inject:""`
}

func (n *NotifyContext) Result(context *gin.Context) {
	var request domain.NotifyRQ
	if context.Bind(&request) == nil {
		response, _ := n.NotifyService.HandleNotification(request)
		context.JSON(http.StatusOK, response)
	}
}
