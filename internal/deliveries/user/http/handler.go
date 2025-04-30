package http

import (
	"duo_elements/internal/app/config"
	"duo_elements/internal/data"
	usecases "duo_elements/internal/usecases/user"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(usecase usecases.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req data.CreateUserRequest

	// Логируем начало обработки запроса
	log.Println("Получен запрос на создание пользователя")

	// Логируем тело запроса
	if err := c.Bind(&req); err != nil {
		log.Printf("Ошибка привязки данных: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	log.Printf("Данные запроса: %+v", req)

	user, err := h.usecase.CreateUser(req)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Логируем успешное создание пользователя
	log.Printf("Пользователь успешно создан: %+v", user)
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.usecase.GetUserByID(intID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	log.Printf("Пользователь найден: %+v", user)
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByUsername(c echo.Context) error {
	username := c.Param("username")

	user, err := h.usecase.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdatePassword(c echo.Context) error {
	var req struct {
		Username    string `json:"username"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := h.usecase.UpdatePasswordByUsername(req.Username, req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "password updated"})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	err = h.usecase.DeleteUserByID(intID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}

func RegisterUserRoutes(e *echo.Echo, cfg *config.Config, usecase usecases.UserUsecase) {
	handler := NewUserHandler(usecase)
	e.POST("/api/v1/users", handler.CreateUser)
	e.GET("/api/v1/users/:id", handler.GetUserByID)
	e.GET("/api/v1/users/username/:username", handler.GetUserByUsername)
	e.PUT("/api/v1/users/password", handler.UpdatePassword)
	e.DELETE("/api/v1/users/:id", handler.DeleteUser)
}
