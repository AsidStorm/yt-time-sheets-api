package v1

import (
	"net/http"

	"yandex.tracker.api/domain"

	"github.com/gorilla/mux"
)

const RouterPrefix = "/api/v1"

func Register(r *mux.Router, c domain.Context) {
	router := r.PathPrefix(RouterPrefix).Subrouter()

	handle := wrapHandler(c, router)

	handle("/boot", http.MethodGet, Boot, false)

	handle("/ping", http.MethodGet, Ping, false)
	handle("/config", http.MethodGet, Config, true)
	handle("/config", http.MethodPost, SetConfig, true)
	handle("/reset_config", http.MethodGet, ResetConfig, true)
	handle("/i_am_token", http.MethodGet, IAmToken, false)
	handle("/users_and_groups", http.MethodGet, UsersAndGroups, false)
	handle("/queues", http.MethodGet, Queues, false)
	handle("/issue_types", http.MethodGet, IssueTypes, false)
	handle("/issue_statuses", http.MethodGet, IssueStatuses, false)
	handle("/result", http.MethodPost, Result, false)
	handle("/result_v2", http.MethodPost, ResultV2, false)
	handle("/boards", http.MethodGet, Boards, false)
	handle("/projects", http.MethodGet, Projects, false)
	handle("/boards/{boardId}/sprints", http.MethodGet, BoardSprints, false)
	handle("/boards/{boardId}/sprints/{sprintId}/tasks", http.MethodGet, BoardSprintTasks, false)
	handle("/work_logs", http.MethodPatch, PatchWorkLog, false)
	handle("/work_logs", http.MethodDelete, DeleteWorkLog, false)
	handle("/me", http.MethodGet, Me, false)
	handle("/work_logs", http.MethodPost, CreateWorkLog, false)
	handle("/filter_tasks", http.MethodPost, FilterTasks, false)

	handle("/result_v3", http.MethodPost, ResultV3, false)
}
