package routers

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SubmitFlag struct {
	Flag string `json:"flag"`
}

func SubmitFlagHandler(c *gin.Context) {
	var flag *SubmitFlag

	jsonErr := c.BindJSON(&flag)
	if jsonErr != nil {
		c.JSON(http.StatusOK, jsonErr)
	}

	team, _ := c.Get("id")
	teamName := team.(*models.JWTTeam).TeamName
	answer, redisErr := database.GetInfo(flag.Flag)
	if redisErr != nil {
		log.Println(redisErr)
	}
	log.Println(answer[0], teamName)
	if answer[0] != nil {
		if answer[0] != teamName {
			submittedFlags, redisErr := database.GetSubmitFlags(flag.Flag)
			if redisErr != nil {
				log.Println(redisErr)
			}
			if submittedFlags[0] == nil && submittedFlags[1] == nil {
				serviceName := answer[1].(string)
				database.AddSubmitFlag(&database.FlagStruct{
					Flag:    flag.Flag,
					Team:    teamName,
					Service: serviceName,
				})
				database.AddAttackFlag(teamName, serviceName)
				database.AddDefenceFlag(answer[0].(string), serviceName)
				c.JSON(http.StatusOK, gin.H{
					"result":  true,
					"message": "Flag is submit successfully",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"result":  false,
				"message": "you have already submit this flag",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "you can't submit your flags",
		})
		return
	}
	// TODO add count to scoreboard
	c.JSON(http.StatusOK, gin.H{
		"result":  false,
		"message": "invalid flag",
	})
}
