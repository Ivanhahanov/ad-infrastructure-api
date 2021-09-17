package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTeamInfo(c *gin.Context) {
	team := models.Team{
		Name:    "Test",
		Players: []string{
			"test1",
			"test2",
			"test3",
		},
	}
	c.JSON(http.StatusOK, team)
}
func CreateTeam(c *gin.Context) {
	var team models.Team
	jsonErr := c.BindJSON(&team)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("The team %s created", team.Name),
	})
}

func DeleteTeam(c *gin.Context) {
	teamName := "Test"
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", teamName),
	})
}
