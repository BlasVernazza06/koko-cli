package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func Connect(dataSourceName string) {
	var err error
	DB, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	log.Println("Database connection established successfully")
}
