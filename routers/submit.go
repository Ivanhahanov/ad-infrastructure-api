package routers

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SubmitFlagHandler(c *gin.Context) {
	result := false
	team, _ := c.Get("id")
	answer, redisErr := database.GetInfo("")
	if redisErr != nil{
		log.Println(redisErr)
	}
	log.Println(answer)
	log.Println(team.(*models.JWTTeam).TeamName)
	if answer[0] == team.(*models.JWTTeam).TeamName{
		result = true
	}
	// TODO add count to scoreboard
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
