package walker

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/walker/providers"
	"log"
	"reflect"
	"strings"
)

func PutFlags() map[string]int {
	var c providers.ConfigProviders
	err := c.Parse("providers.yml")
	if err != nil {
		log.Fatalln(err)
	}
	promResult := make(map[string]int)
	for _, team := range c.Teams {
		var address string
		if team.IP != "" {
			address = team.IP
		} else if team.Domain != "" {
			address = team.IP
		}
		if !reflect.ValueOf(c.HTTP).IsZero() {
			for _, http := range c.HTTP {
				result, httpErr := http.Run(address)
				if httpErr != nil {
					log.Println(httpErr)
				}
				metricName := fmt.Sprintf("http_%s%s", team.Name, strings.Replace(http.Route, "/", "_", -1))
				promResult[metricName] = result
			}
		}
	}
	return promResult
}
