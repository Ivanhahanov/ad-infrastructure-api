package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []string{"test"},
	})
}

func TeamsList(c *gin.Context) {
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

func DeleteUsers(c *gin.Context) {
	user := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("user %s deleted", user),
	})
}

func DeleteTeams(c *gin.Context) {
	team := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("team %s deleted", team),
	})
}

func CountTeamsHandler(c *gin.Context) {
	teamsCount := database.CountTeams()
	c.JSON(http.StatusOK, gin.H{"count": teamsCount})
}
