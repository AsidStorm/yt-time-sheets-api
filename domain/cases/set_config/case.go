package set_config

import (
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	OAuthClientId    string
	OrganizationId   string
	FederationId     string
	TrackerAuthUrl   string
	AllowManualInput bool
}

func Run(c domain.Context, r Request) error {
	if err := validate(c, r); err != nil {
		return fmt.Errorf("unable to initialize case [set_config] due [%s]", err)
	}

	config, err := c.Services().ConfigCache(c.Session()).Get()
	if err != nil {
		return fmt.Errorf("unable to get config from cache due [%s]", err)
	}

	if config != nil {
		return fmt.Errorf("case [set_config] only for initial config set, see docs")
	}

	if err := c.Services().ConfigCache(c.Session()).Store(&models.Config{
		OAuthClientId:    r.OAuthClientId,
		OrganizationId:   r.OrganizationId,
		FederationId:     r.FederationId,
		TrackerAuthUrl:   r.TrackerAuthUrl,
		AllowManualInput: r.AllowManualInput || (r.OAuthClientId == "" && r.OrganizationId == "" && r.FederationId == ""),
	}); err != nil {
		return fmt.Errorf("unable to store config due [%s]", err)
	}

	return nil
}

func validate(c domain.Context, r Request) error {
	return domain.ValidateContext(c)
}
