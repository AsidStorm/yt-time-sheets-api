package yandex_tracker

import (
	"fmt"
	"time"

	"yandex.tracker.api/domain/models"
)

type issue struct {
	Key     string `json:"key"`
	Summary string `json:"summary"`
	Project *struct {
		Id      string `json:"id"`
		Display string `json:"display"`
	} `json:"project"`
	CreatedAt string `json:"createdAt"`
	Spent     string `json:"spent"`
	Epic      *struct {
		Key     string `json:"key"`
		Display string `json:"display"`
	} `json:"epic"`
	Type *struct {
		Id      string `json:"id"`
		Key     string `json:"key"`
		Display string `json:"display"`
	} `json:"type"`
}

func makeDomainIssue(issue issue) (*models.Issue, error) {
	layout := "2006-01-02T15:04:05.000-0700"

	createdAt, err := time.Parse(layout, issue.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse issue [%s] created at [%s] due [%s]", issue.Key, issue.CreatedAt, err)
	}

	is := &models.Issue{
		Key:       issue.Key,
		Summary:   issue.Summary,
		CreatedAt: createdAt,
	}

	if issue.Epic != nil {
		is.Epic = &models.IssueEpic{
			Key:     issue.Epic.Key,
			Summary: issue.Epic.Display,
		}
	}

	if issue.Project != nil {
		is.Project = &models.IssueProject{
			Id:   issue.Project.Id,
			Name: issue.Project.Display,
		}
	}

	if issue.Type != nil {
		is.Type = &models.IssueType{
			Id:      issue.Type.Id,
			Key:     issue.Type.Key,
			Display: issue.Type.Display,
		}
	}

	return is, nil
}
