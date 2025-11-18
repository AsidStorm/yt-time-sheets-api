package yandex_tracker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"yandex.tracker.api/domain/helpers"
	"yandex.tracker.api/domain/models"
)

type workLog struct {
	Id    int64 `json:"id"`
	Issue struct {
		Key     string `json:"key"`
		Display string `json:"display"`
	} `json:"issue"`
	Comment   string `json:"comment"`
	CreatedBy struct {
		Id      string `json:"id"`
		Display string `json:"display"`
	} `json:"createdBy"`
	Start     string `json:"start"` // Время начала работ. Вся аналитика строится по нему. А createdAt - это уже для всяких хитростей.
	CreatedAt string `json:"createdAt"`
	Duration  string `json:"duration"`
}

func makeDomainWorkLog(wl workLog) (*models.RawWorkLog, error) {
	duration, err := helpers.ParseDuration(wl.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration [%s] due [%s]", wl.Duration, err)
	}

	start, err := time.Parse("2006-01-02T15:04:05.999-0700", wl.Start)
	if err != nil {
		return nil, fmt.Errorf("unable to parse start time [%s] due [%s]", wl.Start, err)
	}

	log := &models.RawWorkLog{
		Id:              wl.Id,
		Duration:        duration,
		CreatedById:     wl.CreatedBy.Id,
		CreateByDisplay: wl.CreatedBy.Display,
		IssueKey:        wl.Issue.Key,
		IssueDisplay:    wl.Issue.Display,
		Comment:         wl.Comment,
		CreatedAt:       start, // 2018-06-06T08:42:06.258+0000
		Queue:           extractQueue(wl.Issue.Key),
	}

	return log, nil
}

func extractQueue(issueKey string) string {
	split := strings.Split(issueKey, "-")

	if len(split) > 0 {
		return split[0]
	}

	return "NONE"
}

func (s *service) FilterWorkLogs(dateFrom, dateTo time.Time, userIdentities map[string]bool, filter func(log models.RawWorkLog) bool) ([]models.RawWorkLog, error) {
	var out []models.RawWorkLog
	var err error

	dateFromParts := strings.Split(dateFrom.String(), " ")
	dateToParts := strings.Split(dateTo.String(), " ")

	dateFrom, err = time.Parse("2006-01-02", dateFromParts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse date from [%s] due [%s]", dateFromParts[0], err)
	}

	dateTo, err = time.Parse("2006-01-02", dateToParts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse date to [%s] due [%s]", dateToParts[0], err)
	}

	dateFrom = dateFrom.Truncate(time.Hour * 24).UTC()
	dateTo = dateTo.Truncate(time.Hour * 24).UTC().Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	stored := make(map[int64]bool)

	handler := func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []workLog

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal work logs due [%s]", err)
		}

		for _, l := range data {
			if _, ok := stored[l.Id]; ok {
				continue
			}

			if len(userIdentities) > 0 && !userIdentities[l.CreatedBy.Id] { // Фильтр по UserIdentities вшит, т.к. это позволяет значительно оптимизировать при одном юзере (наиболее частый кейс) запрос
				continue
			}

			wl, err := makeDomainWorkLog(l)
			if err != nil {
				return fmt.Errorf("unable to convert tracker work log to domain due [%s]", err)
			}

			if !filter(*wl) {
				continue
			}

			stored[l.Id] = true

			out = append(out, *wl)
		}

		return nil
	}

	limit := 200

	interval := time.Hour * 24 * 3
	if dateTo.Sub(dateFrom).Abs() > time.Hour*24*30 {
		interval = time.Hour * 24 * 15
	}

	for dateFrom.Before(dateTo) {
		to := dateFrom.Add(interval).Add(time.Second * -1)
		if to.After(dateTo) {
			to = dateTo
		}

		input := []byte(fmt.Sprintf(`{"start":{"from":"%s%s","to":"%s%s"}}`, dateFrom.Format("2006-01-02T15:04:05.999"), dateFromParts[2], to.Format("2006-01-02T15:04:05.999"), dateFromParts[2]))

		if len(userIdentities) == 1 {
			for identity, _ := range userIdentities {
				input = []byte(fmt.Sprintf(`{"createdBy":"%s","start":{"from":"%s%s","to":"%s%s"}}`, identity, dateFrom.Format("2006-01-02T15:04:05.999"), dateFromParts[2], to.Format("2006-01-02T15:04:05.999"), dateFromParts[2]))
			}
		}

		if err := s.paginatorTracker(http.MethodPost, "/v3/worklog/_search", input, handler); err != nil {
			return nil, fmt.Errorf("unable to paginate work logs due [%s]", err)
		}

		dateFrom = dateFrom.Add(interval)
		limit--

		if limit <= 0 {
			return nil, errors.New("200 days limit")
		}
	}

	return out, nil
}

