package jsonapi

import (
	"encoding/json"
	"fmt"
)

type PatchWorkLogRequest struct {
	WorkLogId int64  `json:"workLogId"`
	IssueKey  string `json:"issueKey"`
	Duration  string `json:"duration"`
	Comment   string `json:"comment"`
}

func UnmarshalPatchWorkLogRequest(in []byte) (*PatchWorkLogRequest, error) {
	out := &PatchWorkLogRequest{}

	if err := json.Unmarshal(in, &out); err != nil {
		return nil, fmt.Errorf("unable to unmarshal patch work log request due [%s]", err)
	}

	return out, nil
}
