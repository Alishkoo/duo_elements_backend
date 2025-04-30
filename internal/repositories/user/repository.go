package user

import (
	"duo_elements/internal/data"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	IsUsernameTaken(username string) (bool, error)
	CreateUser(user data.User) error
	GetUserByID(id int) (*data.User, error)
	GetUserByUsername(username string) (*data.User, error)
	UpdatePasswordByUsername(username, newPasswordHash string) error
	DeleteUserByID(id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) IsUsernameTaken(username string) (bool, error) {
	log.Printf("Проверка уникальности имени пользователя: %s", username)

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := r.db.Get(&exists, query, username)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return false, err
	}

	log.Printf("Результат проверки уникальности: %v", exists)
	return exists, nil
}

func (r *userRepository) CreateUser(user data.User) error {
	log.Printf("Создание пользователя в базе данных: %+v", user)

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err := r.db.Exec(query, user.Username, user.PasswordHash)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return err
	}

	log.Println("Пользователь успешно добавлен в базу данных")
	return nil
}

func (r *userRepository) GetUserByID(id int) (*data.User, error) {
	log.Printf("Получение пользователя по ID: %d", id)

	var user data.User
	query := "SELECT username, password_hash FROM users WHERE id = $1"
	err := r.db.Get(&user, query, id)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return nil, err
	}
	log.Printf("Пользователь найден: %+v", user)
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*data.User, error) {
	log.Printf("Получение пользователя по username: %s", username)

	var user data.User
	query := "SELECT username, password_hash FROM users WHERE username = $1"
	err := r.db.Get(&user, query, username)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdatePasswordByUsername(username, newPasswordHash string) error {
	log.Printf("Обновление пароля для пользователя: %s", username)

	query := "UPDATE users SET password_hash = $1 WHERE username = $2"
	_, err := r.db.Exec(query, newPasswordHash, username)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return err
	}

	log.Println("Пароль успешно обновлен")
	return nil
}

func (r *userRepository) DeleteUserByID(id int) error {
	log.Printf("Удаление пользователя по ID: %d", id)

	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return err
	}

	log.Println("Пользователь успешно удален")
	return nil
}
