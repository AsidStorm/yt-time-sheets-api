package models

import "time"

type Issue struct {
	Key       string
	Summary   string
	Project   *IssueProject
	Epic      *IssueEpic
	Type      *IssueType
	CreatedAt time.Time
}

type IssueProject struct {
	Id   string
	Name string
}

type IssueEpic struct {
	Key     string
	Summary string
}

type IssueType struct {
	Id      string
	Key     string
	Display string
}
