package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/watermark/endpoints"
)

func getAllDocs(t *testing.T) {
	d := endpoints.GetRequest{}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Get Request")

	resp, err := http.Post("http://localhost:8081/get", "application/json", b)
	require.Equal(t, err, nil, "Error completing Get Request")

	var getResp endpoints.GetResponse
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	require.Equal(t, err, nil, "Error decoding Get Response")

	fmt.Println(getResp)
}

func TestGet(t *testing.T) {
	filters := []internal.Filter{}
	d := endpoints.GetRequest{Filters: filters}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Get Request")

	resp, err := http.Post("http://localhost:8081/get", "application/json", b)
	require.Equal(t, err, nil, "Error completing Get Request")

	var getResp endpoints.GetResponse
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	require.Equal(t, err, nil, "Error decoding Get Response")

	fmt.Println(getResp)
}

func TestAddDocument(t *testing.T) {
	doc := internal.Document{
		Id:        -1,
		Title:     "Dracula 2",
		Content:   "book",
		Author:    "Bram Stoker",
		Topic:     "Fantasy and Fiction",
		Watermark: "True",
	}
	d := endpoints.AddDocumentRequest{Document: &doc}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Add Request")

	resp, err := http.Post("http://localhost:8081/addDocument", "application/json", b)
	require.Equal(t, err, nil, "Error completing add request")

	var addResp endpoints.AddDocumentResponse
	_ = json.NewDecoder(resp.Body).Decode(&addResp)

	fmt.Println(addResp)
	getAllDocs(t)
}

func TestWatermark(t *testing.T) {
	title := "Dracula 2"
	d := endpoints.WatermarkRequest{TicketID: title, Mark: "My mark"}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Update Request")

	resp, err := http.Post("http://localhost:8081/Watermark", "application/json", b)
	require.Equal(t, err, nil, "Error completing update request")

	var addResp endpoints.WatermarkResponse
	_ = json.NewDecoder(resp.Body).Decode(&addResp)

	fmt.Println(addResp)
	getAllDocs(t)
}