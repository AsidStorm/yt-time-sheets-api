package yandex_tracker

import (
	"fmt"
	"net/http"
)

func (s *service) Ping() (bool, error) {
	request, err := s.newTrackerRequest(http.MethodGet, "/v2/myself", nil)
	if err != nil {
		return false, fmt.Errorf("unable to create ping request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return false, fmt.Errorf("unable to do ping request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return false, nil
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected response code [%d]", response.StatusCode)
	}

	return true, nil
}
