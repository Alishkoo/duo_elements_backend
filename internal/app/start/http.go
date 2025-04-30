package start

import (
	"context"
	"log"
	"time"

	"duo_elements/internal/app/config"
	userHttp "duo_elements/internal/deliveries/user/http"
	"duo_elements/internal/deliveries/websocket"
	userRepo "duo_elements/internal/repositories/user"
	userService "duo_elements/internal/services/user"
	userUsecases "duo_elements/internal/usecases/user"
	"duo_elements/pkg/graceful"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HTTP(errs chan<- error, cfg *config.Config, db *sqlx.DB) graceful.Service {
	startType := "http"

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// Создаем WebSocket-сервер
	wsServer := websocket.NewWebSocketServer()

	// Регистрируем WebSocket-эндпоинт
	e.GET("/ws", func(c echo.Context) error {
		wsServer.HandleConnections(c.Response(), c.Request())
		return nil
	})

	// Используем слои
	userRepo := userRepo.NewUserRepository(db)
	userService := userService.NewUserService()
	userUsecase := userUsecases.NewUserUsecase(userRepo, userService)
	userHttp.RegisterUserRoutes(e, cfg, userUsecase)

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "test_is_good"})
	})

	go func() {
		errs <- e.Start(cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port)
	}()

	return graceful.NewService(startType, func() {
		// Завершаем работу WebSocket-сервера
		wsServer.Stop()

		// Завершаем работу HTTP-сервера
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	})
}
