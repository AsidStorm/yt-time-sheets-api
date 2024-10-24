package jsonapi

import (
	"encoding/json"
	"fmt"
	"time"
)

type CreateWorkLogRequest struct {
	UserIdentity string    `json:"userIdentity"`
	MyUser       bool      `json:"myUser"`
	IssueKey     string    `json:"issueKey"`
	Duration     string    `json:"duration"`
	Comment      string    `json:"comment"`
	Date         time.Time `json:"date"`
	IssueComment string    `json:"issueComment"`
}

func UnmarshalCreateWorkLogRequest(in []byte) (*CreateWorkLogRequest, error) {
	request := &CreateWorkLogRequest{}

	if err := json.Unmarshal(in, &request); err != nil {
		return nil, fmt.Errorf("unable to unmarshal create work log request due [%s]", err)
	}

	return request, nil
}
