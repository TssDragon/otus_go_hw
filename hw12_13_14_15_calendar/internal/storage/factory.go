package storage

import (
	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/common"
	"github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/config"
	memorystorage "github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/TssDragon/otus_go_hw/hw_12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(config config.StorageConf) common.Storage {
	switch config.Type {
	case "SQL":
		return sqlstorage.New(config.Login, config.Password)
	default:
		return memorystorage.New()
	}
}
