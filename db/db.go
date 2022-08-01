package db

import(
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/stdlib"

	"auctions-service/models"
)

func InitDatabase() *sqlx.DB {
	configs := models.GetDBConfigs()
	connStr := fmt.Sprintf("host=%s, port=%s, user=%s, name=%s, password=%s, sslmode=%s", 
		configs.Host, configs.Port, configs.User, configs.DBName, configs.Password, configs.SSLMode)
	
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil{
		log.Fatal(err)
	}

	return db
}