package watermark

import (
	"context"
	"net/http"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ethancox127/WatermarkService/internal"
)

type watermarkService struct{
	ctx context.Context
	db  *sqlx.DB
}

func NewService(ctx context.Context, db *sqlx.DB) Service {
	return &watermarkService{ctx: ctx, db: db}
}

func testConnection(db *sqlx.DB) error {
	log.Println("Testing db connection")
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

func (w *watermarkService) Get(filters ...internal.Filter) ([]internal.Document, error) {
	log.Println("Getting records")

	err := testConnection(w.db)
	if err != nil {
		return nil, err
	}

	query := buildQuery(filters...)
	log.Println(query)

	rows, err := w.db.Queryx(query)
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

func (w *watermarkService) Status(ticketID string) (internal.Status, error) {
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(ticketID, mark string) (int, error) {
	err := testConnection(w.db)
	if err != nil {
		return -1, err
	}

	_, err = w.db.Queryx("UPDATE books SET watermark = " + `'` + mark + `'` + " WHERE title = " + `'` + ticketID + `'`)
	if err != nil {
		return -1, err
	}

	return 1, nil
}

func (w *watermarkService) AddDocument(doc *internal.Document) (string, error) {
	log.Println("Adding record")

	err := testConnection(w.db)
	if err != nil {
		return "", err
	}

	log.Println("Executing query")
	_, err = w.db.Queryx("INSERT INTO books(id, title, content, author, topic, watermark) VALUES (" + "DEFAULT, " + `'` + doc.Title + `'` + ", " + `'` + doc.Content + `'` + ", " + `'` + doc.Author + `'` + ", " + `'` + doc.Topic + `'` + ", " + `'` + doc.Watermark + `'` + ")")
	if err != nil {
		return "", err
	}

	return "", nil
}

func (w *watermarkService) ServiceStatus() (int, error) {
	log.Println("Checking the Service health...")
	return http.StatusOK, nil
}
