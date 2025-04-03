package user

import (
	"duo_elements/internal/data"
	"duo_elements/internal/repositories/user"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req data.CreateUserRequest) (*data.UserResponse, error) {
	log.Printf("Начало создания пользователя: %+v", req)

	// Проверка уникальности имени пользователя
	exists, err := s.repo.IsUsernameTaken(req.Username)
	if err != nil {
		log.Printf("Ошибка проверки уникальности имени пользователя: %v", err)
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		log.Printf("Имя пользователя уже занято: %s", req.Username)
		return nil, fmt.Errorf("username already taken")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	log.Println("Пароль успешно захеширован")

	// Создание пользователя
	user := data.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}
	if err := s.repo.CreateUser(user); err != nil {
		log.Printf("Ошибка создания пользователя в базе данных: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("Пользователь успешно создан: %+v", user)
	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (s *UserService) GetUserByID(id int) (*data.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (s *UserService) GetUserByUsername(username string) (*data.UserResponse, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &data.UserResponse{
		Username: user.Username,
	}, nil
}

func (s *UserService) UpdatePasswordByUsername(username, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.repo.UpdatePasswordByUsername(username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *UserService) DeleteUserByID(id int) error {
	err := s.repo.DeleteUserByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
