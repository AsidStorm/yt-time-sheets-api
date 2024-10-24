package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type board struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *service) Boards() ([]models.Board, error) {
	var out []models.Board

	request, err := s.newTrackerRequest(http.MethodGet, "/v2/boards", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create boards request due [%s]", err)
	}

	r, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do boards request due [%s]", err)
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusUnauthorized {
		return nil, nil
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("board response status code [%d] invalid", r.StatusCode)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read board body due [%s]", err)
	}

	var data []board

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal boards body due [%s]", err)
	}

	for _, b := range data {
		out = append(out, models.Board{
			Id:   b.Id,
			Name: b.Name,
		})
	}

	return out, nil
}
