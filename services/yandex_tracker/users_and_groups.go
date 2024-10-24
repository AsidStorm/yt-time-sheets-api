package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"yandex.tracker.api/domain/models"
)

type group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type int64  `json:"type"`
}

// Enum: "UNKNOWN" "DEPARTMENT" "WIKI" "SERVICE" "CONNECT_GROUP" "CONNECT_DEPARTMENT" "CLOUD" "COM_ALL_USERS"

const yandexCloudGroupTypeId = 6
const yandex360DepartmentTypeId = 5
const yandex360GroupTypeId = 4

func (s *service) UsersAndGroups() ([]models.User, []models.Group, error) {
	var users []models.User

	// 1. Забираем группы
	groups := make(map[string]*models.Group)

	if err := s.paginatorTracker(http.MethodGet, "/v2/groups", nil, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []group

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal groups due [%s]", err)
		}

		for _, g := range data {
			if g.Type != yandexCloudGroupTypeId && g.Type != yandex360GroupTypeId && g.Type != yandex360DepartmentTypeId {
				// Пропускаем странные группы
				continue
			}

			groups[fmt.Sprintf("%d", g.Id)] = &models.Group{
				Id:    fmt.Sprintf("%d", g.Id),
				Label: g.Name,
			}
		}

		return nil
	}); err != nil {
		if err != ErrForbidden {
			return nil, nil, fmt.Errorf("unable to paginate 360 groups due [%s]", err)
		}
	}

	if err := s.paginatorTracker(http.MethodGet, "/v2/users?expand=groups", nil, func(r *http.Response) error {
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
			users = append(users, models.User{
				Id:         u.Uid,
				TrackerId:  u.TrackerUid,
				Email:      u.Email,
				Display:    u.Display,
				HasLicense: u.HasLicense,
			})

			for _, g := range u.Groups {
				if _, ok := groups[g.Id]; !ok {
					continue
				}

				groups[g.Id].Members = append(groups[g.Id].Members, fmt.Sprintf("%d", u.Uid))
			}
		}

		return nil
	}); err != nil {
		if err != ErrForbidden {
			return nil, nil, fmt.Errorf("unable to paginate users due [%s]", err)
		}
	}

	var out []models.Group

	for _, g := range groups {
		if len(g.Members) == 0 {
			continue
		}

		out = append(out, *g)
	}

	return users, out, nil
}
