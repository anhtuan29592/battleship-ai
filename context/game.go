package context

import (
	"github.com/gin-gonic/gin"
	"github.com/anhtuan29592/battleship-ai/domain"
	"fmt"
	"net/http"
	"encoding/json"
)

func Start(c *gin.Context) {
	var request domain.GameStartRQ
	if c.Bind(&request) == nil {
		jsonRQ, _ := json.MarshalIndent(&request, "", "\t")
		fmt.Println(string(jsonRQ))
		c.JSON(http.StatusOK, domain.GameStartRS{})
	}
}

func Turn(c *gin.Context) {
	var request domain.TurnRQ
	if c.Bind(&request) == nil {
		jsonRQ, _ := json.MarshalIndent(&request, "", "\t")
		fmt.Println(string(jsonRQ))
		c.JSON(http.StatusOK, domain.TurnRS{})
	}
}
