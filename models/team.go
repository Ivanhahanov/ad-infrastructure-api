package models

type Team struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}
