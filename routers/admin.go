package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"users": []string{"test"},
	})
}

func TeamsList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"teams": []string{"Test"},
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

func GenerateVariables(c *gin.Context) {
	filename := "variables.tf"
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("file %s generated", filename),
	})
}
