package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yandex.tracker.api/domain/models"
)

type sprint struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (s *service) Sprints(boardId int64) ([]models.Sprint, error) {
	var out []models.Sprint

	request, err := s.newTrackerRequest(http.MethodGet, fmt.Sprintf("/v2/boards/%d/sprints", boardId), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create sprints request due [%s]", err)
	}

	r, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do sprints request due [%s]", err)
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusBadRequest {
		// Возникает когда у доски нету спринтов

		return nil, nil
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sprints response status code [%d] invalid", r.StatusCode)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read sprints body due [%s]", err)
	}

	var data []sprint

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal sprints body due [%s]", err)
	}

	for _, s := range data {
		startDate, err := time.Parse("2006-01-02", s.StartDate)
		if err != nil {
			return nil, fmt.Errorf("unable to parse start date [%s] due [%s]", s.StartDate, err)
		}

		endDate, err := time.Parse("2006-01-02", s.EndDate)
		if err != nil {
			return nil, fmt.Errorf("unable to parse end date [%s] due [%s]", s.EndDate, err)
		}

		out = append(out, models.Sprint{
			Id:        s.Id,
			Name:      s.Name,
			Status:    s.Status,
			StartDate: startDate,
			EndDate:   endDate,
		})
	}

	return out, nil
}
