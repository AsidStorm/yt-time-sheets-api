package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type queue struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (s *service) Queues() ([]models.Queue, error) {
	var out []models.Queue

	if err := s.paginatorTracker(http.MethodGet, "/v2/queues", nil, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []queue

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal queues due [%s]", err)
		}

		for _, q := range data {
			out = append(out, models.Queue{
				Key:  q.Key,
				Name: q.Name,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to retrieve queues due [%s]", err)
	}

	return out, nil
}
