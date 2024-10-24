package issue_types

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	IssueTypes []models.DictionaryIssueType
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [issue_types] due [%s]", err)
	}

	types, err := c.Services().YandexTracker(c.Session()).IssueTypes()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve issue types due [%s]", err)
	}

	return &Response{
		IssueTypes: types,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
