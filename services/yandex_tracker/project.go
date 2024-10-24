package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s *service) Projects() ([]models.Project, error) {
	var out []models.Project

	if err := s.paginatorTracker(http.MethodGet, "/v2/projects", nil, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read projects body due [%s]", err)
		}
		defer r.Body.Close()

		var data []project

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal projects due [%s]", err)
		}

		for _, p := range data {
			out = append(out, models.Project{
				Id:   p.Id,
				Name: p.Name,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to retrieve projects due [%s]", err)
	}

	return out, nil
}
