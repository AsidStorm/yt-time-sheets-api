package get_i_am_token

import (
	"fmt"
	"time"
	"yandex.tracker.api/domain"
)

type Response struct {
	IAmToken  string
	ExpiresAt time.Time
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [get_i_am_token] due [%s]", err)
	}

	token, expires, err := c.Services().YandexTracker(c.Session()).GetIAmToken()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve i am token due [%s]", err)
	}

	return &Response{
		IAmToken:  token,
		ExpiresAt: expires,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
