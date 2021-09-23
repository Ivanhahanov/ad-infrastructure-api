package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SubmitFlagHandler(c *gin.Context) {
	var message string
	var result bool
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"result":  result,
	})
}
