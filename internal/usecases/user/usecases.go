package usecases

import (
	"duo_elements/internal/data"
	"duo_elements/internal/repositories/user"
	userService "duo_elements/internal/services/user"
	"fmt"
	"log"
)

type UserUsecase interface {
	CreateUser(req data.CreateUserRequest) (*data.UserResponse, error)
	GetUserByID(id int) (*data.UserResponse, error)
	GetUserByUsername(username string) (*data.UserResponse, error)
	UpdatePasswordByUsername(username, newPassword string) error
	DeleteUserByID(id int) error
}

type UserUsecaseImpl struct {
	repo    user.UserRepository
	service *userService.UserService
}

func NewUserUsecase(repo user.UserRepository, service *userService.UserService) UserUsecase {
	return &UserUsecaseImpl{
		repo:    repo,
		service: service,
	}
}

func (uc *UserUsecaseImpl) CreateUser(req data.CreateUserRequest) (*data.UserResponse, error) {
	log.Printf("Начало создания пользователя: %+v", req)

	// Проверка уникальности имени пользователя (бизнес-правило)
	exists, err := uc.repo.IsUsernameTaken(req.Username)
	if err != nil {
		log.Printf("Ошибка проверки уникальности имени пользователя: %v", err)
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		log.Printf("Имя пользователя уже занято: %s", req.Username)
		return nil, fmt.Errorf("username already taken")
	}

	// Хеширование пароля (вызов сервиса)
	hashedPassword, err := uc.service.HashPassword(req.Password)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создание пользователя
	user := data.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}
	if err := uc.repo.CreateUser(user); err != nil {
		log.Printf("Ошибка создания пользователя в базе данных: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (uc *UserUsecaseImpl) GetUserByID(id int) (*data.UserResponse, error) {
	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (uc *UserUsecaseImpl) GetUserByUsername(username string) (*data.UserResponse, error) {
	user, err := uc.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (uc *UserUsecaseImpl) UpdatePasswordByUsername(username, newPassword string) error {
	hashedPassword, err := uc.service.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = uc.repo.UpdatePasswordByUsername(username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (uc *UserUsecaseImpl) DeleteUserByID(id int) error {
	err := uc.repo.DeleteUserByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
