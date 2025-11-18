package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"yandex.tracker.api/domain/models"
)

func (s *service) IssuesWhereFilter(queues, projects []string, issueTypes []int64) ([]models.Issue, error) {
	var requestParts []string

	if len(queues) > 0 {
		var queuesParts []string

		for _, q := range queues {
			queuesParts = append(queuesParts, fmt.Sprintf(`"%s"`, q))
		}

		requestParts = append(requestParts, fmt.Sprintf("(Queue: %s)", strings.Join(queuesParts, ",")))
	}

	if len(projects) > 0 {
		requestParts = append(requestParts, fmt.Sprintf("(Project: %s)", strings.Join(projects, ",")))
	}

	if len(issueTypes) > 0 {
		requestParts = append(requestParts, fmt.Sprintf("(Type: %s)", strings.Join(int64SliceToString(issueTypes), ",")))
	}

	if len(requestParts) == 0 {
		return nil, models.EmptyFilterErr
	}

	return s.IssuesWhereQuery(strings.Join(requestParts, " AND "))
}

func (s *service) IssuesWhereQuery(query string) ([]models.Issue, error) {
	var out []models.Issue

	query = strings.TrimSpace(query)

	if query == "" {
		return nil, models.EmptyFilterErr
	}

	in, err := json.Marshal(filterTaskRequest{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal [filterTaskRequest] due [%s]", err)
	}

	if err := s.paginatorTracker(http.MethodPost, "/v3/issues/_search", in, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read issues body due [%s]", err)
		}
		defer r.Body.Close()

		var response []issue

		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("unable to unmarshal issues due [%s]", err)
		}

		for _, i := range response {
			if i.Spent == "" || i.Spent == "PT0S" {
				continue
			}

			is, err := makeDomainIssue(i)
			if err != nil {
				return fmt.Errorf("unable to make domain issue due [%s]", err)
			}

			out = append(out, *is)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to paginate through issues due [%s]", err)
	}

	return out, nil
}

func int64SliceToString(in []int64) []string {
	out := make([]string, 0, len(in))

	for _, v := range in {
		out = append(out, fmt.Sprintf("%d", v))
	}

	return out
}
