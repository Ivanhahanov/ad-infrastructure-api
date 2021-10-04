package routers

import (
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ShowTeamStatus(c *gin.Context) {
	teamName := c.Param("name")
	teamStatus, sources := database.GetTeamStatus(teamName)
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

func ShowScoreboard(c *gin.Context) {
	var status string
	teams, dbErr := database.GetTeams()
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
		teamHistory := database.GetTeamHistory(team.Name)
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
