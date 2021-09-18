package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Teams Teams `yaml:"teams"`
	Users Users `yaml:"users"`
}

type Teams struct {
	Number    int       `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Users struct {
	Number int `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Resources struct {
	Memory int `yaml:"memory"`
	VCPU   int `yaml:"vcpu"`
}

func ReadConf(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
