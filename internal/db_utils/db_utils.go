package db_utils

import (
	"log"

	"github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func ConnectDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=newuser dbname=documents sslmode=disable password=password host=localhost")
    if err != nil {
        log.Fatalln(err)
		return nil, err
    }
  
    // Test the connection to the database
    if err := db.Ping(); err != nil {
        log.Fatal(err)
		return nil, err
    } else {
        log.Println("Successfully Connected")
		return db, nil
    }
}