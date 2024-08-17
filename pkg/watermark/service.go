package watermark

import (
	"github.com/ethancox127/WatermarkService/internal"
)

type Service interface {
	// Get the list of all documents
	Get(filters ...internal.Filter) ([]internal.Document, error)
	Status(ticketID string) (internal.Status, error)
	Watermark(ticketID, mark string) (int, error)
	AddDocument(doc *internal.Document) (string, error)
	ServiceStatus() (int, error)
}
