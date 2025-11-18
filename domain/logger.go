package domain

type Logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
