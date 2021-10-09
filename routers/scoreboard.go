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

type ScoreboardTeamJson struct {
	TeamName string                  `json:"name"`
	SLA      float64                 `json:"sla"`
	Score    float64                 `json:"score"`
	Services []ScoreboardServiceJson `json:"services"`
}

type ScoreboardServiceJson struct {
	Name         string  `json:"name"`
	Value        string  `json:"value"`
	SLA          float64 `json:"sla"`
	FP           float64 `json:"fp"`
	Gained       float64 `json:"gained"`
	Lost         float64 `json:"lost"`
	ServiceScore float64 `json:"score"`
}

func ShowScoreboard(c *gin.Context) {
	var status string
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		log.Println(dbErr)
	}
	var scoreboard []ScoreboardTeamJson
	for _, team := range teams {
		var serviceNum = 0.0
		var totalStatus = 0.0
		var totalScore = 0.0
		sTeam := ScoreboardTeamJson{
			TeamName: team.Name,
		}
		teamHistory := database.GetTeamHistory(team.Name)
		for serviceName, values := range teamHistory.RoundsHistory {
			sService := ScoreboardServiceJson{}
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
			flags := database.GetServiceFlagsStats(team.Name, serviceName)
			sService.Gained = flags.Gained
			sService.Lost = flags.Lost
			sService.FP = flags.Gained - flags.Lost
			if sService.FP >= 0 {
				sService.ServiceScore = sService.FP * (totalServiceOKStatus / teamHistory.TotalRounds)
			} else if sService.FP < 0 {
				sService.ServiceScore = sService.FP * (1 - totalServiceOKStatus/teamHistory.TotalRounds)
			}
			totalScore += sService.ServiceScore
			sTeam.Services = append(sTeam.Services, sService)
		}
		sTeam.Score = totalScore / serviceNum
		sTeam.SLA = totalStatus / serviceNum * 100
		scoreboard = append(scoreboard, sTeam)
		log.Println(teamHistory)
	}
	c.JSON(http.StatusOK, gin.H{"scoreboard": scoreboard})
}
