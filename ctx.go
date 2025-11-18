package main

import (
	"go.uber.org/zap"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/services"
	"yandex.tracker.api/services/yandex_tracker"
)

type ctx struct {
	logger  *zap.SugaredLogger
	session domain.Session
	svs     *svs
}

type svs struct {
	configCache services.ConfigCache
}

type session struct {
}

func (s session) AuthToken() string {
	return ""
}

func (s session) IAmToken() string {
	return ""
}

func (s session) OrganizationID() string {
	return ""
}

func (s session) IsAuthorized() bool {
	return false
}

func (s session) TraceId() string {
	return "runtime"
}

func (c *ctx) Logger() domain.Logger {
	return c.logger
}

func (c *ctx) Services() domain.Services {
	return c.svs
}

func (s *svs) YandexTracker(session domain.Session) services.YandexTracker {
	return yandex_tracker.MakeService(session.AuthToken(), session.IAmToken(), session.OrganizationID())
}

func (s *svs) ConfigCache(session domain.Session) services.ConfigCache {
	return s.configCache
}

func (c *ctx) Session() domain.Session {
	return c.session
}

func (c *ctx) WithSession(session domain.Session) domain.Context {
	return &ctx{
		session: session,
		svs:     c.svs,
		logger:  c.logger.With(zap.String("trace_id", session.TraceId())),
	}
}
