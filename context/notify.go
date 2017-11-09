package context

import (
	"github.com/anhtuan29592/paladin/domain"
	"github.com/anhtuan29592/paladin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotifyContext struct {
	NotifyService *service.NotifyService `inject:""`
}

func (n *NotifyContext) Result(context *gin.Context) {
	var request domain.NotifyRQ
	if context.Bind(&request) == nil {
		//log.Printf("notify request %s" ,request)
		response, _ := n.NotifyService.HandleNotification(request)
		context.JSON(http.StatusOK, response)
	}
}
