package ping

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	PingResult models.PingResult
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [ping] due [%s]", err)
	}

	haveAccess, err := c.Services().YandexTracker(c.Session()).Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to ping due [%s]", err)
	}

	if !haveAccess {
		return &Response{
			PingResult: models.PingResultNeedAuthorize,
		}, nil
	}

	return &Response{
		PingResult: models.PingResultHaveAccess,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
