package models

type User struct {
	Id              int64  `json:"id"`
	TrackerId       int64  `json:"trackerId"`
	Email           string `json:"email"`
	Display         string `json:"display"`
	HasLicense      bool   `json:"hasLicense"`
	IsAdministrator bool   `json:"isAdministrator"`
}
