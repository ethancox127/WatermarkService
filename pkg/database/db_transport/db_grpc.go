package db_transport

import (
	"context"
	"fmt"

	db "github.com/ethancox127/WatermarkService/api/v1/pb/database"
	"github.com/ethancox127/WatermarkService/internal"
	endpoints "github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"
	"github.com/ethancox127/WatermarkService/pkg/database"
)

type grpcServer struct {
	db.UnimplementedDatabaseServer
	svc		database.Service
}

func NewGRPCServer(svc database.Service) db.DatabaseServer {
	return &grpcServer{svc: svc}
}

func (g *grpcServer) Get(ctx context.Context, r *db.GetRequest) (*db.GetReply, error) {
	req, err := decodeGRPCGetRequest(ctx, r)
	if err != nil {
		return &db.GetReply{Documents: nil, Err: err.Error()}, err
	}

	docs, err := g.svc.Get(req.Filters...)
	if err != nil {
		return &db.GetReply{Documents: nil, Err: err.Error()}, err
	}

	myDocs := make([]*db.Document, len(docs))
	for i, doc := range docs {
		myDocs[i] = &db.Document{Id: int32(doc.Id), Content: doc.Content, Title: doc.Title, Author: doc.Author, Topic: doc.Topic, Watermark: doc.Watermark}
	}

	return &db.GetReply{Documents: myDocs, Err: ""}, nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *db.ServiceStatusRequest) (*db.ServiceStatusReply, error) {
	_, err := decodeGRPCServiceStatusRequest(ctx, r)
	if err != nil {
		return &db.ServiceStatusReply{Code: -1, Err: err.Error()}, err
	}

	code, err := g.svc.ServiceStatus()
	if err != nil {
		return &db.ServiceStatusReply{Code: int64(code), Err: err.Error()}, err
	}

	return &db.ServiceStatusReply{Code: int64(code), Err: ""}, nil
}

func (g *grpcServer) Add(ctx context.Context, r *db.AddRequest) (*db.AddReply, error) {
	req, err := decodeGRPCAddRequest(ctx, r)
	if err != nil {
		return &db.AddReply{Err: err.Error()}, err
	}

	err = g.svc.Add(&req.Document)
	if err != nil {
		return &db.AddReply{Err: err.Error()}, err
	}

	return &db.AddReply{Err: ""}, nil
}

func (g *grpcServer) Update(ctx context.Context, r *db.UpdateRequest) (*db.UpdateReply, error) {
	req, err := decodeGRPCUpdateRequest(ctx, r)
	if err != nil {
		return &db.UpdateReply{Err: err.Error()}, err
	}

	fmt.Println(req.Document)
	err = g.svc.Update(req.Title, &req.Document)
	if err != nil {
		return &db.UpdateReply{Err: err.Error()}, err
	}

	return &db.UpdateReply{Err: ""}, nil
}

func (g *grpcServer) Remove(ctx context.Context, r *db.RemoveRequest) (*db.RemoveReply, error) {
	req, err := decodeGRPCRemoveRequest(ctx, r)
	if err != nil {
		return &db.RemoveReply{Err: err.Error()}, err
	}

	err = g.svc.Remove(req.Title)
	if err != nil {
		return &db.RemoveReply{Err: err.Error()}, err
	}

	return &db.RemoveReply{Err: ""}, nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (endpoints.GetRequest, error) {
	req := grpcReq.(*db.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (endpoints.UpdateRequest, error) {
	req := grpcReq.(*db.UpdateRequest)
	doc := internal.Document{
		Id:		   int(req.Document.Id),
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoints.UpdateRequest{Title: req.Title, Document: doc}, nil
}

func decodeGRPCRemoveRequest(_ context.Context, grpcReq interface{}) (endpoints.RemoveRequest, error) {
	req := grpcReq.(*db.RemoveRequest)
	return endpoints.RemoveRequest{Title: req.Title}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (endpoints.AddRequest, error) {
	req := grpcReq.(*db.AddRequest)
	doc := internal.Document{
		Id:		   int(req.Document.Id),
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoints.AddRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, _ interface{}) (endpoints.ServiceStatusRequest, error) {
	return endpoints.ServiceStatusRequest{}, nil
}