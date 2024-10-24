package issue_statuses

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	IssueStatuses []models.IssueStatus `json:"issueStatuses"`
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [issue_statuses] due [%s]", err)
	}

	statuses, err := c.Services().YandexTracker(c.Session()).IssueStatuses()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve issue statuses due [%s]", err)
	}

	return &Response{
		IssueStatuses: statuses,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
