package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

func (s *service) MyUser() (*models.User, error) {
	request, err := s.newTrackerRequest(http.MethodGet, "/v2/myself", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create my user request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do my user request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code [%d]", response.StatusCode)
	}

	in, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read my user body due [%s]", err)
	}

	user := &user{}

	if err := json.Unmarshal(in, &user); err != nil {
		return nil, fmt.Errorf("unable to unmarshal my user due [%s]", err)
	}

	// Мы запрашиваем так же информацию из 360, для определения, является ли пользователь администратором. Если пользователь не администратор, то мы не можем создавать записи данных от имени другого пользователя.

	isAdministrator := false

	if user.HasLicense {
		if s.isCloudOrgId() {
			for {
				request, err := s.newCloudOrganizationRequest(http.MethodGet, fmt.Sprintf("/organization-manager/v1/organizations/%s:listAccessBindings", s.orgId), nil)
				if err != nil {
					break
				}

				response, err := s.doRequest(request, true)
				if err != nil {
					response.Body.Close()
					break
				}

				fmt.Println(response.StatusCode)

				if response.StatusCode == http.StatusOK {
					isAdministrator = true
				}

				response.Body.Close()

				break
			}
		} else {
			for {
				request, err := s.new360Request(http.MethodGet, fmt.Sprintf("/directory/v1/org/%s/groups", s.orgId), nil)
				if err != nil {
					break
				}

				response, err := s.doRequest(request, false)
				if err != nil {
					response.Body.Close()
					break
				}

				if response.StatusCode == http.StatusOK {
					isAdministrator = true
				}

				response.Body.Close()

				break
			}
		}
	}

	return &models.User{
		Id:              user.Uid,
		TrackerId:       user.TrackerUid,
		Email:           user.Email,
		Display:         user.Display,
		HasLicense:      user.HasLicense,
		IsAdministrator: isAdministrator,
	}, nil
}
