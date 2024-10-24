package yandex_tracker

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
