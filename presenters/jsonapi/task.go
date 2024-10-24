package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

const taskType = "task"

type task struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes taskAttributes `json:"attributes"`
}

type taskAttributes struct {
	Summary string `json:"summary"`
}

func makeTask(t models.Task) task {
	return task{
		Id:   t.Key,
		Type: taskType,
		Attributes: taskAttributes{
			Summary: t.Summary,
		},
	}
}

func MarshalTasks(in []models.Task) ([]byte, error) {
	response := struct {
		Data []task `json:"data"`
	}{make([]task, len(in))}

	for i, t := range in {
		response.Data[i] = makeTask(t)
	}

	return json.Marshal(response)
}
