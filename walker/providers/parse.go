package providers

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigProviders struct {
	Teams []Team `yaml:"teams"`
	HTTP  []HTTP `yaml:"http"`
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
	for i, http := range p.HTTP {

		if http.Route == "" {
			p.HTTP[i].Route = "/"
		}
		if http.Port == 0 {
			p.HTTP[i].Port = 80
		}

		if http.Schema == "" {
			p.HTTP[i].Schema = "http"
		}
		if http.Method == "" {
			p.HTTP[i].Method = "get"
		}
	}

	return nil
}
