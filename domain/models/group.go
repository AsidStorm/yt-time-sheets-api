package models

type Group struct {
	Id      string   `json:"id"`
	Label   string   `json:"label"`
	Members []string `json:"members"`
}
