package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethancox127/WatermarkService/pkg/watermark"
	"github.com/ethancox127/WatermarkService/pkg/watermark/endpoints"
)

func NewHTTPHandler(svc watermark.Service) http.Handler {
	httpsrv := newHttpServer(svc)
	if httpsrv == nil {
		return nil
	}

	m := http.NewServeMux()
	m.HandleFunc("/healthz", httpsrv.ServiceStatus)
	m.HandleFunc("/get", httpsrv.Get)
	m.HandleFunc("/addDocument", httpsrv.AddDocument)
	m.HandleFunc("/watermark", httpsrv.Watermark)
	m.HandleFunc("/status", httpsrv.Status)

	return m
}

type httpServer struct {
	svc watermark.Service
}

func newHttpServer(svc watermark.Service) *httpServer {
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

	res := endpoints.GetResponse{Documents: docs, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) Watermark(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPWatermarkRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := s.svc.Watermark(req.TicketID, req.Mark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := endpoints.WatermarkResponse{Code: code, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) AddDocument(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPAddDocumentRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := s.svc.AddDocument(req.Document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := endpoints.AddDocumentResponse{TicketID: ticket, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) Status(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPStatusRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := s.svc.Status(req.TicketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := endpoints.StatusResponse{Status: status, Err: ""}
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

	res := endpoints.ServiceStatusResponse{Code: code, Err: ""}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func decodeHTTPGetRequest(r *http.Request) (endpoints.GetRequest, error) {
	var req endpoints.GetRequest
	if r.ContentLength == 0 {
		log.Println("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func decodeHTTPStatusRequest(r *http.Request) (endpoints.StatusRequest, error) {
	var req endpoints.StatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func decodeHTTPWatermarkRequest(r *http.Request) (endpoints.WatermarkRequest, error) {
	var req endpoints.WatermarkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func decodeHTTPAddDocumentRequest(r *http.Request) (endpoints.AddDocumentRequest, error) {
	var req endpoints.AddDocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ *http.Request) (endpoints.ServiceStatusRequest, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}
