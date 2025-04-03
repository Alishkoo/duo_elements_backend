package usecases

import "duo_elements/internal/data"

type UserUsecase interface {
	CreateUser(req data.CreateUserRequest) (*data.UserResponse, error)
	GetUserByID(id int) (*data.UserResponse, error)
	GetUserByUsername(username string) (*data.UserResponse, error)
	UpdatePasswordByUsername(username, newPassword string) error
	DeleteUserByID(id int) error
}
