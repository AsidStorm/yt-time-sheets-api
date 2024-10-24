package result_v2

import (
	"fmt"
	"strings"
	"time"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/cases/result"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	UserIdentities []string           `json:"userIdentities"`
	Queues         []string           `json:"queues"`
	Projects       []string           `json:"projects"`
	Month          int                `json:"month"`
	Year           int                `json:"year"`
	Statuses       []string           `json:"statuses"`
	IssueTypes     []int64            `json:"issueTypes"`
	ResultGroup    models.ResultGroup `json:"resultGroup"`
}

func Run(c domain.Context, r Request) (*result.Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [result_v2] due [%s]", err)
	}

	issues, err := c.Services().YandexTracker(c.Session()).IssuesInStatus(r.Statuses, r.Queues, r.Month, r.Year)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve issues in statuses [%s] due [%s] (monthd: %d; year: %d)", err, strings.Join(r.Statuses, ", "), r.Month, r.Year)
	}

	if len(issues) == 0 {
		return &result.Response{}, nil
	}

	startDate := time.Date(r.Year, time.Month(r.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	issuesFilter := make(map[string]bool, len(issues))

	for _, i := range issues {
		issuesFilter[i.Key] = true

		if i.CreatedAt.Before(startDate) {
			startDate = i.CreatedAt
		}
	}

	return result.Run(c, result.Request{
		UserIdentities: r.UserIdentities,
		Queues:         r.Queues,
		Projects:       r.Projects,
		DateFrom:       startDate,
		DateTo:         endDate,
		IssueTypes:     r.IssueTypes,
		ResultGroup:    r.ResultGroup,
		IssuesFilter:   issuesFilter,
	})
}

func validate(c domain.Context, r Request) error {
	return domain.ValidateContext(c)
}
