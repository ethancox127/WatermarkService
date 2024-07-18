package main

import (
	"fmt"
	"log"
	"context"
	"testing"
	"github.com/jmoiron/sqlx"

	"github.com/ethancox127/WatermarkService/internal/db_utils"
	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/database"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB
var err error
var ctx context.Context

func TestDBConnection(t *testing.T) {
	db, err = db_utils.ConnectDB()
    if err != nil {
        log.Fatalln(err)
    }
  
	ctx = context.TODO()
}

func TestAdd(t *testing.T) {
	dbService := database.NewService()
	doc := internal.Document{Id: -1, Title: "Test", Content: "Test", Author: "Test", Topic: "Test", Watermark: "test"}
	_, err := dbService.Add(ctx, db, &doc)
	require.Equal(t, err, nil, "Error truncating Event Log")
}

func TestGet(t *testing.T) {
	filter := internal.Filter{Key: "title", Value: "Mary Had a Little Lamb"}
	dbService := database.NewService()
	docs, err := dbService.Get(ctx, db, filter)
	require.Equal(t, err, nil, "Error truncating Event Log")
	fmt.Println(docs)

	filter = internal.Filter{Key: "title"}
	docs, err = dbService.Get(ctx, db, filter)
	require.Equal(t, err, nil, "Error truncating Event Log")
	fmt.Println(docs)
}

func TestUpdate(t *testing.T) {
	dbService := database.NewService()
	
}