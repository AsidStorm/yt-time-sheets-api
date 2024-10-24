package jsonapi

import (
	"encoding/json"
	"time"
)

type IAmTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func MarshalIAmTokenResponse(token string, expires time.Time) ([]byte, error) {
	out := struct {
		Data IAmTokenResponse `json:"data"`
	}{IAmTokenResponse{
		Token:   token,
		Expires: expires,
	}}

	return json.Marshal(out)
}
