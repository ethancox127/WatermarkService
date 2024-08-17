package db_transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethancox127/WatermarkService/pkg/database"
	"github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"
)

func NewHTTPHandler(svc database.Service) http.Handler {
	httpsrv := newHttpServer(svc)
	if httpsrv == nil {
		return nil
	}

	m := http.NewServeMux()
	m.HandleFunc("/healthz", httpsrv.ServiceStatus)
	m.HandleFunc("/get", httpsrv.Get)
	m.HandleFunc("/add", httpsrv.Add)
	m.HandleFunc("/update", httpsrv.Update)
	m.HandleFunc("/remove", httpsrv.Remove)

	return m
}

type httpServer struct {
	svc database.Service
}

func newHttpServer(svc database.Service) *httpServer {
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

	docs, err := s.svc.Get(req.Filters...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.GetResponse{Documents: docs, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) Update(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPUpdateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.svc.Update(req.Title, &req.Document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.UpdateResponse{Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) Add(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPAddRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.svc.Add(&req.Document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.AddResponse{Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) Remove(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPRemoveRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.svc.Remove(req.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.RemoveResponse{Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) ServiceStatus(w http.ResponseWriter, r *http.Request) {
	_, err := decodeHTTPServiceStatusRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := s.svc.ServiceStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := db_endpoints.ServiceStatusResponse{Code: code, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func decodeHTTPGetRequest(r *http.Request) (db_endpoints.GetRequest, error) {
	fmt.Println("Get request")
	var req db_endpoints.GetRequest
	if r.ContentLength == 0 {
		log.Println("Get all documents")
		return req, nil
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return req, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPUpdateRequest(r *http.Request) (db_endpoints.UpdateRequest, error) {
	fmt.Println("Update request")
	var req db_endpoints.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return req, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPAddRequest(r *http.Request) (db_endpoints.AddRequest, error) {
	fmt.Println("Add request")
	var req db_endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return req, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPRemoveRequest(r *http.Request) (db_endpoints.RemoveRequest, error) {
	fmt.Println("Remove request")
	var req db_endpoints.RemoveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	fmt.Println("Request built properly")
	fmt.Println(req)
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ *http.Request) (db_endpoints.ServiceStatusRequest, error) {
	var req db_endpoints.ServiceStatusRequest
	return req, nil
}
