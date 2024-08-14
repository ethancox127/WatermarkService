package db_transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"

	"github.com/ethancox127/WatermarkService/internal/util"
	"github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(svc database.Service) http.Handler {
	httpsrv := newHttpServer()
	if httpsrv == nil {
		return nil
	}

	m := http.NewServeMux()
	m.HandleFunc("/healthz", )
	m.HandleFunc("/get", )
	m.HandleFunc("/add", )
	m.HandleFunc("/update", )
	m.HandleFunc("/remove", )

	return m
}

type httpServer struct {
	svc database.Service
}

func newHttpServer(svc database.Server) *httpServer {
	return &httpServer{
		svc: svc,
	}
}

func (s *httpServer) Get(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPGetRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	docs, err := svc.Get(req.Filters...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.GetResponse{Documents: docs, err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("Get request")
	var req db_endpoints.GetRequest
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
	fmt.Println("Update request")
	bytedata, err := ioutil.ReadAll(r.Body)
	reqBodyString := string(bytedata)
	fmt.Println(reqBodyString)
	var req db_endpoints.UpdateRequest
	/*err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}*/
	err = json.Unmarshal(bytedata, &req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("Add request")
	var req db_endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	/*err = json.Unmarshal(bytedata, &req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}*/
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPRemoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("Remove request")
	var req db_endpoints.RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req db_endpoints.ServiceStatusRequest
	return req, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	fmt.Println("encode response")
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(&response)
}

func encodeError(e error, w http.ResponseWriter) {
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
