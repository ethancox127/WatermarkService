package database

import (
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/ethancox127/WatermarkService/internal"
)

type Service interface {
	// Get the list of all documents
	Get(ctx context.Context, db *sqlx.DB, filters ...internal.Filter) ([]internal.Document, error)
	Update(ctx context.Context, db *sqlx.DB, title string, doc *internal.Document) (error)
	Add(ctx context.Context, db *sqlx.DB, doc *internal.Document) (string, error)
	Remove(ctx context.Context, db *sqlx.DB, title string) (error)
	ServiceStatus(ctx context.Context) (int, error)
}
