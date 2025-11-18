package yandex_tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"yandex.tracker.api/domain/models"
)

type issueSearchRequest struct {
	Keys []string `json:"keys"`
}

func (s *service) IssuesByKeys(issueKeys []string) (map[string]models.Issue, error) {
	out := make(map[string]models.Issue, len(issueKeys))

	if len(issueKeys) == 0 {
		return out, nil
	}

	chunks := chunkSlice(issueKeys, 150)

	for idx, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}

		if idx%10 == 0 { // to avoid requests limit
			time.Sleep(time.Second)
		}

		in, err := json.Marshal(issueSearchRequest{Keys: chunk})
		if err != nil {
			return nil, fmt.Errorf("unable to marshal [issueSearchRequest] due [%s]", err)
		}

		body := bytes.NewReader(in)

		request, err := s.newTrackerRequest(http.MethodPost, "/v2/issues/_search", body)
		if err != nil {
			return nil, fmt.Errorf("unable to create request due [%s]", err)
		}

		response, err := s.doRequest(request, false)
		if err != nil {
			return nil, fmt.Errorf("unable to do request due [%s]", err)
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()

			return nil, fmt.Errorf("unexpected status code [%d]", response.StatusCode)
		}

		outBody, err := io.ReadAll(response.Body)
		if err != nil {
			response.Body.Close()

			return nil, fmt.Errorf("unable to read issues body due [%s]", err)
		}

		response.Body.Close()

		var data []issue

		if err := json.Unmarshal(outBody, &data); err != nil {
			return nil, fmt.Errorf("unable to unmarshal issues due [%s]", err)
		}

		for _, i := range data {
			is := models.Issue{
				Key:     i.Key,
				Summary: i.Summary,
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

			out[is.Key] = is
		}
	}

	return out, nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
