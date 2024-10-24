package models

import "time"

type Sprint struct {
	Id        int64
	Name      string
	Status    string
	StartDate time.Time
	EndDate   time.Time
}
