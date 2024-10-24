package projects

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Projects []models.Project
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [projects] due [%s]", err)
	}

	projects, err := c.Services().YandexTracker(c.Session()).Projects()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve projects due [%s]", err)
	}

	return &Response{
		Projects: projects,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
