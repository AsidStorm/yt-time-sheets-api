package worklogs_retriever

import (
	"fmt"
	"time"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	DateFrom time.Time
	DateTo   time.Time

	UserIdentities map[string]bool
}

type Response struct {
	WorkLogs []models.RawWorkLog
	Issues   map[string]models.Issue
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [worlogs_retriever] due [%s]", err)
	}

	workLogs, err := c.Services().YandexTracker(c.Session()).FilterWorkLogs(r.DateFrom, r.DateTo, r.UserIdentities, func(log models.RawWorkLog) bool {
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve work logs due [%s] (date_from: %s; date_to: %s)", err, r.DateFrom, r.DateTo)
	}

	// Собираем issueKeys по workLogs

	issueKeys := make(map[string]bool)

	for _, log := range workLogs {
		issueKeys[log.IssueKey] = true
	}

	issues, err := c.Services().YandexTracker(c.Session()).IssuesByKeys(mapToSlice(issueKeys))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve issues due [%s] (len(issue_keys): %d)", err, len(issueKeys))
	}

	return &Response{
		WorkLogs: workLogs,
		Issues:   issues,
	}, nil
}

func validate(c domain.Context, r Request) error {
	return domain.ValidateContext(c)
}

func mapToSlice(in map[string]bool) []string {
	out := make([]string, 0, len(in))

	for k, _ := range in {
		out = append(out, k)
	}

	return out
}
