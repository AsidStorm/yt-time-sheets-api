package result

import (
	"fmt"
	"time"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/models"
)

type Request struct {
	UserIdentities []string
	Queues         []string
	Projects       []string
	DateFrom       time.Time
	DateTo         time.Time
	IssueTypes     []int64
	ResultGroup    models.ResultGroup
	IssuesFilter   map[string]bool
}

type Response struct {
	WorkLogs []models.WorkLog
	DateFrom time.Time
	DateTo   time.Time
}

func Run(c domain.Context, r Request) (*Response, error) {
	if err := validate(c, r); err != nil {
		return nil, fmt.Errorf("unable to initialize case [result] due [%s]", err)
	}

	response := &Response{
		DateTo:   r.DateTo,
		DateFrom: r.DateFrom,
	}

	workLogs, issuesKeys, err := c.Services().YandexTracker(c.Session()).WorkLogs(r.UserIdentities, r.Queues, r.DateFrom, r.DateTo)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve work logs due [%s]", err)
	}

	if len(issuesKeys) > 0 {
		issues, err := c.Services().YandexTracker(c.Session()).IssuesByKeys(issuesKeys)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve issues by keys due [%s] (len(issues_keys): %d)", err, len(issuesKeys))
		}

		haveProjectsFilter := len(r.Projects) > 0
		haveIssueTypesFilter := len(r.IssueTypes) > 0

		projectsMap := make(map[string]bool, len(r.Projects))

		for _, v := range r.Projects {
			projectsMap[v] = true
		}

		issueTypesMap := make(map[string]bool, len(r.IssueTypes))

		for _, v := range r.IssueTypes {
			issueTypesMap[fmt.Sprintf("%d", v)] = true
		}

		for _, log := range workLogs {
			issue := issues[log.IssueKey]

			if issue.Project != nil {
				log.ProjectId = issue.Project.Id
				log.ProjectName = issue.Project.Name
			}

			if issue.Type != nil {
				log.TypeId = issue.Type.Id
				log.TypeKey = issue.Type.Key
				log.TypeDisplay = issue.Type.Display
			}

			if issue.Epic != nil {
				log.EpicKey = issue.Epic.Key
				log.EpicDisplay = issue.Epic.Summary
			}

			if haveProjectsFilter && !projectsMap[log.ProjectId] {
				continue
			}

			if haveIssueTypesFilter && !issueTypesMap[log.TypeId] {
				continue
			}

			response.WorkLogs = append(response.WorkLogs, log)
		}
	}

	if len(r.IssuesFilter) > 0 {
		var filteredLogs []models.WorkLog

		for _, log := range response.WorkLogs {
			if !r.IssuesFilter[log.IssueKey] {
				continue
			}

			filteredLogs = append(filteredLogs, log)
		}

		response.WorkLogs = filteredLogs
	}

	return response, nil
}

func validate(c domain.Context, r Request) error {
	return domain.ValidateContext(c)
}
