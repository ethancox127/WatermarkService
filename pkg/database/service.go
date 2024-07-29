package database

import (
	"github.com/ethancox127/WatermarkService/internal"
)

type Service interface {
	// Get the list of all documents
	Get(filters ...internal.Filter) ([]internal.Document, error)
	Update(title string, doc *internal.Document) (error)
	Add(doc *internal.Document) (error)
	Remove(title string) (error)
	ServiceStatus() (int, error)
}
