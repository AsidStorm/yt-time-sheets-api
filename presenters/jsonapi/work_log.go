package jsonapi

import (
	"encoding/json"
	"time"
	"yandex.tracker.api/domain/models"
)

type workLog struct {
	WorkLogId       int64         `json:"workLogId"`
	Duration        time.Duration `json:"duration"`
	CreatedById     string        `json:"createdById"`
	CreateByDisplay string        `json:"createByDisplay"`
	Comment         string        `json:"comment"`
	IssueKey        string        `json:"issueKey"`
	IssueDisplay    string        `json:"issueDisplay"`
	CreatedAt       time.Time     `json:"createdAt"`
	Queue           string        `json:"queue"`
	ProjectId       string        `json:"projectId"`
	ProjectName     string        `json:"projectName"`
	EpicKey         string        `json:"epicKey"`
	EpicDisplay     string        `json:"epicDisplay"`
	TypeId          string        `json:"typeId"`
	TypeKey         string        `json:"typeKey"`
	TypeDisplay     string        `json:"typeDisplay"`
}

func makeWorkLog(l models.WorkLog) workLog {
	return workLog{
		WorkLogId:       l.Id,
		Duration:        l.Duration,
		CreatedById:     l.CreatedById,
		CreateByDisplay: l.CreateByDisplay,
		Comment:         l.Comment,
		IssueKey:        l.IssueKey,
		IssueDisplay:    l.IssueDisplay,
		CreatedAt:       l.CreatedAt,
		Queue:           l.ExtractQueue(),
		ProjectId:       l.ProjectId,
		ProjectName:     l.ProjectName,
		EpicKey:         l.EpicKey,
		EpicDisplay:     l.EpicDisplay,
		TypeId:          l.TypeId,
		TypeKey:         l.TypeKey,
		TypeDisplay:     l.TypeDisplay,
	}
}

func MarshalWorkLog(in models.WorkLog) ([]byte, error) {
	out := struct {
		Data workLog `json:"data"`
	}{makeWorkLog(in)}

	return json.Marshal(out)
}

func MarshalWorkLogs(in []models.WorkLog) ([]byte, error) {
	out := struct {
		Data []workLog `json:"data"`
	}{make([]workLog, len(in))}

	for i, l := range in {
		out.Data[i] = makeWorkLog(l)
	}

	return json.Marshal(out)
}