func (s *service) WorkLogs(userIdentities, queues []string, dateFrom, dateTo time.Time) ([]models.WorkLog, []string, error) {
	var out []models.WorkLog
	var issuesKeys []string
	var err error

	issuesKeysStorage := make(map[string]bool)

	dateFromParts := strings.Split(dateFrom.String(), " ")
	dateToParts := strings.Split(dateTo.String(), " ")

	dateFrom, err = time.Parse("2006-01-02", dateFromParts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse date from [%s] due [%s]", dateFromParts[0], err)
	}

	dateTo, err = time.Parse("2006-01-02", dateToParts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse date to [%s] due [%s]", dateToParts[0], err)
	}

	dateFrom = dateFrom.Truncate(time.Hour * 24).UTC()
	dateTo = dateTo.Truncate(time.Hour * 24).UTC().Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	haveQueuesFilter := len(queues) > 0
	haveUsersFilter := len(userIdentities) > 0

	usersFilter := make(map[string]bool)

	for _, userId := range userIdentities {
		usersFilter[userId] = true
	}

	stored := make(map[int64]bool)

	handler := func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read body due [%s]", err)
		}
		defer r.Body.Close()

		var data []workLog

		if err := json.Unmarshal(body, &data); err != nil {
			return fmt.Errorf("unable to unmarshal work logs due [%s]", err)
		}

		for _, l := range data {
			if haveUsersFilter && !usersFilter[l.CreatedBy.Id] {
				continue
			}

			if haveQueuesFilter {
				ok := false

				for _, q := range queues {
					if strings.HasPrefix(l.Issue.Key, q) {
						ok = true
						break
					}
				}

				if !ok {
					continue
				}
			}

			if _, ok := stored[l.Id]; ok {
				continue
			}

			stored[l.Id] = true

			d, err := helpers.ParseDuration(l.Duration)
			if err != nil {
				return fmt.Errorf("unable to parse duration [%s] due [%s]", l.Duration, err)
			}

			t, err := time.Parse("2006-01-02T15:04:05.999-0700", l.Start)
			if err != nil {
				return fmt.Errorf("unable to parse time [%s] due [%s]", l.Start, err)
			}

			out = append(out, models.WorkLog{
				Id:              l.Id,
				Duration:        d,
				CreatedById:     l.CreatedBy.Id,
				CreateByDisplay: l.CreatedBy.Display,
				IssueKey:        l.Issue.Key,
				IssueDisplay:    l.Issue.Display,
				Comment:         l.Comment,
				CreatedAt:       t, // 2018-06-06T08:42:06.258+0000
			})

			if !issuesKeysStorage[l.Issue.Key] {
				issuesKeysStorage[l.Issue.Key] = true
				issuesKeys = append(issuesKeys, l.Issue.Key)
			}
		}

		return nil
	}

	limit := 150

	interval := time.Hour * 24
	if dateTo.Sub(dateFrom).Abs() > time.Hour*24*30 {
		interval = time.Hour * 24 * 30
	}

	for dateFrom.Before(dateTo) {
		to := dateFrom.Add(interval).Add(time.Second * -1)

		input := []byte(fmt.Sprintf(`{"start":{"from":"%s%s","to":"%s%s"}}`, dateFrom.Format("2006-01-02T15:04:05.999"), dateFromParts[2], to.Format("2006-01-02T15:04:05.999"), dateFromParts[2]))

		if len(userIdentities) == 1 {
			input = []byte(fmt.Sprintf(`{"createdBy":"%s","start":{"from":"%s%s","to":"%s%s"}}`, userIdentities[0], dateFrom.Format("2006-01-02T15:04:05.999"), dateFromParts[2], to.Format("2006-01-02T15:04:05.999"), dateFromParts[2]))
		}

		if err := s.paginatorTracker(http.MethodPost, "/v2/worklog/_search", input, handler); err != nil {
			return nil, nil, fmt.Errorf("unable to paginate work logs due [%s]", err)
		}

		dateFrom = dateFrom.Add(interval)
		limit--

		if limit <= 0 {
			return nil, nil, errors.New("150 days limit")
		}
	}

	return out, issuesKeys, nil
}

type patchWorkLogRequest struct {
	Duration string `json:"duration"`
	Comment  string `json:"comment"`
}

type createMyWorkLogRequest struct {
	Start    string `json:"start"`
	Duration string `json:"duration"`
	Comment  string `json:"comment"`
}

