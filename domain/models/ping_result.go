package models

type PingResult string

const (
	PingResultHaveAccess    PingResult = "HAVE_ACCESS"
	PingResultNeedAuthorize PingResult = "NEED_AUTHORIZE"
)
