package walker

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/walker/providers"
	"log"
	"reflect"
	"strings"
)

func PutFlags() (map[string]int, error) {
	var c providers.ConfigProviders
	err := c.Parse("walker.yml")
	if err != nil {
		return nil, err
	}
	promResult := make(map[string]int)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		return nil, dbErr
	}
	for _, team := range teams {
		for _, service := range c.Service {
			if !reflect.ValueOf(service.HTTP).IsZero() {
				for _, http := range service.HTTP {
					metricName := fmt.Sprintf("http_%s_%s%s", team.Name, service.Name, strings.Replace(http.Route, "/", "_", -1))
					flag := providers.GenerateFlag(20)
					//response, httpErr := http.Run(team.Address, flag)
					response, _, httpErr := http.Run("localhost", flag)
					if httpErr != nil {
						log.Println(team.Address, team.Name, service.Name, httpErr)
						promResult[metricName] = 0
						break
					}
					if response.StatusCode == 200 {
						database.PutFlag(database.FlagStruct{
							Flag:    flag,
							Service: service.Name,
							Team:    team.Name,
						})
						promResult[metricName] = 1
					}
				}
			}
		}
	}
	return promResult, nil
}

func CheckFlags() (map[string]int, error) {
	var c providers.ConfigProviders
	err := c.Parse("checker.yml")
	if err != nil {
		return nil, err
	}
	promResult := make(map[string]int)
	teams, dbErr := database.GetTeams()
	if dbErr != nil {
		return nil, dbErr
	}
	for _, team := range teams {
		for _, service := range c.Service {
			if !reflect.ValueOf(service.HTTP).IsZero() {
				for _, http := range service.HTTP {
					metricName := fmt.Sprintf("http_%s_%s%s", team.Name, service.Name, strings.Replace(http.Route, "/", "_", -1))
					//response, httpErr := http.Run(team.Address, "")
					_, body, httpErr := http.Run("localhost", "")
					if httpErr != nil {
						log.Println(team.Address, team.Name, service.Name, httpErr)
						promResult[metricName] = 0
						break
					}
					log.Println(string(body))
					flag, redisErr := database.GetInfo(string(body))
					if redisErr != nil {
						log.Println(team.Address, team.Name, service.Name, redisErr)
					}
					log.Println(flag)
					if flag[0] == team.Name && flag[1] == service.Name {
						promResult[metricName] = 1
					}
				}
			}
		}
	}
	return promResult, nil
}
