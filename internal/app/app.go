package app

import (
	"log"

	"duo_elements/internal/app/config"
	"duo_elements/internal/app/connections"
	"duo_elements/internal/app/start"
	"duo_elements/pkg/graceful"
)

func Run(configFiles ...string) {

	cfg, err := config.NewConfig(configFiles...)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	conns, err := connections.NewConnections(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	errs := make(chan error, 100)

	grace := graceful.New(
		start.HTTP(errs, cfg),
	)

	grace.Shutdown(errs, log.Default(), conns)
}
