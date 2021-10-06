package v1

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	_manager "github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SubmitFlag struct {
	Flag string `json:"flag"`
}

type flagRoutes struct {
	l         logger.Interface
	manager  *_manager.CtfManager
}
func NewFlagRoutes(handler *gin.RouterGroup, l logger.Interface, manager *_manager.CtfManager, jwtMiddleware *jwt.GinJWTMiddleware) {
	r := &flagRoutes{l, manager}

	handler.POST("/flag/submit", jwtMiddleware.MiddlewareFunc(), r.SubmitFlagHandler)
}

func (r *flagRoutes) SubmitFlagHandler(c *gin.Context) {
	var flag *SubmitFlag

	jsonErr := c.BindJSON(&flag)
	if jsonErr != nil {
		c.JSON(http.StatusOK, jsonErr)
	}

	team, _ := c.Get("id")
	teamName := team.(*entity.JWTTeam).TeamName
	answer, redisErr := r.manager.GetFlags(flag.Flag)
	if redisErr != nil {
		log.Println(redisErr)
	}
	log.Println(answer[0], teamName)
	if answer[0] == nil {
		// TODO add count to scoreboard
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "invalid flag",
		})
		return
	}

	if answer[0] == teamName {
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "you can't submit your flags",
		})
		return
	}

	submittedFlags, redisErr := r.manager.GetSubmitFlags(flag.Flag)
	if redisErr != nil {
		log.Println(redisErr)
	}

	if submittedFlags[0] != nil || submittedFlags[1] != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  false,
			"message": "you have already submit this flag",
		})
		return
	}

	serviceName := answer[1].(string)
	r.manager.PutFlag(flag.Flag, teamName, serviceName)
	r.manager.FlagRepo.AddAttackFlag(teamName, serviceName)
	r.manager.FlagRepo.AddDefenceFlag(answer[0].(string), serviceName)
	c.JSON(http.StatusOK, gin.H{
		"result":  true,
		"message": "Flag is submit successfully",
	})
}
