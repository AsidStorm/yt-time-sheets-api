package domain

type Session interface {
	AuthToken() string
	IAmToken() string
	OrganizationID() string
	IsAuthorized() bool
	TraceId() string
}
