package db_transport

import (
	"fmt"
	"testing"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/require"

	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"
)

func TestGet(t *testing.T) {
	filters := []internal.Filter{internal.Filter{Key: "author", Value: "JK Rowling"}}
	d := db_endpoints.GetRequest{Filters: filters}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Get Request")

	resp, err := http.Post("http://localhost:8081/get", "application/json", b)
	require.Equal(t, err, nil, "Error completing Get Request")

	var getResp db_endpoints.GetResponse
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	require.Equal(t, err, nil, "Error decoding Get Response")

	fmt.Println(getResp)
}

func TestAdd(t *testing.T) {
	doc := internal.Document{
		Id: -1,
		Title: "Dracula 2",
		Content: "book",
		Author: "Bram Stoker",
		Topic: "Fantasy and Fiction",
		Watermark: "True",
	}
	d := db_endpoints.AddRequest{Document: doc}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Add Request")

	resp, err := http.Post("http://localhost:8081/add", "application/json", b)
	require.Equal(t, err, nil, "Error completing add request")

	var addResp db_endpoints.AddResponse
	err = json.NewDecoder(resp.Body).Decode(&addResp)

	fmt.Println(addResp)
}

func TestUpdate(t *testing.T) {
	title := "Dracula 2"
	doc := internal.Document{
		Id: -1,
		Title: "Dracula 2",
		Content: "book",
		Author: "Bram Stoker 2",
		Topic: "Fantasy and Fiction",
		Watermark: "True",
	}
	d := db_endpoints.UpdateRequest{Title: title, Document: doc}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Update Request")

	resp, err := http.Post("http://localhost:8081/update", "application/json", b)
	require.Equal(t, err, nil, "Error completing update request")

	var addResp db_endpoints.AddResponse
	err = json.NewDecoder(resp.Body).Decode(&addResp)

	fmt.Println(addResp)
}

func TestRemove(t *testing.T) {
	title := "Dracula 2"
	d := db_endpoints.RemoveRequest{Title: title}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(d)
	require.Equal(t, err, nil, "Error encoding Remove Request")

	resp, err := http.Post("http://localhost:8081/remove", "application/json", b)
	require.Equal(t, err, nil, "Error truncating Event Log")

	var addResp db_endpoints.RemoveResponse
	err = json.NewDecoder(resp.Body).Decode(&addResp)

	fmt.Println(addResp)
}

