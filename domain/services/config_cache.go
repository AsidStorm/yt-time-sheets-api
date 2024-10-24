package services

import "yandex.tracker.api/domain/models"

type ConfigCache interface {
	Get() (*models.Config, error)
	Reset() error
	Store(in *models.Config) error
}
