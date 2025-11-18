package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"yandex.tracker.api/domain/models"
)

type issueStatus struct {
	Id          int64  `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *service) IssueStatuses() ([]models.IssueStatus, error) {
	var out []models.IssueStatus

	request, err := s.newTrackerRequest(http.MethodGet, "/v2/statuses", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create statuses request due [%s]", err)
	}

	r, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do statuses request due [%s]", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		fmt.Printf("statuses response status code [%d] invalid", r.StatusCode)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("unable to read statuses body due [%s]", err)
	}

	var data []issueStatus

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal statuses body due [%s] (raw: %s)", err, string(body))
	}

	for _, s := range data {
		out = append(out, models.IssueStatus{
			Id:          s.Id,
			Key:         s.Key,
			Name:        s.Name,
			Description: s.Description,
		})
	}

	return out, nil
}
