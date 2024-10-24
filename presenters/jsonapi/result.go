package jsonapi

import (
	"encoding/json"
	"fmt"
	"time"
	"yandex.tracker.api/domain/models"
)

type ResultRequest struct {
	UserIdentities []string           `json:"userIdentities"`
	Queues         []string           `json:"queues"`
	DateFrom       time.Time          `json:"dateFrom"`
	DateTo         time.Time          `json:"dateTo"`
	Projects       []string           `json:"projects"`
	IssueTypes     []int64            `json:"issueTypes"`
	ResultGroup    models.ResultGroup `json:"resultGroup"`
}

func UnmarshalResultRequest(in []byte) (*ResultRequest, error) {
	out := &ResultRequest{}

	if err := json.Unmarshal(in, &out); err != nil {
		return nil, fmt.Errorf("unable to unmarshal result request due [%s]", err)
	}

	return out, nil
}

type ResultResponse struct {
	WorkLogs []workLog `json:"workLogs"`
	DateFrom time.Time `json:"dateFrom"`
	DateTo   time.Time `json:"dateTo"`
}

func MarshalResultResponse(workLogs []models.WorkLog, dateFrom, dateTo time.Time) ([]byte, error) {
	out := &ResultResponse{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		WorkLogs: make([]workLog, len(workLogs)),
	}

	for i, l := range workLogs {
		out.WorkLogs[i] = makeWorkLog(l)
	}

	return json.Marshal(out)
}
