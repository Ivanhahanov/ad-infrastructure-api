package models

type Team struct {
	Name    string   `json:"name"`
	SshPubKey string `json:"ssh_pub_key"`
	Players []string `json:"players"`
}
