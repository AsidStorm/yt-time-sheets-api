package sprints

import (
	"errors"
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	BoardId int64
}

type Response struct {
	Sprints []models.Sprint
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [sprints] due [%s]", err)
	}

	sprints, err := c.Services().YandexTracker(c.Session()).Sprints(r.BoardId)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve board [%d] sprints due [%s]", r.BoardId, err)
	}

	return &Response{
		Sprints: sprints,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if r.BoardId <= 0 {
		return errors.New("request.BoardId must be greater than 0")
	}

	return domain.ValidateContext(c)
}
