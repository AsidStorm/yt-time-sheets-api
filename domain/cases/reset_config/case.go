package reset_config

import (
	"fmt"
	"yandex.tracker.api/domain"
)

func Run(c domain.Context) error {
	if err := validate(c); err != nil {
		return fmt.Errorf("unable to initialize case [reset_config] due [%s]", err)
	}

	if err := c.Services().ConfigCache(c.Session()).Reset(); err != nil {
		return fmt.Errorf("unable to reset config cache due [%s]", err)
	}

	return nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
