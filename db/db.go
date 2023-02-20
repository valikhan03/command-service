package db

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/valikhan03/command-service/models"
)

func InitDatabase() *sqlx.DB {
	configs := models.GetDBConfigs()
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		configs.Host, configs.Port, configs.User, configs.DBName, configs.Password, configs.SSLMode)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		log.Fatalf("DB connect error: %s", err.Error())
	}

	return db
}
