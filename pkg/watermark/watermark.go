package watermark

import (
	"net/http"
	"os"

	"github.com/ethancox127/WatermarkService/internal"
	"github.com/go-kit/log"

	"github.com/lithammer/shortuuid/v3"
)

type watermarkService struct{}

func NewService() Service { return &watermarkService{} }

func (w *watermarkService) Get(filters ...internal.Filter) ([]internal.Document, error) {
	// query the database using the filters and return the list of documents
	// return error if the filter (key) is invalid and also return error if no item found
	doc := internal.Document{
		Content: "book",
		Title:   "Harry Potter and Half Blood Prince",
		Author:  "J.K. Rowling",
		Topic:   "Fiction and Magic",
	}
	return []internal.Document{doc}, nil
}

func (w *watermarkService) Status(ticketID string) (internal.Status, error) {
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(ticketID, mark string) (int, error) {
	return http.StatusOK, nil
}

func (w *watermarkService) AddDocument(doc *internal.Document) (string, error) {
	newTicketID := shortuuid.New()
	return newTicketID, nil
}

func (w *watermarkService) ServiceStatus() (int, error) {
	logger.Log("Checking the Service health...")
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
