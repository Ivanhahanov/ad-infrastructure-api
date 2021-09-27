package providers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigProviders struct {
	Teams   []Team    `yaml:"teams"`
	Service []Service `yaml:"services"`
}

type Service struct {
	Name string `yaml:"name"`
	HTTP []HTTP `yaml:"http"`
}

type Team struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	IP     string `yaml:"ip"`
}

type HTTP struct {
	Route    string                 `yaml:"route"`
	Schema   string                 `yaml:"schema"`
	Method   string                 `yaml:"method"`
	Port     int                    `yaml:"port"`
	Params   map[string]string      `yaml:"params"`
	Header   map[string]string      `yaml:"header"`
	JsonBody map[string]interface{} `yaml:"json_body"`
}

func (p *ConfigProviders) Parse(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &p)
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}
	for i, team := range p.Service {
		for j, http := range team.HTTP {

			if http.Route == "" {
				p.Service[i].HTTP[j].Route = "/"
			}
			if http.Port == 0 {
				p.Service[i].HTTP[j].Port = 80
			}

			if http.Schema == "" {
				p.Service[i].HTTP[j].Schema = "http"
			}
			if http.Method == "" {
				p.Service[i].HTTP[j].Method = "get"
			}
		}
	}
	return nil
}
