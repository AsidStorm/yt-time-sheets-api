package jsonapi

import (
	"encoding/json"
	"fmt"
)

type DeleteWorkLogRequest struct {
	WorkLogId int64  `json:"workLogId"`
	IssueKey  string `json:"issueKey"`
}

func UnmarshalDeleteWorkLogRequest(in []byte) (*DeleteWorkLogRequest, error) {
	request := &DeleteWorkLogRequest{}

	if err := json.Unmarshal(in, &request); err != nil {
		return nil, fmt.Errorf("unable to unmarshal delete work log request due [%s]", err)
	}

	return request, nil
}
