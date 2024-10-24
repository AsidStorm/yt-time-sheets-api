package yandex_tracker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"yandex.tracker.api/domain/services"
)

type service struct {
	hostTracker           string
	host360               string
	hostCloudOrganization string
	hostIAm               string
	authToken             string
	iAmToken              string
	orgId                 string
}

var ErrForbidden = errors.New("forbidden")

func MakeService(authToken, iAmToken, orgId string) services.YandexTracker {
	return &service{
		hostTracker:           "https://api.tracker.yandex.net",
		host360:               "https://api360.yandex.net",
		hostIAm:               "https://iam.api.cloud.yandex.net",
		hostCloudOrganization: "https://organization-manager.api.cloud.yandex.net",
		authToken:             authToken,
		orgId:                 orgId,
		iAmToken:              iAmToken,
	}
}

func (s *service) newIAmRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", s.hostIAm, url), body)
}

func (s *service) newTrackerRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", s.hostTracker, url), body)
}

func (s *service) new360Request(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", s.host360, url), body)
}

func (s *service) newCloudOrganizationRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, fmt.Sprintf("%s%s", s.hostCloudOrganization, url), body)
}

func (s *service) isCloudOrgId() bool {
	if _, err := strconv.Atoi(s.orgId); err != nil {
		return true
	}

	return false
}

func (s *service) doRequest(r *http.Request, forceIAmToken bool) (*http.Response, error) {
	if s.iAmToken != "" || forceIAmToken {
		if s.iAmToken == "" {
			iAmToken, _, err := s.GetIAmToken()
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve iAmToken due [%s]", err)
			}

			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", iAmToken))
		} else {
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.iAmToken))
		}
	} else {
		r.Header.Set("Authorization", fmt.Sprintf("OAuth %s", s.authToken))
	}

	if s.isCloudOrgId() {
		r.Header.Set("X-Cloud-Org-ID", s.orgId)
	} else {
		r.Header.Set("X-Org-ID", s.orgId)
	}

	return http.DefaultClient.Do(r)
}

func (s *service) paginatorTracker(method, url string, in []byte, handler func(response *http.Response) error) error {
	page := 1

	for {
		realUrl := fmt.Sprintf("%s?perPage=%d&page=%d", url, 500, page)
		if strings.Contains(url, "?") {
			realUrl = fmt.Sprintf("%s&perPage=%d&page=%d", url, 500, page)
		}

		body := bytes.NewReader(in)

		request, err := s.newTrackerRequest(method, realUrl, body)
		if err != nil {
			return fmt.Errorf("unable to create request due [%s]", err)
		}

		response, err := s.doRequest(request, false)
		if err != nil {
			return fmt.Errorf("unable to do request due [%s]", err)
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()

			return fmt.Errorf("unexpected status code [%d]", response.StatusCode)
		}

		totalPages, err := strconv.Atoi(response.Header.Get("X-Total-Pages"))
		if err != nil {
			response.Body.Close()
			return fmt.Errorf("unable to extract X-Total-Pages due [%s]", err)
		}

		if response.Header.Get("X-Total-Count") == "10000" {
			response.Body.Close()
			return fmt.Errorf("API limit for 10000, try to specify your query")
		}

		if err := handler(response); err != nil {
			response.Body.Close()
			return fmt.Errorf("unable to process handler due [%s]", err)
		}

		response.Body.Close()

		if page >= totalPages {
			break
		}

		page++
	}

	return nil
}

func (s *service) paginator360(method, url string, body io.Reader, handler func(response *http.Response) (int, error)) error {
	page := 1

	for {
		realUrl := fmt.Sprintf("%s?perPage=%d&page=%d", url, 100, page)
		if strings.Contains(url, "?") {
			realUrl = fmt.Sprintf("%s&perPage=%d&page=%d", url, 100, page)
		}

		request, err := s.new360Request(method, realUrl, body)
		if err != nil {
			return fmt.Errorf("unable to create request due [%s]", err)
		}

		response, err := s.doRequest(request, false)
		if err != nil {
			return fmt.Errorf("unable to do request due [%s]", err)
		}

		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusUnauthorized {
			response.Body.Close()
			return ErrForbidden
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return fmt.Errorf("unexpected status code [%d]", response.StatusCode)
		}

		totalPages, err := handler(response)
		if err != nil {
			response.Body.Close()
			return fmt.Errorf("unable to process handler due [%s]", err)
		}

		response.Body.Close()

		if page >= totalPages {
			break
		}

		page++
	}

	return nil
}
