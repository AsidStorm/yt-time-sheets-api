package sprint_tasks

import (
	"errors"
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	BoardId  int64
	SprintId int64
}

type Response struct {
	Tasks []models.Task
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialzie case [sprint_tasks] due [%s]", err)
	}

	tasks, err := c.Services().YandexTracker(c.Session()).SprintTasks(r.BoardId, r.SprintId)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve board [%d] sprint [%d] tasks due [%s]", r.BoardId, r.SprintId, err)
	}

	return &Response{
		Tasks: tasks,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if r.SprintId <= 0 {
		return errors.New("request.SprintId must be greater than 0")
	}

	if r.BoardId <= 0 {
		return errors.New("request.BoardId must be greater than 0")
	}

	return domain.ValidateContext(c)
}
