package common

import "github.com/google/uuid"

type Event struct {
	ID             uuid.UUID
	Title          string
	DateStart      string
	DateEnd        string
	Description    string
	UserID         uuid.UUID
	TimeToRemember uint
}
