package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

const queueType = "queue"

type queue struct {
	Id         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes queueAttributes `json:"attributes"`
}

type queueAttributes struct {
	Name string `json:"name"`
}

func makeQueue(q models.Queue) queue {
	return queue{
		Id:   q.Key,
		Type: queueType,
		Attributes: queueAttributes{
			Name: q.Name,
		},
	}
}

func MarshalQueues(queues []models.Queue) ([]byte, error) {
	response := struct {
		Data []queue `json:"data"`
	}{make([]queue, len(queues))}

	for i, u := range queues {
		response.Data[i] = makeQueue(u)
	}

	return json.Marshal(response)
}
