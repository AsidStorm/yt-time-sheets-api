package models

import (
	"strings"
	"time"
)

type WorkLog struct {
	Id              int64
	Duration        time.Duration
	CreatedById     string
	CreateByDisplay string
	Comment         string
	IssueKey        string
	IssueDisplay    string

	ProjectId   string
	ProjectName string

	EpicKey     string
	EpicDisplay string

	TypeId      string
	TypeKey     string
	TypeDisplay string

	CreatedAt time.Time
}

func (l WorkLog) ExtractQueue() string {
	split := strings.Split(l.IssueKey, "-")

	if len(split) > 0 {
		return split[0]
	}

	return "NONE"
}
