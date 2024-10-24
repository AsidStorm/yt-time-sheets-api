package yandex_tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"yandex.tracker.api/domain/models"
)

type task struct {
	Key     string `json:"key"`
	Summary string `json:"summary"`
}

type filterTaskRequest struct {
	Query string `json:"query"`
}

func (s *service) SprintTasks(boardId, sprintId int64) ([]models.Task, error) {
	var out []models.Task

	in, err := json.Marshal(filterTaskRequest{
		Query: fmt.Sprintf(`"Boards": %d "Sprint": %d`, boardId, sprintId),
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

		var response []task

		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("unable to unmarshal board sprint tasks due [%s]", err)
		}

		for _, t := range response {
			out = append(out, models.Task{
				Key:     t.Key,
				Summary: t.Summary,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("unable to paginate throught issues due [%s]", err)
	}

	return out, nil
}

func (s *service) FilterTasks(query string) ([]models.Task, error) {
	in, err := json.Marshal(filterTaskRequest{
		Query: makeFilterTask(query),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal filter tasks request due [%s]", err)
	}

	request, err := s.newTrackerRequest(http.MethodPost, "/v2/issues/_search?perPage=30&page=1", bytes.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("unable to make filter tasks request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return nil, fmt.Errorf("unable to do filter tasks request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code [%d] not ok", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read filter tasks body due [%s]", err)
	}

	var out []task

	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("unable to unmarshal filter tasks due [%s]", err)
	}

	tasks := make([]models.Task, len(out))

	for i, t := range out {
		tasks[i] = models.Task{
			Key:     t.Key,
			Summary: t.Summary,
		}
	}

	return tasks, nil
}

func makeFilterTask(in string) string {
	if strings.Contains(in, "-") && !strings.Contains(in, " ") && !strings.HasPrefix(in, "-") && !strings.HasSuffix(in, "-") {
		return fmt.Sprintf(`"Key": "%s"`, in) // Поиск по ключу
	}

	return fmt.Sprintf(`"Summary": ~"%s" OR "Description": "%s"`, in, in)
}
