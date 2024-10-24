package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

type project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func makeProject(p models.Project) project {
	return project{
		Id:   p.Id,
		Name: p.Name,
	}
}

func MarshalProjects(in []models.Project) ([]byte, error) {
	out := struct {
		Data []project `json:"data"`
	}{make([]project, len(in))}

	for i, p := range in {
		out.Data[i] = makeProject(p)
	}

	return json.Marshal(out)
}
