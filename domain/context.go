package domain

type Context interface {
	Session() Session
	Services() Services
	WithSession(Session) Context
	Logger() Logger
}

func ValidateContext(c Context) error {
	return nil
}
