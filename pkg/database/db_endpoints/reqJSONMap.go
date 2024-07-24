package db_endpoints

import "github.com/ethancox127/WatermarkService/internal"

type GetRequest struct {
	Filters   []internal.Filter `json:"filters,omitempty"`
}

type GetResponse struct {
	Documents []internal.Document `json:"documents"`
	Err       string              `json:"err,omitempty"`
}

type UpdateRequest struct {
	Title	  string			  `json:"title"`
	Document  *internal.Document   `json:"documents"`
}

type UpdateResponse struct {
	Err       error              `json:"err,omitempty"`
}

type AddRequest struct {
	Document  *internal.Document   `json:"documents"`
}

type AddResponse struct {
	Success   string			  `json:"success"`
	Err		  error				  `json:"err,omitempty"`
}

type RemoveRequest struct {
	Title	  string			  `json:"title"`
}

type RemoveResponse struct {
	Err		  error				  `json:"err,omitempty"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}
