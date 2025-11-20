package models

import (
	"strings"
	"time"
)

type RawWorkLog struct {
	Id              int64
	Duration        time.Duration
	CreatedById     string
	CreateByDisplay string
	Comment         string
	IssueKey        string
	IssueDisplay    string
	Queue           string
	CreatedAt       time.Time
}

func CombineWorkLog(log RawWorkLog, issue Issue) WorkLog {
	wl := WorkLog{
		Id:              log.Id,
		Duration:        log.Duration,
		CreatedById:     log.CreatedById,
		IssueKey:        log.IssueKey,
		IssueDisplay:    log.IssueDisplay,
		CreateByDisplay: log.CreateByDisplay,
		Comment:         log.Comment,
		CreatedAt:       log.CreatedAt,
	}

	if issue.Epic != nil {
		wl.EpicKey = issue.Epic.Key
		wl.EpicDisplay = issue.Epic.Summary
	}

	if issue.Project != nil {
		wl.ProjectId = issue.Project.Id
		wl.ProjectName = issue.Project.Name
	}

	if issue.Type != nil {
		wl.TypeId = issue.Type.Id
		wl.TypeKey = issue.Type.Key
		wl.TypeDisplay = issue.Type.Display
	}

	return wl
}

type WorkLog struct {
	Id              int64         `json:"workLogId"`
	Duration        time.Duration `json:"duration"`
	CreatedById     string        `json:"createdById"`
	CreateByDisplay string        `json:"createByDisplay"`
	Comment         string        `json:"comment"`
	IssueKey        string        `json:"issueKey"`
	IssueDisplay    string        `json:"issueDisplay"`
	Queue           string        `json:"queue"`

	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`

	EpicKey     string `json:"epicKey"`
	EpicDisplay string `json:"epicDisplay"`

	TypeId      string `json:"typeId"`
	TypeKey     string `json:"typeKey"`
	TypeDisplay string `json:"typeDisplay"`

	CreatedAt time.Time `json:"createdAt"`
}

func (l WorkLog) ExtractQueue() string {
	split := strings.Split(l.IssueKey, "-")

	if len(split) > 0 {
		return split[0]
	}

	return "NONE"
}
