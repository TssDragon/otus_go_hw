package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/common"
	"github.com/google/uuid"
)

const (
	dateHumanFormat     = "2006-01-02"
	dateTimeHumanFormat = "2006-01-02 15:04:05"
)

var (
	ErrDateBusy          = errors.New("нельзя создать событие на дату. время занято")
	ErrEventDoesntExists = errors.New("нельзя работать с событием, которого не существует")
	ErrUserUpdate        = errors.New("нельзя изменять пользователя события")
)

type Storage struct {
	events map[uuid.UUID]common.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: map[uuid.UUID]common.Event{},
	}
}

func (s *Storage) Create(newEvent common.Event) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.events {
		datesIntersects, err := intersectionDates(newEvent.DateStart, newEvent.DateEnd, event)
		if err != nil {
			return uuid.Nil, err
		}
		if event.UserID == newEvent.UserID && datesIntersects {
			return uuid.UUID{}, ErrDateBusy
		}
	}
	newEvent.ID = uuid.New()
	s.events[newEvent.ID] = newEvent
	return newEvent.ID, nil
}

func (s *Storage) Update(eventID uuid.UUID, newEvent common.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	event, ok := s.events[eventID]
	if !ok {
		return ErrEventDoesntExists
	}
	if event.UserID != newEvent.UserID {
		return ErrUserUpdate
	}
	datesIntersects, err := intersectionDates(newEvent.DateStart, newEvent.DateEnd, event)
	if err != nil {
		return err
	}
	if datesIntersects {
		return ErrDateBusy
	}

	s.events[eventID] = newEvent
	return nil
}

func (s *Storage) Delete(eventID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.events[eventID]
	if !ok {
		return ErrEventDoesntExists
	}
	delete(s.events, eventID)
	return nil
}

func (s *Storage) EventsListOnDate(checkDate string) ([]common.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var listToReturn []common.Event
	date, err := time.Parse(dateHumanFormat, checkDate)
	if err != nil {
		return listToReturn, err
	}
	for _, event := range s.events {
		currDateStart, err := time.Parse(dateTimeHumanFormat, event.DateStart)
		if err != nil {
			return listToReturn, err
		}
		currDateStart.Truncate(24 * time.Hour)
		currDateEnd, err := time.Parse(dateTimeHumanFormat, event.DateEnd)
		if err != nil {
			return listToReturn, err
		}
		currDateEnd.Truncate(24 * time.Hour)
		if date == currDateStart || date == currDateEnd {
			listToReturn = append(listToReturn, event)
		}
	}
	return listToReturn, nil
}

func (s *Storage) EventsListOnWeek(weekStart string) ([]common.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var listToReturn []common.Event
	weekStartDate, err := time.Parse(dateHumanFormat, weekStart)
	if err != nil {
		return listToReturn, err
	}
	weekEndSecond := weekStartDate.AddDate(0, 0, 7).Second()
	weekStartSecond := weekStartDate.Second()

	for _, event := range s.events {
		currDateTimeStart, err := time.Parse(dateTimeHumanFormat, event.DateStart)
		if err != nil {
			return listToReturn, err
		}
		currDateTimeStart.Truncate(24 * time.Hour)
		currDateTimeEnd, err := time.Parse(dateTimeHumanFormat, event.DateEnd)
		if err != nil {
			return listToReturn, err
		}
		currDateTimeEnd.Truncate(24 * time.Hour)
		if compareDateSeconds(weekStartSecond, weekEndSecond, currDateTimeStart.Second(), currDateTimeEnd.Second()) {
			listToReturn = append(listToReturn, event)
		}
	}
	return listToReturn, nil
}

func compareDateSeconds(weekStartDate int, weekEndDate int, currDateStart int, currDateEnd int) bool {
	cmp1 := currDateStart >= weekStartDate && currDateStart <= weekEndDate
	cmp2 := currDateEnd >= weekStartDate && currDateEnd <= weekEndDate
	return cmp1 || cmp2
}

func (s *Storage) EventsListOnMonth(monthStart string) ([]common.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var listToReturn []common.Event
	monthStartDate, err := time.Parse(dateHumanFormat, monthStart)
	if err != nil {
		return listToReturn, err
	}
	monthEndSecond := monthStartDate.AddDate(0, 1, 0).Second()
	monthStartSecond := monthStartDate.Second()

	for _, event := range s.events {
		currDateTime, err := time.Parse(dateTimeHumanFormat, event.DateStart)
		if err != nil {
			return listToReturn, err
		}
		if currDateTime.Second() >= monthStartSecond && currDateTime.Second() <= monthEndSecond {
			listToReturn = append(listToReturn, event)
		}
	}
	return listToReturn, nil
}

func intersectionDates(dateStart string, dateEnd string, eventToCompare common.Event) (bool, error) {
	dateTimeStart, err := time.Parse(dateTimeHumanFormat, dateStart)
	if err != nil {
		return false, err
	}
	dateTimeEnd, err := time.Parse(dateTimeHumanFormat, dateEnd)
	if err != nil {
		return false, err
	}
	eventDateTimeStart, err := time.Parse(dateTimeHumanFormat, eventToCompare.DateStart)
	if err != nil {
		return false, err
	}
	eventDateTimeEnd, err := time.Parse(dateTimeHumanFormat, eventToCompare.DateEnd)
	if err != nil {
		return false, err
	}

	before := eventDateTimeStart.Before(dateTimeStart) && eventDateTimeEnd.Before(dateTimeStart)
	after := eventDateTimeStart.After(dateTimeEnd) && eventDateTimeEnd.After(dateTimeEnd)
	return !(before || after), nil
}
