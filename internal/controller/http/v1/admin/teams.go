package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *routes) TeamsList(c *gin.Context) {
	teams, err := r.teamRepo.GetTeams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teams": teams,
	})
}

func (r *routes) DeleteTeam(c *gin.Context) {
	teamName := c.Param("name")
	err := r.teamRepo.DeleteTeam(teamName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", teamName),
	})
}