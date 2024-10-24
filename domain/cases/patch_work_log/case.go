package patch_work_log

import (
	"errors"
	"fmt"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/helpers"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	WorkLogId int64
	IssueKey  string
	Duration  string
	Comment   string
}

type Response struct {
	WorkLog models.WorkLog
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [patch_work_log] due [%s]", err)
	}

	duration, err := helpers.DurationFromString(r.Duration)
	if err != nil {
		return nil, fmt.Errorf("unable to parse duration [%s] due [%s]", r.Duration, err)
	}

	workLog, err := c.Services().YandexTracker(c.Session()).PatchWorkLog(r.WorkLogId, r.IssueKey, duration, r.Comment)
	if err != nil {
		return nil, fmt.Errorf("unable to patch work log due [%s]", err)
	}

	return &Response{
		WorkLog: *workLog,
	}, nil
}

func validate(c domain.Context, r Request) error {
	if r.WorkLogId <= 0 {
		return errors.New("request.WorkLogId must be greater than 0")
	}

	if r.IssueKey == "" {
		return errors.New("request.IssueKey is empty string")
	}

	if r.Duration == "" {
		return errors.New("request.Duration is empty string")
	}

	return domain.ValidateContext(c)
}
