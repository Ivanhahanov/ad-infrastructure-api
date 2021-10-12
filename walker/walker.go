package walker

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/database"
	"github.com/Ivanhahanov/ad-infrastructure-api/walker/providers"
	"log"
	"reflect"
	"strings"
)

func formatLabels(labels map[string]string) string {
	var d []string
	for key,value := range labels{
		d = append(d, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	return strings.Join(d, ",")
}
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
					metricLabels := map[string]string{
						"proto": "http",
						"team": team.Name,
						"service": service.Name,
						"route": http.Route,
					}
					metricNameStr := fmt.Sprintf("walker{%s}", formatLabels(metricLabels))
					flag := providers.GenerateFlag(20)
					//response, httpErr := http.Run(team.Address, flag)
					response, _, httpErr := http.Run("localhost", flag)
					if httpErr != nil {
						log.Println(team.Address, team.Name, service.Name, httpErr)
						promResult[metricNameStr] = 0
						break
					}
					if response.StatusCode == 200 {
						database.PutFlag(&database.FlagStruct{
							Flag:    flag,
							Service: service.Name,
							Team:    team.Name,
						})
						promResult[metricNameStr] = 1
					}
				}
			}
			if !reflect.ValueOf(service.Script).IsZero() {
				for _, script := range service.Script{
					flag := providers.GenerateFlag(20)
					// response, _ := script.RunScript(team.Address, flag)
					response, _ := script.RunScript("localhost", flag)
					log.Println(response)
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
					metricLabels := map[string]string{
						"proto": "http",
						"team": team.Name,
						"service": service.Name,
						"route": http.Route,
					}
					metricNameStr := fmt.Sprintf("checker{%s}", formatLabels(metricLabels))
					//response, httpErr := http.Run(team.Address, "")
					_, body, httpErr := http.Run("localhost", "")
					if httpErr != nil {
						log.Println(team.Address, team.Name, service.Name, httpErr)
						promResult[metricNameStr] = 0
						break
					}
					log.Println(string(body))
					flag, redisErr := database.GetInfo(string(body))
					if redisErr != nil {
						log.Println(team.Address, team.Name, service.Name, redisErr)
					}
					log.Println(flag)
					promResult[metricNameStr] = 0
					if flag[0] == team.Name && flag[1] == service.Name {
						promResult[metricNameStr] = 1
					}
				}
			}
			if !reflect.ValueOf(service.Script).IsZero() {
				for _, script := range service.Script{
					id := "123"
					// response, _ := script.RunScript(team.Address, flag)
					response, _ := script.RunScript("localhost", id)
					log.Println(response)
				}
			}
		}
	}
	return promResult, nil
}
