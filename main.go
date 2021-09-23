package main

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/config"
	"github.com/Ivanhahanov/ad-infrastructure-api/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := config.ReadConf("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Conf)
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/api/v1")
	{
		team := v1.Group("/team")
		{
			team.GET("/", routers.GetTeamInfo)
			// team.POST("/")
			team.PUT("/", routers.CreateTeam)
			team.DELETE("/", routers.DeleteTeam)
		}
		admin := v1.Group("/admin")
		{
			admin.GET("/teams", routers.TeamsList)
			admin.DELETE("/team/:name", routers.DeleteTeams)
			admin.POST("/generate/terraform", routers.GenerateTerraformConfig)
			admin.POST("/generate/variables", routers.GenerateVariables)

		}
		v1.GET("/walker", routers.RunWalkerHandler)
		v1.POST("/submit", routers.SubmitFlagHandler)
	}
	router.Run(":8080")
}
