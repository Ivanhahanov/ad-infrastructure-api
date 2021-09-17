package main

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.GET("/", routers.GetUserInfo)
			user.POST("/", routers.UpdateUsersKey)
			user.PUT("/", routers.CreateUser)
			user.DELETE("/", routers.DeleteUser)
		}
		team := v1.Group("/command")
		{
			team.GET("/", routers.GetTeamInfo)
			// team.POST("/")
			team.PUT("/", routers.CreateTeam)
			team.DELETE("/", routers.DeleteTeam)

		}

	}
	router.Run(":8080")
}
