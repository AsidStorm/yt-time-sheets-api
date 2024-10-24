package yandex_tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type iAmToken struct {
	IAmToken  string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (s *service) GetIAmToken() (string, time.Time, error) {
	request, err := s.newIAmRequest(http.MethodPost, "/iam/v1/tokens", bytes.NewReader([]byte(fmt.Sprintf(`{"yandexPassportOauthToken":"%s"}`, s.authToken))))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to create i am token request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to do i am request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("unexpected status code [%d]", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to read i am response body due [%s]", err)
	}

	out := iAmToken{}

	if err := json.Unmarshal(body, &out); err != nil {
		return "", time.Time{}, fmt.Errorf("unable to unmarshal i am body due [%s]", err)
	}

	return out.IAmToken, out.ExpiresAt, nil
}
