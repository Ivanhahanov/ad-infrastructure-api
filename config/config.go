package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf *Config

type Config struct {
	TerraformProjectPath string `yaml:"terraform_project_path"`
	SshKeys              string `yaml:"ssh_keys"`
	AdminMachine         bool   `yaml:"admin_machine"`
	Teams                Teams  `yaml:"teams"`
	Users                Users  `yaml:"users"`
}

type Teams struct {
	Number    int       `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Users struct {
	Number    int       `yaml:"number"`
	Resources Resources `yaml:"resources"`
}

type Resources struct {
	Memory int `yaml:"memory"`
	VCPU   int `yaml:"vcpu"`
}

func ReadConf(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	Conf = &Config{}
	err = yaml.Unmarshal(buf, Conf)
	if Conf.SshKeys == "" {
		Conf.SshKeys = "ssh_keys/"
	}
	if err != nil {
		return fmt.Errorf("in file %q: %v", filename, err)
	}

	return nil
}
