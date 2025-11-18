package issues_collector

import (
	"errors"
	"fmt"
	"time"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	Filter         *RequestFilter         `json:"filter"`
	FilterStatuses *RequestFilterStatuses `json:"filterStatuses"`
	Query          *string                `json:"query"`
}

type RequestFilter struct {
	Queues     []string `json:"queues"`
	Projects   []string `json:"projects"`
	IssueTypes []int64  `json:"issueTypes"`
}

type RequestFilterStatuses struct {
	Month    int      `json:"month"`
	Year     int      `json:"year"`
	Statuses []string `json:"statuses"`
}

type Response struct {
	Issues    map[string]models.Issue `json:"issues"`
	IssueKeys map[string]bool         `json:"issueKeys"`

	DateFrom *time.Time `json:"dateFrom"`
	DateTo   *time.Time `json:"dateTo"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [issues_collector] due [%s]", err)
	}

	var issues []models.Issue
	var err error

	if r.Query != nil {
		issues, err = c.Services().YandexTracker(c.Session()).IssuesWhereQuery(*r.Query)
	} else if r.Filter != nil {
		if r.FilterStatuses != nil {
			issues, err = c.Services().YandexTracker(c.Session()).IssuesInStatus(r.FilterStatuses.Statuses, r.Filter.Queues, r.Filter.Projects, r.FilterStatuses.Month, r.FilterStatuses.Year)
		} else {
			issues, err = c.Services().YandexTracker(c.Session()).IssuesWhereFilter(r.Filter.Queues, r.Filter.Projects, r.Filter.IssueTypes)
		}
	}

	if err != nil {
		if errors.Is(err, models.EmptyFilterErr) {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to retrieve issues due [%s]", err)
	}

	response := &Response{
		Issues:    make(map[string]models.Issue, len(issues)),
		IssueKeys: make(map[string]bool, len(issues)),
	}

	if r.Filter != nil && r.FilterStatuses != nil {
		startDate := time.Date(r.FilterStatuses.Year, time.Month(r.FilterStatuses.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

		response.DateFrom = &startDate
		response.DateTo = &endDate
	}

	for _, issue := range issues {
		response.IssueKeys[issue.Key] = true
		response.Issues[issue.Key] = issue

		if r.Filter != nil && r.FilterStatuses != nil {
			if issue.CreatedAt.Before(*response.DateFrom) {
				response.DateFrom = &issue.CreatedAt
			}
		}
	}

	return response, nil
}

func validate(c domain.Context, r Request) error {
	if r.Filter != nil && r.Query != nil {
		return fmt.Errorf("request.Filter and request.Query cannot be used together")
	}

	if r.Filter == nil && r.Query == nil {
		return fmt.Errorf("request.Filter or request.Query must be used")
	}

	if r.FilterStatuses != nil && r.Filter == nil {
		return fmt.Errorf("request.FilterStatuses can only be used with request.Filter")
	}

	return domain.ValidateContext(c)
}
