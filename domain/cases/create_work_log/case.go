package create_work_log

import (
	"errors"
	"fmt"
	"time"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/helpers"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	UserIdentity string
	MyUser       bool
	IssueKey     string
	Duration     string
	Comment      string
	Date         time.Time
	IssueComment string
}

type Response struct {
	WorkLog models.WorkLog
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [create_work_log] due [%s]", err)
	}

	duration, err := helpers.DurationFromString(r.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration [%s] due [%s]", r.Duration, err)
	}

	issues, err := c.Services().YandexTracker(c.Session()).IssuesByKeys([]string{r.IssueKey})
	if err != nil {
		return nil, fmt.Errorf("unable to extract work log issue by key [%s] due [%s]", r.IssueKey, err)
	}

	issue, ok := issues[r.IssueKey]
	if !ok {
		return nil, fmt.Errorf("unable to determine issue with key [%s] due [issue not found in response]", r.IssueKey)
	}

	workLog, err := c.Services().YandexTracker(c.Session()).CreateWorkLog(r.IssueKey, r.UserIdentity, r.MyUser, duration, r.Comment, r.Date)
	if err != nil {
		return nil, fmt.Errorf("unable to create work log due [%s]", err)
	}

	if issue.Epic != nil {
		workLog.EpicDisplay = issue.Epic.Summary
		workLog.EpicKey = issue.Epic.Key
	}

	if issue.Project != nil {
		workLog.ProjectId = issue.Project.Id
		workLog.ProjectName = issue.Project.Name
	}

	if r.IssueComment != "" {
		// У нас есть комментарий, который мы хотим опубликовать
		err := c.Services().YandexTracker(c.Session()).CreateIssueComment(r.IssueKey, r.IssueComment)
		if err != nil {
			if rollbackErr := c.Services().YandexTracker(c.Session()).DeleteWorkLog(workLog.Id, r.IssueKey); rollbackErr != nil {
				return nil, fmt.Errorf("unable to create issue comment due [%s] (also unable to rollback workLog [%d] due [%s])", err, workLog.Id, rollbackErr)
			}

			return nil, fmt.Errorf("unable to create issue comment due [%s]", err)
		}
	}

	return &Response{
		WorkLog: *workLog,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if r.IssueKey == "" {
		return errors.New("request.IssueKey is empty string")
	}

	if r.Duration == "" {
		return errors.New("request.Duration is empty string")
	}

	if r.UserIdentity == "" {
		return errors.New("request.UserIdentity is empty string")
	}

	if r.Date.IsZero() {
		return errors.New("request.Date is zero")
	}

	if r.IssueComment != "" && !r.MyUser {
		return errors.New("request.MyUser is false and request.IssueComment not empty")
	}

	return domain.ValidateContext(c)
}
