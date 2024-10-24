package jsonapi

import (
	"encoding/json"
	"yandex.tracker.api/domain/models"
)

func MarshalPingResponse(result models.PingResult) ([]byte, error) {
	out := struct {
		Data models.PingResult `json:"data"`
	}{result}

	return json.Marshal(out)
}
