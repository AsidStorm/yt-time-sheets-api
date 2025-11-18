package result_v3

import (
	"fmt"
	"time"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/cases/issues_collector"
	"yandex.tracker.api/domain/cases/worklogs_collector"
	"yandex.tracker.api/domain/cases/worklogs_retriever"
	"yandex.tracker.api/domain/models"
)

// Переработанный result, будем использовать его в дальнейшем

type Request struct {
	UserIdentities []string                                `json:"userIdentities"`
	DateFrom       time.Time                               `json:"dateFrom"`
	DateTo         time.Time                               `json:"dateTo"`
	Filter         *issues_collector.RequestFilter         `json:"filter"`
	FilterStatuses *issues_collector.RequestFilterStatuses `json:"filterStatuses"`
	Query          *string                                 `json:"query"`
}

type Response struct {
	WorkLogs []models.WorkLog `json:"workLogs"`
	DateFrom time.Time        `json:"dateFrom"`
	DateTo   time.Time        `json:"dateTo"`
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [result_v3] due [%s]", err)
	}

	issues, err := issues_collector.Run(c, issues_collector.Request{
		Filter:         r.Filter,
		FilterStatuses: r.FilterStatuses,
		Query:          r.Query,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve issues due [%s]", err)
	}

	var rawWorkLogs []models.RawWorkLog
	var rawIssues map[string]models.Issue

	response := &Response{
		DateFrom: r.DateFrom,
		DateTo:   r.DateTo,
	}

	if issues != nil {
		if issues.DateFrom != nil && issues.DateTo != nil {
			response.DateFrom = *issues.DateFrom
			response.DateTo = *issues.DateTo
		}

		if len(issues.IssueKeys) == 0 {
			return response, nil
		}

		workLogs, err := worklogs_collector.Run(c, worklogs_collector.Request{
			DateFrom:       response.DateFrom,
			DateTo:         response.DateTo,
			IssueKeys:      issues.IssueKeys,
			UserIdentities: stringSliceToMap(r.UserIdentities),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to collect work logs due [%s]", err)
		}

		rawWorkLogs = workLogs.WorkLogs
		rawIssues = issues.Issues
	} else {
		workLogs, err := worklogs_retriever.Run(c, worklogs_retriever.Request{
			DateFrom:       r.DateFrom,
			DateTo:         r.DateTo,
			UserIdentities: stringSliceToMap(r.UserIdentities),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve work logs due [%s]", err)
		}

		rawWorkLogs = workLogs.WorkLogs
		rawIssues = workLogs.Issues
	}

	response.WorkLogs = make([]models.WorkLog, 0, len(rawWorkLogs))

	// Собираем вместе реальный WorkLog (с полным набором данных для обратной совместимости)
	for _, log := range rawWorkLogs {
		issue, ok := rawIssues[log.IssueKey]
		if !ok { // Не должно быть никогда
			continue
		}

		response.WorkLogs = append(response.WorkLogs, models.CombineWorkLog(log, issue))
	}

	return response, nil
}

func validate(c domain.Context, r Request) error {
	return domain.ValidateContext(c)
}

func stringSliceToMap(in []string) map[string]bool {
	out := make(map[string]bool, len(in))

	for _, v := range in {
		out[v] = true
	}

	return out
}
