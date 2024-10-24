package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

const groupType = "group"

type group struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes groupAttributes `json:"attributes"`
}

type groupAttributes struct {
	Label   string   `json:"label"`
	Members []string `json:"members"`
}

func makeGroup(g models.Group) group {
	return group{
		Id:   g.Id,
		Type: groupType,
		Attributes: groupAttributes{
			Label:   g.Label,
			Members: g.Members,
		},
	}
}

func MarshalGroups(in []models.Group) ([]byte, error) {
	response := struct {
		Data []group `json:"data"`
	}{make([]group, len(in))}

	for i, g := range in {
		response.Data[i] = makeGroup(g)
	}

	return json.Marshal(response)
}
