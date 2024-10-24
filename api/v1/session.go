package v1

type session struct {
	authToken      string
	iAmToken       string
	organizationId string
}

func (s session) AuthToken() string {
	return s.authToken
}

func (s session) IAmToken() string {
	return s.iAmToken
}

func (s session) OrganizationID() string {
	return s.organizationId
}

func (s session) IsAuthorized() bool {
	return (s.IAmToken() != "" || s.AuthToken() != "") && s.OrganizationID() != ""
}
