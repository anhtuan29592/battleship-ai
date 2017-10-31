package context

import (
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/anhtuan29592/battleship-ai/domain"
	"net/http"
)

func Invite(c *gin.Context) {
	var request domain.GameInvitationRQ
	if c.Bind(&request) == nil {
		jsonRQ, _ := json.MarshalIndent(&request, "", "\t")
		fmt.Println(string(jsonRQ))
		c.JSON(http.StatusOK, domain.GameInvitationRS{true})
	}

}
