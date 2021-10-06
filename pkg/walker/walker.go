package walker

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/internal/entity"

	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/manager/repositories"
	"github.com/Ivanhahanov/ad-infrastructure-api/pkg/walker/models"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"strings"
)

type Walker struct {
	teamRepo *repositories.TeamRepository
	manager  *manager.CtfManager

	teams   []models.Team    `yaml:"teams"`
	services []models.Service `yaml:"services"`
}

type configProviders struct {
	Teams   []models.Team    `yaml:"teams"`
	Services []models.Service `yaml:"services"`
}

func parseYaml(filename string) (*configProviders, error) {
	var p configProviders

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = yaml.Unmarshal(buf, &p); err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	for i, team := range p.Services {
		for j, http := range team.HTTP {
			serviceHTTP := p.Services[i].HTTP[j]

			if http.Route == "" {
				serviceHTTP.Route = "/"
			}

			if http.Port == 0 {
				serviceHTTP.Port = 80
			}

			if http.Schema == "" {
				serviceHTTP.Schema = "http"
			}

			if http.Method == "" {
				serviceHTTP.Method = "get"
			}
		}
	}
	return &p, nil
}

func New(manager *manager.CtfManager) (*Walker, error) {
	providers, err := parseYaml("walker.yml")

	if err != nil {return nil, err}

	return &Walker{
		teamRepo: manager.TeamRepo,
		manager: manager,

		teams: providers.Teams,
		services: providers.Services,
	}, nil
}

func formatLabels(labels map[string]string) string {
	var d []string
	for key, value := range labels{
		d = append(d, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	return strings.Join(d, ",")
}

func createMetricLabels(team *entity.Team, service models.Service, http models.HTTP) map[string]string {
	return map[string]string{
		"proto": "http",
		"team": team.Name,
		"service": service.Name,
		"route": http.Route,
	}
}

func (w *Walker) PutFlags() (map[string]int, error) {
	promResult := make(map[string]int)

	teams, dbErr := w.teamRepo.GetTeams()
	if dbErr != nil {return nil, dbErr}

	for _, team := range teams {
		for _, service := range w.services {
			if reflect.ValueOf(service.HTTP).IsZero() {continue}

			for _, http := range service.HTTP {
				metricLabels := createMetricLabels(team, service, http)
				metricNameStr := fmt.Sprintf("walker{%s}", formatLabels(metricLabels))

				flag := w.generateFlag(20)

				response, _, httpErr := http.Run("localhost", flag)
				if httpErr != nil {
					log.Println(team.Address, team.Name, service.Name, httpErr)
					promResult[metricNameStr] = 0
					break
				}
				if response.StatusCode == 200 {
					w.manager.PutFlag(flag, service.Name, team.Name)
					promResult[metricNameStr] = 1
				}
			}
		}
	}
	return promResult, nil
}

func (w *Walker) CheckFlags() (map[string]int, error) {
	promResult := make(map[string]int)

	teams, dbErr := w.teamRepo.GetTeams()
	if dbErr != nil {return nil, dbErr}

	for _, team := range teams {
		for _, service := range w.services {
			if reflect.ValueOf(service.HTTP).IsZero() {continue}

			for _, http := range service.HTTP {
				metricLabels := createMetricLabels(team, service, http)
				metricNameStr := fmt.Sprintf("checker{%s}", formatLabels(metricLabels))

				_, body, httpErr := http.Run("localhost", "")
				if httpErr != nil {
					log.Println(team.Address, team.Name, service.Name, httpErr)
					promResult[metricNameStr] = 0
					break
				}
				log.Println(string(body))
				flag, redisErr := w.manager.GetFlags(string(body))
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
	}
	return promResult, nil
}

func (w *Walker) generateFlag(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = models.Letters[rand.Intn(len(models.Letters))]
	}
	return string(b)
}