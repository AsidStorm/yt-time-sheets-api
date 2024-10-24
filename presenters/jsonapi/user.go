package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

type user struct {
	Id              int64  `json:"id"`
	TrackerId       int64  `json:"trackerId"`
	Email           string `json:"email"`
	Display         string `json:"display"`
	HasLicense      bool   `json:"hasLicense"`
	IsAdministrator bool   `json:"isAdministrator"`
}

func makeUser(u models.User) user {
	return user{
		Id:              u.Id,
		TrackerId:       u.TrackerId,
		Email:           u.Email,
		Display:         u.Display,
		HasLicense:      u.HasLicense,
		IsAdministrator: u.IsAdministrator,
	}
}

func MarshalUsersAndGroups(users []models.User, groups []models.Group) ([]byte, error) {
	response := struct {
		Users  []user  `json:"users"`
		Groups []group `json:"groups"`
	}{
		Users:  make([]user, len(users)),
		Groups: make([]group, len(groups)),
	}

	for i, u := range users {
		response.Users[i] = makeUser(u)
	}

	for i, g := range groups {
		response.Groups[i] = makeGroup(g)
	}

	return json.Marshal(response)
}

func MarshalUser(u models.User) ([]byte, error) {
	return json.Marshal(struct {
		Data user `json:"data"`
	}{makeUser(u)})
}
