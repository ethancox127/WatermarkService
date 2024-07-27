package db_transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ethancox127/WatermarkService/internal/util"
	"github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(ep db_endpoints.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/healthz", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeResponse,
	))
	m.Handle("/get", httptransport.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	m.Handle("/add", httptransport.NewServer(
		ep.AddEndpoint,
		decodeHTTPAddRequest,
		encodeResponse,
	))
	m.Handle("/update", httptransport.NewServer(
		ep.UpdateEndpoint,
		decodeHTTPUpdateRequest,
		encodeResponse,
	))
	m.Handle("/remove", httptransport.NewServer(
		ep.RemoveEndpoint,
		decodeHTTPRemoveRequest,
		encodeResponse,
	))

	return m
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("Get request")
	var req db_endpoints.GetRequest
	fmt.Println(r.Body)
	if r.ContentLength == 0 {
		logger.Log("Get all documents")
		return req, nil
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req db_endpoints.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req db_endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPRemoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req db_endpoints.RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req db_endpoints.ServiceStatusRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("encode response")
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(&response)
}

func encodeError(ctx context.Context, e error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch e {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": e.Error(),
	})
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
