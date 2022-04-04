package app

import (
	"context"

	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/common"
)

type App struct {
	storage common.Storage
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
}

func New(logger Logger, storage common.Storage) *App {
	return &App{
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
