package yandex_tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *service) CreateIssueComment(issueKey, issueComment string) error {
	in, err := json.Marshal(createIssueCommentRequest{
		Text: issueComment,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal create issue comment request due [%s]", err)
	}

	request, err := s.newTrackerRequest(http.MethodPost, fmt.Sprintf("/v2/issues/%s/comments", issueKey), bytes.NewReader(in))
	if err != nil {
		return fmt.Errorf("unable to create issue comment request due [%s]", err)
	}

	response, err := s.doRequest(request, false)
	if err != nil {
		return fmt.Errorf("unable to do create issue comment request due [%s]", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected create issue comment response code [%d] (expected: %d)", response.StatusCode, http.StatusCreated)
	}

	return nil
}

type createIssueCommentRequest struct {
	Text string `json:"text"`
}
