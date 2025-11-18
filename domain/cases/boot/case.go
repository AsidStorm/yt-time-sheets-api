package boot

import (
	"fmt"
	"sync"

	"yandex.tracker.api/domain"
	"yandex.tracker.api/domain/cases/boards"
	"yandex.tracker.api/domain/cases/issue_statuses"
	"yandex.tracker.api/domain/cases/issue_types"
	"yandex.tracker.api/domain/cases/my_user"
	"yandex.tracker.api/domain/cases/projects"
	"yandex.tracker.api/domain/cases/queues"
	"yandex.tracker.api/domain/cases/users_and_groups"
	"yandex.tracker.api/domain/models"
)

type Response struct {
	Me            models.User                  `json:"me"`
	Users         []models.User                `json:"users"`
	Groups        []models.Group               `json:"groups"`
	Queues        []models.Queue               `json:"queues"`
	Boards        []models.Board               `json:"boards"`
	IssueTypes    []models.DictionaryIssueType `json:"issueTypes"`
	IssueStatuses []models.IssueStatus         `json:"issueStatuses"`
	Projects      []models.Project             `json:"projects"`
}

func Run(c domain.Context) (*Response, error) {
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("unable to initialize case [boot] due [%s]", err)
	}

	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	errCh := make(chan error, 7) // Buffered channel to prevent blocking

	wg.Add(7)

	out := &Response{}

	go func() {
		response, err := my_user.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve my user due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.Me = response.User

		wg.Done()
	}()

	go func() {
		response, err := users_and_groups.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve users and groups due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.Users = response.Users
		out.Groups = response.Groups

		wg.Done()
	}()

	go func() {
		response, err := queues.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve queues due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.Queues = response.Queues

		wg.Done()
	}()

	go func() {
		response, err := boards.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve boards due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.Boards = response.Boards

		wg.Done()
	}()

	go func() {
		response, err := issue_statuses.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve issue statuses due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.IssueStatuses = response.IssueStatuses

		wg.Done()
	}()

	go func() {
		response, err := issue_types.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve issue types due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.IssueTypes = response.IssueTypes

		wg.Done()
	}()

	go func() {
		response, err := projects.Run(c)
		if err != nil {
			errCh <- fmt.Errorf("unable to retrieve projects due [%s]", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		out.Projects = response.Projects

		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func validate(c domain.Context) error {
	return domain.ValidateContext(c)
}
