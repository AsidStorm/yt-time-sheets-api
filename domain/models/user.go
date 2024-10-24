package models

type User struct {
	Id              int64
	TrackerId       int64
	Email           string
	Display         string
	HasLicense      bool
	IsAdministrator bool
}
