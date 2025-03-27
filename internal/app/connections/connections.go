package connections

import (
	"duo_elements/internal/app/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connections struct {
	DB *sqlx.DB
}

func (c *Connections) Close() {
	if c.DB != nil {
		c.DB.Close()
		fmt.Println("Соединение с БД закрыто")
	}
}

func NewConnections(cfg *config.Config) (*Connections, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	fmt.Println("Подключение к БД успешно")

	return &Connections{DB: db}, nil
}
