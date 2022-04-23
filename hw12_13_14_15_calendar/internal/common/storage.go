package common

import "github.com/google/uuid"

type Storage interface {
	Create(event Event) (uuid.UUID, error)
	Update(eventID uuid.UUID, newEvent Event) error
	Delete(eventID uuid.UUID) error
	EventsListOnDate(checkDate string) ([]Event, error)
	EventsListOnWeek(date string) ([]Event, error)
	EventsListOnMonth(date string) ([]Event, error)
}
