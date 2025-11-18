package yandex_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"yandex.tracker.api/domain/models"
)

func getStartAndEndDates(month, year int) (time.Time, time.Time, error) {
	// Validate month
	if month < 1 || month > 12 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %d", month)
	}
	// Validate year
	if year < 1 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %d", year)
	}

	// Start of the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// End of the month: set date to the first day of the next month and subtract one second
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	return startDate, endDate, nil
}

func (s *service) IssuesInStatus(statuses []string, queues, projects []string, month, year int) ([]models.Issue, error) {
	var out []models.Issue

	start, end, err := getStartAndEndDates(month, year)
	if err != nil {
		return nil, fmt.Errorf("unable to retreive start and end dates due [%s] (month: %d; year: %d)", err, month, year)
	}

	// ((Status: changed(to: "done" date: 01.06.2024 .. 30.06.2024))  OR  (Status: changed(to: "wontdo" date: 01.06.2024 .. 30.06.2024))) AND (Status: "done","wontdo")

	var parts []string
	var wrappedStatuses []string
	var wrappedQueues []string

	for _, status := range statuses {
		parts = append(parts, fmt.Sprintf(`(Status: changed(to: "%s" date: %s .. %s))`, status, start.Format("02.01.2006"), end.Format("02.01.2006")))
		wrappedStatuses = append(wrappedStatuses, fmt.Sprintf(`"%s"`, status))
	}

	for _, q := range queues {
		wrappedQueues = append(wrappedQueues, fmt.Sprintf(`"%s"`, q))
	}

	queuesFilter := ""
	if len(wrappedQueues) > 0 {
		queuesFilter = " AND (Queue: " + strings.Join(wrappedQueues, ",") + ")"
	}

	projectsFilter := ""
	if len(projects) > 0 {
		projectsFilter = " AND (Project: " + strings.Join(projects, ",") + ")"
	}

	in, err := json.Marshal(filterTaskRequest{
		Query: fmt.Sprintf(`(%s AND (Status: %s))%s%s`, strings.Join(parts, " OR "), strings.Join(wrappedStatuses, ","), queuesFilter, projectsFilter),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal filter tasks request due [%s]", err)
	}

	if err := s.paginatorTracker(http.MethodPost, "/v2/issues/_search", in, func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("unable to read board sprint tasks body due [%s]", err)
		}
		defer r.Body.Close()

		var response []issue

		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("unable to unmarshal board sprint tasks due [%s]", err)
		}

		layout := "2006-01-02T15:04:05.000-0700"

		for _, i := range response {
			if i.Spent == "" || i.Spent == "PT0S" {
				continue
			}

			createdAt, err := time.Parse(layout, i.CreatedAt)
			if err != nil {
				return fmt.Errorf("unable to parse issue [%s] created at [%s] due [%s]", i.Key, i.CreatedAt, err)
			}

			is := models.Issue{
				Key:       i.Key,
				Summary:   i.Summary,
				CreatedAt: createdAt,
			}

			if i.Epic != nil {
				is.Epic = &models.IssueEpic{
					Key:     i.Epic.Key,
					Summary: i.Epic.Display,
				}
			}

			if i.Project != nil {
				is.Project = &models.IssueProject{
					Id:   i.Project.Id,
					Name: i.Project.Display,
				}
			}

			if i.Type != nil {
				is.Type = &models.IssueType{
					Id:      i.Type.Id,
					Key:     i.Type.Key,
					Display: i.Type.Display,
				}
			}

			out = append(out, is)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to paginate throught issues due [%s]", err)
	}

	return out, nil
}
