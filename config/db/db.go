package db

import (
	env "AuthInGo/config/env"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = setupDB()

	if err != nil {
		log.Fatal("Error connecting to db", err)
	}
}

func setupDB() (*sql.DB, error) {
	cfg := mysql.NewConfig()

	cfg.User = env.GetString("DB_USER", "root")
	cfg.Addr = env.GetString("ADDR", "127.0.0.1:3306")
	cfg.Passwd = env.GetString("DB_PASSWORD", "rootpassword")
	cfg.DBName = env.GetString("DBName", "auth_dev")
	cfg.Net = env.GetString("DB_NET", "tcp")

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Println("Error connecting to db : ", err)
		return nil, err
	}

	pingErr := db.Ping()

	if pingErr != nil {
		fmt.Println("Error pinging db : ", pingErr)
		return nil, pingErr
	}

	fmt.Println("✅ Successfully connected to db :", cfg.DBName)

	return db, nil
}