package users_and_groups

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Users  []models.User
	Groups []models.Group
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [users_and_groups] due [%s]", err)
	}

	users, groups, err := c.Services().YandexTracker(c.Session()).UsersAndGroups()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve users and groups due [%s]", err)
	}

	c.Services().YandexTracker(c.Session()).IssueStatuses()

	return &Response{
		Users:  users,
		Groups: groups,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
