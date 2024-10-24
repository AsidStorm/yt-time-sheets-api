package filter_tasks

import (
	"errors"
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	Query string
}

type Response struct {
	Tasks []models.Task
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [filter_tasks] due [%s]", err)
	}

	tasks, err := c.Services().YandexTracker(c.Session()).FilterTasks(r.Query)
	if err != nil {
		return nil, fmt.Errorf("unable to filter tasks due [%s] (query: %s)", err, r.Query)
	}

	return &Response{
		Tasks: tasks,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if len(r.Query) < 3 {
		return errors.New("request.Query must be greater or equal 3 symbols")
	}

	return domain.ValidateContext(c)
}
