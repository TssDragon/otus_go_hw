package sqlstorage

import (
	"context"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/common"
	"github.com/google/uuid"
)

type Storage struct {
	events map[uuid.UUID]common.Event
}

func New(string, string) *Storage {
	return &Storage{
		events: map[uuid.UUID]common.Event{},
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Create(newEvent common.Event) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (s *Storage) Update(eventID uuid.UUID, newEvent common.Event) error {
	return nil
}

func (s *Storage) Delete(eventID uuid.UUID) error {
	return nil
}

func (s *Storage) EventsListOnDate(checkDate string) ([]common.Event, error) {
	return nil, nil
}

func (s *Storage) EventsListOnWeek(date string) ([]common.Event, error) {
	return nil, nil
}

func (s *Storage) EventsListOnMonth(date string) ([]common.Event, error) {
	return nil, nil
}
