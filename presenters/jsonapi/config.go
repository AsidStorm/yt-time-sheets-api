package jsonapi

import (
	"encoding/json"
	"fmt"
	"yandex.tracker.api/domain/models"
)

type ConfigRequest struct {
	OAuthClientId    string `json:"OAuthClientId"`
	OrganizationId   string `json:"organizationId"`
	FederationId     string `json:"federationId"`
	TrackerAuthUrl   string `json:"trackerAuthUrl"`
	AllowManualInput bool   `json:"allowManualInput"`
}

func UnmarshalConfigRequest(in []byte) (*ConfigRequest, error) {
	out := &ConfigRequest{}

	if err := json.Unmarshal(in, &out); err != nil {
		return nil, fmt.Errorf("unable to unmarshal ConfigRequest due [%s]", err)
	}

	return out, nil
}

type ConfigResponse struct {
	OAuthClientId    string `json:"OAuthClientId"`
	OrganizationId   string `json:"organizationId"`
	FederationId     string `json:"federationId"`
	TrackerAuthUrl   string `json:"trackerAuthUrl"`
	AllowManualInput bool   `json:"allowManualInput"`
	HaveConfig       bool   `json:"haveConfig"`
}

func MarshalConfigResponse(in *models.Config, haveConfig bool) ([]byte, error) {
	value := models.Config{}

	if in != nil {
		value = *in
	}

	return json.Marshal(ConfigResponse{
		OAuthClientId:    value.OAuthClientId,
		OrganizationId:   value.OrganizationId,
		FederationId:     value.FederationId,
		TrackerAuthUrl:   value.TrackerAuthUrl,
		AllowManualInput: value.AllowManualInput,
		HaveConfig:       haveConfig,
	})
}
