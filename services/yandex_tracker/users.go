package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type user struct {
	Uid        int64  `json:"uid"`
	TrackerUid int64  `json:"trackerUid"`
	Display    string `json:"display"`
	Email      string `json:"email"`
	HasLicense bool   `json:"hasLicense"`
	Groups     []struct {
		Id string `json:"id"`
	} `json:"groups"`
}

func (s *service) Users() ([]models.User, error) {
	var out []models.User

	if err := s.paginatorTracker(http.MethodGet, "/v2/users", nil, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []user

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal users due [%s]", err)
		}

		for _, u := range data {
			out = append(out, models.User{
				Id:         u.Uid,
				TrackerId:  u.TrackerUid,
				Email:      u.Email,
				Display:    u.Display,
				HasLicense: u.HasLicense,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to retrieve users due [%s]", err)
	}

	return out, nil
}
