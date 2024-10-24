package config

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Config     *models.Config
	HaveConfig bool
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [config] due [%s]", err)
	}

	config, err := c.Services().ConfigCache(c.Session()).Get()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve config from cache due [%s]", err)
	}

	return &Response{
		Config:     config,
		HaveConfig: config != nil,
	}, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
