package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/cases/boards"
	"yandex.tracker.api/domain/cases/config"
	"yandex.tracker.api/domain/cases/create_work_log"
	"yandex.tracker.api/domain/cases/delete_work_log"
	"yandex.tracker.api/domain/cases/filter_tasks"
	"yandex.tracker.api/domain/cases/get_i_am_token"
	"yandex.tracker.api/domain/cases/issue_statuses"
	"yandex.tracker.api/domain/cases/issue_types"
	"yandex.tracker.api/domain/cases/my_user"
	"yandex.tracker.api/domain/cases/patch_work_log"
	"yandex.tracker.api/domain/cases/ping"
	"yandex.tracker.api/domain/cases/projects"
	"yandex.tracker.api/domain/cases/queues"
	"yandex.tracker.api/domain/cases/reset_config"
	"yandex.tracker.api/domain/cases/result"
	"yandex.tracker.api/domain/cases/result_v2"
	"yandex.tracker.api/domain/cases/set_config"
	"yandex.tracker.api/domain/cases/sprint_tasks"
	"yandex.tracker.api/domain/cases/sprints"
	"yandex.tracker.api/domain/cases/users_and_groups"
	"yandex.tracker.api/presenters/jsonapi"
)

func Me(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := my_user.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalUser(response.User)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func Ping(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := ping.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalPingResponse(response.PingResult)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func ResetConfig(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	err := reset_config.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(nil)
}

func Config(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := config.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalConfigResponse(response.Config, response.HaveConfig)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func SetConfig(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalConfigRequest(in)
	if err != nil {
		return BadRequest(err)
	}

	err = set_config.Run(c, set_config.Request{
		OAuthClientId:    request.OAuthClientId,
		OrganizationId:   request.OrganizationId,
		FederationId:     request.FederationId,
		TrackerAuthUrl:   request.TrackerAuthUrl,
		AllowManualInput: request.AllowManualInput,
	})
	if err != nil {
		return InternalServerError(err)
	}

	return OK(nil)
}

func IAmToken(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := get_i_am_token.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalIAmTokenResponse(response.IAmToken, response.ExpiresAt)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func UsersAndGroups(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := users_and_groups.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalUsersAndGroups(response.Users, response.Groups)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func Boards(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := boards.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalBoards(response.Boards)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func Projects(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := projects.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalProjects(response.Projects)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func BoardSprintTasks(v map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	boardId, err := strconv.Atoi(v["boardId"])
	if err != nil {
		return BadRequest(fmt.Errorf("unable to parse board id [%s] due [%s]", v["boardId"], err))
	}

	sprintId, err := strconv.Atoi(v["sprintId"])
	if err != nil {
		return BadRequest(fmt.Errorf("unable to parse sprint id [%s] due [%s]", v["sprintId"], err))
	}

	response, err := sprint_tasks.Run(c, sprint_tasks.Request{
		BoardId:  int64(boardId),
		SprintId: int64(sprintId),
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalTasks(response.Tasks)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func BoardSprints(v map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	boardId, err := strconv.Atoi(v["boardId"])
	if err != nil {
		return BadRequest(fmt.Errorf("unable to parse board id [%s] due [%s]", v["boardId"], err))
	}

	response, err := sprints.Run(c, sprints.Request{
		BoardId: int64(boardId),
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalSprints(response.Sprints)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func Queues(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := queues.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalQueues(response.Queues)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func IssueTypes(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := issue_types.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalDictionaryIssueTypes(response.IssueTypes)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func IssueStatuses(_ map[string]string, c domain.Context, _ http.ResponseWriter, _ *http.Request) (int, []byte, error) {
	response, err := issue_statuses.Run(c)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := json.Marshal(response)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func Result(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalResultRequest(in)
	if err != nil {
		return BadRequest(err)
	}

	response, err := result.Run(c, result.Request{
		UserIdentities: request.UserIdentities,
		Queues:         request.Queues,
		DateFrom:       request.DateFrom,
		DateTo:         request.DateTo,
		Projects:       request.Projects,
		ResultGroup:    request.ResultGroup,
		IssueTypes:     request.IssueTypes,
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalResultResponse(response.WorkLogs, response.DateFrom, response.DateTo)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func ResultV2(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	var request result_v2.Request

	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	if err := json.Unmarshal(in, &request); err != nil {
		return BadRequest(err)
	}

	response, err := result_v2.Run(c, request)
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalResultResponse(response.WorkLogs, response.DateFrom, response.DateTo)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func PatchWorkLog(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalPatchWorkLogRequest(in)
	if err != nil {
		return BadRequest(err)
	}

	response, err := patch_work_log.Run(c, patch_work_log.Request{
		WorkLogId: request.WorkLogId,
		IssueKey:  strings.TrimSpace(request.IssueKey),
		Duration:  strings.ReplaceAll(strings.TrimSpace(request.Duration), " ", ""),
		Comment:   strings.TrimSpace(request.Comment),
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalWorkLog(response.WorkLog)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func CreateWorkLog(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalCreateWorkLogRequest(in)
	if err != nil {
		return BadRequest(err)
	}

	response, err := create_work_log.Run(c, create_work_log.Request{
		UserIdentity: strings.TrimSpace(request.UserIdentity),
		MyUser:       request.MyUser,
		IssueKey:     strings.TrimSpace(request.IssueKey),
		Duration:     strings.ReplaceAll(strings.TrimSpace(request.Duration), " ", ""),
		Comment:      strings.TrimSpace(request.Comment),
		Date:         request.Date,
		IssueComment: strings.TrimSpace(request.IssueComment),
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalWorkLog(response.WorkLog)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}

func DeleteWorkLog(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalDeleteWorkLogRequest(in)
	if err != nil {
		return BadRequest(err)
	}

	err = delete_work_log.Run(c, delete_work_log.Request{
		WorkLogId: request.WorkLogId,
		IssueKey:  strings.TrimSpace(request.IssueKey),
	})
	if err != nil {
		return InternalServerError(err)
	}

	return OK(nil)
}

func FilterTasks(_ map[string]string, c domain.Context, _ http.ResponseWriter, r *http.Request) (int, []byte, error) {
	in, err := io.ReadAll(r.Body)
	if err != nil {
		return BadRequest(err)
	}

	request, err := jsonapi.UnmarshalFilterTasks(in)
	if err != nil {
		return BadRequest(err)
	}

	response, err := filter_tasks.Run(c, filter_tasks.Request{
		Query: strings.TrimSpace(request.Query),
	})
	if err != nil {
		return InternalServerError(err)
	}

	out, err := jsonapi.MarshalTasks(response.Tasks)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(out)
}
