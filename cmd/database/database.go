package main

import (
	"log"
	"context"

	"github.com/ethancox127/WatermarkService/internal/db_utils"
	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/database"
)

func main() {
	db, err := db_utils.ConnectDB()
    if err != nil {
        log.Fatalln(err)
    }
  
    defer db.Close()
	ctx := context.TODO()

	filter := internal.Filter{Key: "title", Value: "Mary Had a Little Lamb"}
	dbService := database.NewService()
	docs, err := dbService.Get(ctx, db, filter)
	if err == nil {
		log.Println("test")
		log.Println(docs)
	}

	filter = internal.Filter{Key: "title"}
	docs, err = dbService.Get(ctx, db, filter)
	if err == nil {
		log.Println(docs)
	}
}