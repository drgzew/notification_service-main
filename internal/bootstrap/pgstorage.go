package bootstrap

import (
	"fmt"
	"log"

	"notification_service/config"
	"notification_service/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	storage, err := pgstorage.NewPGStorage(connectionString)
	if err != nil {
		log.Panicf("Ошибка инициализации БД: %v", err)
	}

	return storage
}