type createAnyWorkLogRequest struct {
	Start     string `json:"start"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	Comment   string `json:"comment"`
	Duration  string `json:"duration"`
}

func (s *service) PatchWorkLog(workLogId int64, issueKey string, d string, comment string) (*models.WorkLog, error) {
	in, err := json.Marshal(patchWorkLogRequest{
		Duration: d,
		Comment:  comment,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal patch work log request due [%s]", err)
	}

	request, err := s.newTrackerRequest(http.MethodPatch, fmt.Sprintf("/v2/issues/%s/worklog/%d", issueKey, workLogId), bytes.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("unable to create patch work log request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do patch work log request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response patch work log code [%d]", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read create work log response body due [%s]", err)
	}

	log := &workLog{}

	if err := json.Unmarshal(body, &log); err != nil {
		return nil, fmt.Errorf("unable to unmarshal patch work log response due [%s] (raw: %s)", err, body)
	}

	dur, err := helpers.ParseDuration(log.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration [%s] due [%s]", log.Duration, err)
	}

	t, err := time.Parse("2006-01-02T15:04:05.999-0700", log.Start)
	if err != nil {
		return nil, fmt.Errorf("unable to parse time [%s] due [%s]", log.Start, err)
	}

	return &models.WorkLog{
		Id:              log.Id,
		Duration:        dur,
		CreatedById:     log.CreatedBy.Id,
		CreateByDisplay: log.CreatedBy.Display,
		Comment:         log.Comment,
		IssueKey:        log.Issue.Key,
		IssueDisplay:    log.Issue.Display,
		CreatedAt:       t,
	}, nil
}

func (s *service) CreateWorkLog(issueKey, userIdentity string, myUser bool, d string, comment string, start time.Time) (*models.WorkLog, error) {
	in, url, err := createWorkLogPayload(issueKey, userIdentity, myUser, d, comment, start)
	if err != nil {
		return nil, fmt.Errorf("unable to build create work log request payload due [%s]", err)
	}

	request, err := s.newTrackerRequest(http.MethodPost, url, bytes.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("unable to create create work log request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do create work log request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response create work log code [%d]", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read create work log response body due [%s]", err)
	}

	log, err := unmarshalCreateWorkLogRequest(body, myUser, start)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal work log response due [%s] (raw: %s)", err, body)
	}

	dur, err := helpers.ParseDuration(log.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration [%s] due [%s]", log.Duration, err)
	}

	t, err := time.Parse("2006-01-02T15:04:05.999-0700", log.Start)
	if err != nil {
		return nil, fmt.Errorf("unable to parse time [%s] due [%s]", log.Start, err)
	}

	return &models.WorkLog{
		Id:              log.Id,
		Duration:        dur,
		CreatedById:     log.CreatedBy.Id,
		CreateByDisplay: log.CreatedBy.Display,
		Comment:         log.Comment,
		IssueKey:        log.Issue.Key,
		IssueDisplay:    log.Issue.Display,
		CreatedAt:       t,
	}, nil
}

func unmarshalCreateWorkLogRequest(body []byte, myUser bool, start time.Time) (*workLog, error) {
	if myUser {
		log := workLog{}

		if err := json.Unmarshal(body, &log); err != nil {
			return nil, fmt.Errorf("unable to unmarshal work log response due [%s] (raw: %s)", err, body)
		}

		return &log, nil
	}

	var log []workLog

	if err := json.Unmarshal(body, &log); err != nil {
		return nil, fmt.Errorf("unable to unmarshal work log response due [%s] (raw: %s)", err, body)
	}

	if len(log) == 0 {
		return nil, fmt.Errorf("zero work logs in response")
	}

	return &log[0], nil
}

func createWorkLogPayload(issueKey, userIdentity string, myUser bool, duration string, comment string, start time.Time) ([]byte, string, error) {
	if myUser {
		in, err := json.Marshal(createMyWorkLogRequest{
			Start:    start.Format("2006-01-02T15:04:05.999Z07:00"),
			Duration: duration,
			Comment:  comment,
		})
		if err != nil {
			return nil, "", fmt.Errorf("unable to marshal create my work log request due [%s]", err)
		}

		return in, fmt.Sprintf("/v2/issues/%s/worklog", issueKey), nil
	}

	in, err := json.Marshal(createAnyWorkLogRequest{
		Start:     start.Format("2006-01-02T15:04:05.999Z07:00"),
		CreatedBy: userIdentity,
		CreatedAt: start.Format("2006-01-02T15:04:05.999Z07:00"),
		Comment:   comment,
		Duration:  duration,
	})
	if err != nil {
		return nil, "", fmt.Errorf("unable to marshal create any work log request due [%s]", err)
	}

	return in, fmt.Sprintf("/v2/issues/%s/worklogs/_import", issueKey), nil
}

func (s *service) DeleteWorkLog(workLogId int64, issueKey string) error {
	request, err := s.newTrackerRequest(http.MethodDelete, fmt.Sprintf("/v2/issues/%s/worklog/%d", issueKey, workLogId), nil)
	if err != nil {
		return fmt.Errorf("unable to create delete work log request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return fmt.Errorf("unable to do delete work log request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response delete work log code [%d]", response.StatusCode)
	}

	return nil
}
