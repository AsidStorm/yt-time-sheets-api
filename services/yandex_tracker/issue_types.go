package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type dictionaryIssueType struct {
	Id   int64  `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (s *service) IssueTypes() ([]models.DictionaryIssueType, error) {
	var out []models.DictionaryIssueType

	if err := s.paginatorTracker(http.MethodGet, "/v2/issuetypes", nil, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []dictionaryIssueType

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal issue types due [%s]", err)
		}

		for _, t := range data {
			out = append(out, models.DictionaryIssueType{
				Id:   t.Id,
				Key:  t.Key,
				Name: t.Name,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to retrieve issue types due [%s]", err)
	}

	return out, nil
}
