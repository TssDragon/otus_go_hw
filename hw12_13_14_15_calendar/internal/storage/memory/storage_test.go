package memorystorage

import (
	"sync"
	"testing"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	firstUserUUID := uuid.New()
	thirdUserUUID := uuid.New()

	firstEvent := common.Event{
		Title:          "test title",
		DateStart:      "2022-01-01 00:00:00",
		DateEnd:        "2022-02-03 00:00:00",
		Description:    "test dscr",
		UserID:         firstUserUUID,
		TimeToRemember: 0,
	}
	secondEvent := common.Event{
		Title:          "test title 2",
		DateStart:      "2022-01-02 00:00:00",
		DateEnd:        "2022-02-02 00:00:00",
		Description:    "test dscr 2",
		UserID:         firstUserUUID,
		TimeToRemember: 0,
	}
	thirdEvent := common.Event{
		Title:          "test title 3",
		DateStart:      "2022-02-02 00:00:00",
		DateEnd:        "2022-03-03 00:00:00",
		Description:    "test dscr 3",
		UserID:         thirdUserUUID,
		TimeToRemember: 0,
	}

	storage := New()

	t.Run("create event", func(t *testing.T) {
		eventID, err := storage.Create(firstEvent)

		require.NoError(t, err)
		require.Len(t, storage.events, 1)
		require.NotEqual(t, eventID, uuid.UUID{})
	})

	t.Run("delete event", func(t *testing.T) {
		eventID := uuid.UUID{}
		for _, event := range storage.events {
			eventID = event.ID
		}
		err := storage.Delete(eventID)

		require.NoError(t, err)
		require.Len(t, storage.events, 0)
	})

	t.Run("create several events", func(t *testing.T) {
		_, err1 := storage.Create(firstEvent)
		_, err2 := storage.Create(thirdEvent)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Len(t, storage.events, 2)
	})

	t.Run("create intersection event for same user", func(t *testing.T) {
		_, err := storage.Create(secondEvent)

		require.ErrorIs(t, err, ErrDateBusy)
		require.Equal(t, len(storage.events), 2)
	})

	t.Run("update event wrong date", func(t *testing.T) {
		firstEventID := uuid.UUID{}
		for _, event := range storage.events {
			firstEventID = event.ID
			break
		}
		firstEvent.Title = "new title"
		err := storage.Update(firstEventID, firstEvent)

		require.ErrorIs(t, err, ErrDateBusy)
	})

	t.Run("update event wrong user", func(t *testing.T) {
		firstEventID := uuid.UUID{}
		for _, event := range storage.events {
			firstEventID = event.ID
			break
		}
		thirdEvent.Title = "new title 2"
		err := storage.Update(firstEventID, thirdEvent)

		require.ErrorIs(t, err, ErrUserUpdate)
	})

	t.Run("update non existing event", func(t *testing.T) {
		newEventID := uuid.New()
		err := storage.Update(newEventID, firstEvent)

		require.ErrorIs(t, err, ErrEventDoesntExists)
	})

	t.Run("update event correct", func(t *testing.T) {
		firstEventID := uuid.UUID{}
		for _, event := range storage.events {
			firstEventID = event.ID
			break
		}
		firstEvent.Title = "new title"
		firstEvent.DateStart = "2023-01-01 00:00:00"
		firstEvent.DateEnd = "2023-02-01 00:00:00"
		err := storage.Update(firstEventID, firstEvent)

		require.NoError(t, err)
		require.Equal(t, "new title", storage.events[firstEventID].Title)
	})

	t.Run("get events list on day", func(t *testing.T) {
		_, _ = storage.Create(secondEvent)
		day := "2022-02-02"
		list, err := storage.EventsListOnDate(day)

		require.NoError(t, err)
		require.Len(t, list, 2)
	})

	storage.events = nil
	storage = nil
}

func TestStorageConcurrency(t *testing.T) {
	storage := New()

	t.Run("concurrency test", func(t *testing.T) {
		wg := sync.WaitGroup{}

		for i := 1; i <= 100; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				event := common.Event{
					Title:          "test title",
					DateStart:      "2022-01-01 00:00:00",
					DateEnd:        "2022-02-03 00:00:00",
					Description:    "test dscr",
					UserID:         uuid.New(),
					TimeToRemember: 0,
				}
				id, _ := storage.Create(event)
				_ = storage.Delete(id)
			}(&wg)
		}
		wg.Wait()

		require.Len(t, storage.events, 0)
	})
}
