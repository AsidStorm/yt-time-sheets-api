package domain

import "yandex.tracker.api/domain/services"

type Services interface {
	YandexTracker(Session) services.YandexTracker
	ConfigCache(Session) services.ConfigCache
}
