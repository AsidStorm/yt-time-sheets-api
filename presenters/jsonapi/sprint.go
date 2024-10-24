package jsonapi

import (
	"encoding/json"
	"time"
	"yandex.tracker.api/domain/models"
)

const sprintType = "sprint"

type sprint struct {
	Id         int64            `json:"id"`
	Type       string           `json:"type"`
	Attributes sprintAttributes `json:"attributes"`
}

type sprintAttributes struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func makeSprint(s models.Sprint) sprint {
	return sprint{
		Id:   s.Id,
		Type: sprintType,
		Attributes: sprintAttributes{
			Name:      s.Name,
			Status:    s.Status,
			StartDate: s.StartDate,
			EndDate:   s.EndDate,
		},
	}
}

func MarshalSprints(in []models.Sprint) ([]byte, error) {
	out := struct {
		Data []sprint `json:"data"`
	}{make([]sprint, len(in))}

	for i, s := range in {
		out.Data[i] = makeSprint(s)
	}

	return json.Marshal(out)
}
