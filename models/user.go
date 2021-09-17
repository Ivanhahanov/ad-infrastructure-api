package models

type User struct {
	Nickname string `json:"nickname"`
	FullName string `json:"full_name"`
	PubKey   string `json:"pub_key"`
}
