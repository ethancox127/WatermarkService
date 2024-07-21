package database

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethancox127/WatermarkService/internal"
)

type databaseService struct{}

func NewService() Service { return &databaseService{} }

func testConnection(db *sqlx.DB) error {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return err
	} else {
		log.Println("DB is connected")
		return nil
	}
}

func buildQuery(filters ...internal.Filter) string {
	query := "SELECT * FROM books "
	for i, v := range filters {
		filter := v
		if filter.Value != "" {
			if i == 0 {
				query += " WHERE " + filter.Key + " = '" + filter.Value + "'"
			} else {
				query += " AND " + filter.Key + " = '" + filter.Value + "'"
			}
		} else {
			query += " ORDER BY " + filter.Key
		}
	}
	return query
}

func receiveDocs(rows *sqlx.Rows) ([]internal.Document, error) {

	docs := []internal.Document{}
	for rows.Next() {

		doc := internal.Document{}
		if err := rows.Scan(&doc.Id, &doc.Title, &doc.Content, &doc.Author, &doc.Topic, &doc.Watermark); err != nil {
			log.Fatal(err)
			return nil, err
		}

		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return docs, nil
}

func (d *databaseService) Get(_ context.Context, db *sqlx.DB, filters ...internal.Filter) ([]internal.Document, error) {

	err := testConnection(db)
	if err != nil {
		return nil, err
	}

	query := buildQuery(filters...)

	rows, err := db.Queryx(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	docs, err := receiveDocs(rows)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (d *databaseService) Update(_ context.Context, db *sqlx.DB, title string, doc *internal.Document) error {
	err := testConnection(db)
	if err != nil {
		return err
	}

	_, err = db.Queryx("UPDATE books SET content = " + `'` + doc.Content + `'` + ", author = " + `'` + doc.Author + `'` + ", topic = " + `'` + doc.Topic + `'` + ", watermark = true WHERE title = " + `'` + doc.Title + `'`)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseService) Add(_ context.Context, db *sqlx.DB, doc *internal.Document) (string, error) {
	err := testConnection(db)
	if err != nil {
		return "", err
	}

	_, err = db.Queryx("INSERT INTO books(id, title, content, author, topic, watermark) VALUES (" + "DEFAULT, " + `'` + doc.Title + `'` + ", " + `'` + doc.Content + `'` + ", " + `'` + doc.Author + `'` + ", " + `'` + doc.Topic + `'` + ", " + "true" + ")")
	if err != nil {
		return "", err
	}

	return "Success", nil
}

func (d *databaseService) Remove(_ context.Context, db *sqlx.DB, title string) error {
	err := testConnection(db)
	if err != nil {
		return err
	}

	fmt.Println("DELETE FROM books WHERE title = " + `'` + title + `'`)
	_, err = db.Queryx("DELETE FROM books WHERE title = " + `'` + title + `'`)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseService) ServiceStatus(ctx context.Context) (int, error) {
	log.Println("Checking the Service health...")
	return http.StatusOK, nil
}
