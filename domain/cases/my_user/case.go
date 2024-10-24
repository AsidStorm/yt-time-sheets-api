package my_user

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	User models.User
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initiaize case [my_user] due [%s]", err)
	}

	user, err := c.Services().YandexTracker(c.Session()).MyUser()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve my user due [%s]", err)
	}

	return &Response{
		User: *user,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
