package start

import (
	"context"
	"log"
	"time"

	"duo_elements/internal/app/config"
	"duo_elements/pkg/graceful"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTTP(errs chan<- error, cfg *config.Config) graceful.Service {
	startType := "http"

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "test_is_good"})
	})

	go func() {
		errs <- e.Start(cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port)
	}()

	return graceful.NewService(startType, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	})
}
