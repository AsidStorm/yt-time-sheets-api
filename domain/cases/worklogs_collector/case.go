package worklogs_collector

import (
	"errors"
	"fmt"
	"time"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	DateFrom time.Time
	DateTo   time.Time

	UserIdentities map[string]bool
	IssueKeys      map[string]bool
}

type Response struct {
	WorkLogs []models.RawWorkLog
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [worklogs_collector] due [%s]", err)
	}

	logs, err := c.Services().YandexTracker(c.Session()).FilterWorkLogs(r.DateFrom, r.DateTo, r.UserIdentities, func(log models.RawWorkLog) bool {
		if len(r.IssueKeys) > 0 && !r.IssueKeys[log.IssueKey] {
			return false
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve work logs due [%s] (date_from: %s; date_to: %s)", err, r.DateFrom, r.DateTo)
	}

	return &Response{
		WorkLogs: logs,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if len(r.IssueKeys) == 0 {
		return errors.New("request.IssueKeys is empty")
	}

	return domain.ValidateContext(c)
}
