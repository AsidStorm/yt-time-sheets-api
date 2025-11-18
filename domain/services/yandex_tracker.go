package services

import (
	"time"

	"yandex.tracker.api/domain/models"
)

type YandexTracker interface {
	GetIAmToken() (string, time.Time, error)
	Ping() (bool, error)
	FilterTasks(query string) ([]models.Task, error)
	SprintTasks(boardId, sprintId int64) ([]models.Task, error)
	Projects() ([]models.Project, error)
	UsersAndGroups() ([]models.User, []models.Group, error)
	Queues() ([]models.Queue, error)
	Boards() ([]models.Board, error)
	IssueTypes() ([]models.DictionaryIssueType, error)
	Sprints(boardId int64) ([]models.Sprint, error)
	WorkLogs(userIdentities, queues []string, dateFrom, dateTo time.Time) ([]models.WorkLog, []string, error)
	IssuesByKeys(issueKeys []string) (map[string]models.Issue, error)
	PatchWorkLog(workLogId int64, issueKey string, duration string, comment string) (*models.WorkLog, error)
	DeleteWorkLog(workLogId int64, issueKey string) error
	CreateWorkLog(issueKey, userIdentity string, myUser bool, duration string, comment string, start time.Time) (*models.WorkLog, error)
	MyUser() (*models.User, error)
	CreateIssueComment(issueKey, issueComment string) error
	IssueStatuses() ([]models.IssueStatus, error)
	IssuesInStatus(statuses []string, queues, projects []string, month, year int) ([]models.Issue, error)

	FilterWorkLogs(dateFrom, dateTo time.Time, userIdentities map[string]bool, filter func(workLog models.RawWorkLog) bool) ([]models.RawWorkLog, error)
	IssuesWhereQuery(query string) ([]models.Issue, error)
	IssuesWhereFilter(queues, projects []string, issueTypes []int64) ([]models.Issue, error)
}
