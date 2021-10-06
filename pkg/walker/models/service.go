package models

type Service struct {
	Name string `yaml:"name"`
	HTTP []HTTP `yaml:"http"`
}
