package context

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type FightingController struct {

}


func GetRoute() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	controller := new(FightingController)



	return e
}