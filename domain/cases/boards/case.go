package boards

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Boards []models.Board
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [boards] due [%s]", err)
	}

	boards, err := c.Services().YandexTracker(c.Session()).Boards()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve boards due [%s]", err)
	}

	return &Response{
		Boards: boards,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
