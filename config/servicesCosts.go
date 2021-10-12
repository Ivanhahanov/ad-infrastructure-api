package config

import (
	"fmt"
	"github.com/Ivanhahanov/ad-infrastructure-api/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServicesCost struct {
	Services []*models.Service `yaml:"services"`
}

func (s *ServicesCost) Load() error {
	buf, err := ioutil.ReadFile("services_costs.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &s)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
