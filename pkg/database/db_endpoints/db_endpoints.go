package db_endpoints

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/database"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	UpdateEndpoint        endpoint.Endpoint
	AddEndpoint           endpoint.Endpoint
	RemoveEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(svc database.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(svc),
		UpdateEndpoint:        MakeUpdateEndpoint(svc),
		AddEndpoint:           MakeAddEndpoint(svc),
		RemoveEndpoint:        MakeRemoveEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

func MakeGetEndpoint(svc database.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("Make Get Endpoint")
		req := request.(GetRequest)
		docs, err := svc.Get(req.Filters...)
		if err != nil {
			return GetResponse{docs, err.Error()}, nil
		}
		return GetResponse{docs, ""}, nil
	}
}

func MakeUpdateEndpoint(svc database.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := svc.Update(req.Title, req.Document)
		if err != nil {
			return UpdateResponse{err}, nil
		}
		return UpdateResponse{nil}, nil
	}
}

func MakeAddEndpoint(svc database.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		success, err := svc.Add(req.Document)
		if err != nil {
			return AddResponse{success, err}, nil
		}
		return AddResponse{success, nil}, nil
	}
}

func MakeRemoveEndpoint(svc database.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RemoveRequest)
		err := svc.Remove(req.Title)
		if err != nil {
			return RemoveResponse{err}, nil
		}
		return RemoveResponse{nil}, nil
	}
}

func MakeServiceStatusEndpoint(svc database.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus()
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	fmt.Println("Set Get request")
	resp, err := s.GetEndpoint(ctx, GetRequest{Filters: filters})
	if err != nil {
		return []internal.Document{}, err
	}
	getResp := resp.(GetResponse)
	if getResp.Err != "" {
		return []internal.Document{}, errors.New(getResp.Err)
	}
	return getResp.Documents, nil
}

func (s *Set) Update(ctx context.Context, title string, doc *internal.Document) error {
	resp, err := s.UpdateEndpoint(ctx, UpdateRequest{Title: title, Document: doc})
	if err != nil {
		return err
	}
	updateResp := resp.(UpdateResponse)
	if updateResp.Err != nil {
		return updateResp.Err
	}
	return nil
}

func (s *Set) Add(ctx context.Context, doc *internal.Document) (success string, err error) {
	resp, err := s.AddEndpoint(ctx, AddRequest{Document: doc})
	if err != nil {
		return "Fail", err
	}
	addResp := resp.(AddResponse)
	if addResp.Err != nil {
		return "Fail", addResp.Err
	}
	return "Pass", nil
}

func (s *Set) Remove(ctx context.Context, title string) error {
	resp, err := s.RemoveEndpoint(ctx, RemoveRequest{Title: title})
	if err != nil {
		return err
	}
	removeResp := resp.(RemoveResponse)
	if removeResp.Err != nil {
		return removeResp.Err
	}
	return nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	resp, err := s.ServiceStatusEndpoint(ctx, ServiceStatusRequest{})
	svcStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return svcStatusResp.Code, err
	}
	if svcStatusResp.Err != "" {
		return svcStatusResp.Code, errors.New(svcStatusResp.Err)
	}
	return svcStatusResp.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
