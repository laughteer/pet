package store

import (
	"database/sql"
	"fmt"
	"myapp/internal/config"

	_ "github.com/lib/pq"
)

// NewPostgresDB создаёт подключение к PostgreSQL
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return sql.Open("postgres", connStr)
}
