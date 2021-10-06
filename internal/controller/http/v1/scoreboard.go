package v1

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/logger"
	_manager "github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type STeam struct {
	TeamName string     `json:"name"`
	SLA      float64    `json:"sla"`
	Services []SService `json:"services"`
}

type SService struct {
	Name  string  `json:"name"`
	Value string  `json:"value"`
	SLA   float64 `json:"sla"`
}

type scoreboardRoutes struct {
	l         logger.Interface
	manager  *_manager.CtfManager
}
func NewScoreboardRoutes(handler *gin.RouterGroup, l logger.Interface, manager *_manager.CtfManager, jwtMiddleware *jwt.GinJWTMiddleware) {
	r := &scoreboardRoutes{l, manager}

	scoreboardGroup := handler.Group("/scoreboard")

	scoreboardGroup.GET("/", jwtMiddleware.MiddlewareFunc(), r.ShowScoreboard)
	scoreboardGroup.GET("/:name", jwtMiddleware.MiddlewareFunc(), r.ShowTeamStatus)
}

func (r *scoreboardRoutes) ShowTeamStatus(c *gin.Context) {
	teamName := c.Param("name")
	teamStatus, sources := r.manager.GetTeamStatus(teamName)
	result := map[string][]string{}
	var status string
	var totalServiceOKStatus = 0.0
	var serviceNum = 0.0
	for serviceName, value := range teamStatus {
		if value == sources {
			status = "OK"
			totalServiceOKStatus += 1
		} else if value == 0 {
			status = "DOWN"
		} else if value < sources {
			status = "MUMBLE"
		}
		result[serviceName] = append(result[serviceName], status)
		serviceNum += 1
	}
	log.Println(teamStatus, sources)
	c.JSON(http.StatusOK, gin.H{teamName: result})
}

func (r *scoreboardRoutes) ShowScoreboard(c *gin.Context) {
	var status string
	teams, dbErr := r.manager.TeamRepo.GetTeams()
	if dbErr != nil {
		log.Println(dbErr)
	}
	var scoreboard []STeam
	for _, team := range teams {
		var serviceNum = 0.0
		var totalStatus = 0.0
		sTeam := STeam{
			TeamName: team.Name,
		}
		teamHistory := r.manager.GetTeamHistory(team.Name)
		for serviceName, values := range teamHistory.RoundsHistory {
			sService := SService{}
			var totalServiceOKStatus = 0.0
			for i := 1; i < len(values); i++ {
				if values[i] == teamHistory.Sources {
					status = "OK"
					totalServiceOKStatus += 1
				} else if values[i] == 0 {
					status = "DOWN"
				} else if values[i] < teamHistory.Sources {
					status = "MUMBLE"
				}
				sService.Name = serviceName
				sService.Value = status
			}
			serviceNum += 1
			totalStatus += totalServiceOKStatus / teamHistory.TotalRounds
			sService.SLA = totalServiceOKStatus / teamHistory.TotalRounds * 100
			sTeam.Services = append(sTeam.Services, sService)
		}
		sTeam.SLA = totalStatus / serviceNum * 100
		scoreboard = append(scoreboard, sTeam)
		log.Println(teamHistory)
	}
	c.JSON(http.StatusOK, gin.H{"scoreboard": scoreboard})
}
