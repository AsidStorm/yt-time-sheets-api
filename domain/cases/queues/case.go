package queues

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Queues []models.Queue
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [queues] due [%s]", err)
	}

	queues, err := c.Services().YandexTracker(c.Session()).Queues()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve users due [%s]", err)
	}

	return &Response{
		Queues: queues,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
