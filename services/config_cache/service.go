package config_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"yandex.tracker.api/domain/models"
	"yandex.tracker.api/domain/services"
)

type service struct {
	mutex   *sync.Mutex
	storage *models.Config
}

func Make() services.ConfigCache {
	return &service{
		mutex:   &sync.Mutex{},
		storage: nil,
	}
}

func (s *service) Get() (*models.Config, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.storage != nil {
		value := *s.storage

		return &value, nil
	}

	file, err := os.OpenFile("config.json", os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("unable to read file [config.json] due [%s]", err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read all file [config.json] due [%s]", err)
	}
	defer file.Close()

	data := &data{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal [config.json] due [%s]", err)
	}

	s.storage = &models.Config{
		OAuthClientId:    data.OAuthClientId,
		OrganizationId:   data.OrganizationId,
		FederationId:     data.FederationId,
		TrackerAuthUrl:   data.TrackerAuthUrl,
		AllowManualInput: data.AllowManualInput,
	}

	value := *s.storage

	return &value, nil
}

func (s *service) Reset() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.storage = nil

	return nil
}

func (s *service) Store(in *models.Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data := data{
		OAuthClientId:    in.OAuthClientId,
		OrganizationId:   in.OrganizationId,
		FederationId:     in.FederationId,
		TrackerAuthUrl:   in.TrackerAuthUrl,
		AllowManualInput: in.AllowManualInput,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal [config.json] due [%s]", err)
	}

	if err := os.WriteFile("config.json", body, 0644); err != nil {
		return fmt.Errorf("unable to store [config.json] in fs due [%s]", err)
	}

	value := *in

	s.storage = &value

	return nil
}

type data struct {
	OAuthClientId    string `json:"OAuthClientId"`
	OrganizationId   string `json:"organizationId"`
	FederationId     string `json:"federationId"`
	TrackerAuthUrl   string `json:"trackerAuthUrl"`
	AllowManualInput bool   `json:"allowManualInput"`
}
