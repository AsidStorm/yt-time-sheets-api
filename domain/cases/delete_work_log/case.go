package delete_work_log

import (
	"errors"
	"fmt"
	"yandex.tracker.api/domain"
)

type Request struct {
	WorkLogId int64
	IssueKey  string
}

func Run(c domain.Context, r Request) error {
	if err := validate(c, r); err != nil {
		return fmt.Errorf("unable to initialize case [delete_work_log] due [%s]", err)
	}

	if err := c.Services().YandexTracker(c.Session()).DeleteWorkLog(r.WorkLogId, r.IssueKey); err != nil {
		return fmt.Errorf("unable to delete work log due [%s]", err)
	}

	return nil
}

func validate(c domain.Context, r Request) error {
	if r.WorkLogId <= 0 {
		return errors.New("request.WorkLogId must be greater than 0")
	}

	if r.IssueKey == "" {
		return errors.New("request.IssueKey is empty string")
	}

	return domain.ValidateContext(c)
}
