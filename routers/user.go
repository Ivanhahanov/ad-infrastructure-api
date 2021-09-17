package routers

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserInfo(c *gin.Context) {
	user := models.User{
		Nickname: "",
		FullName: "",
		PubKey:   "",
	}
	c.JSON(http.StatusOK, user)
}
func CreateUser(c *gin.Context) {
	var user models.User
	jsonErr := c.BindJSON(&user)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s created", user.Nickname),
	})
}
func UpdateUsersKey(c *gin.Context) {
	var user models.User
	jsonErr := c.BindJSON(&user)
	if jsonErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": jsonErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s's key updated", "test"),
	})
}
func DeleteUser(c *gin.Context) {
	username := "test"
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s deleted", username),
	})
}
