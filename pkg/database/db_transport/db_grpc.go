package db_transport

import (
	"context"
	"fmt"
	"reflect"

	database "github.com/ethancox127/WatermarkService/api/v1/pb/database"
	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/database/db_endpoints"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpctransport.Handler
	update        grpctransport.Handler
	add           grpctransport.Handler
	remove        grpctransport.Handler
	serviceStatus grpctransport.Handler
}

func NewGRPCServer(ep db_endpoints.Set) database.DatabaseServer {
	return &grpcServer{
		get: grpctransport.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			encodeGRPCGetResponse,
		),
		update: grpctransport.NewServer(
			ep.UpdateEndpoint,
			decodeGRPCUpdateRequest,
			encodeGRPCUpdateResponse,
		),
		add: grpctransport.NewServer(
			ep.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
		),
		remove: grpctransport.NewServer(
			ep.RemoveEndpoint,
			decodeGRPCRemoveRequest,
			encodeGRPCRemoveResponse,
		),
		serviceStatus: grpctransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			encodeGRPCServiceStatusResponse,
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *database.GetRequest) (*database.GetResponse, error) {
	fmt.Println("Get")
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*database.GetResponse), nil
}

func (g *grpcServer) Update(ctx context.Context, r *database.UpdateRequest) (*database.UpdateResponse, error) {
	_, rep, err := g.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*database.UpdateResponse), nil
}

func (g *grpcServer) Add(ctx context.Context, r *database.AddRequest) (*database.AddResponse, error) {
	_, rep, err := g.add.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*database.AddResponse), nil
}

func (g *grpcServer) Remove(ctx context.Context, r *database.RemoveRequest) (*database.RemoveResponse, error) {
	_, rep, err := g.remove.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*database.RemoveResponse), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *database.ServiceStatusRequest) (*database.ServiceStatusResponse, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*database.ServiceStatusResponse), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	fmt.Println("decode Get request")
	req := grpcReq.(*database.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return db_endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*database.UpdateRequest)
	doc := internal.Document{
		Id:        int(req.Document.Id),
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return db_endpoints.UpdateRequest{Title: req.Title, Document: doc}, nil
}

func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*database.AddRequest)
	doc := internal.Document{
		Id:        int(req.Document.Id),
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return db_endpoints.AddRequest{Document: doc}, nil
}

func decodeGRPCRemoveRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*database.RemoveRequest)
	return db_endpoints.RemoveRequest{Title: req.Title}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return db_endpoints.ServiceStatusRequest{}, nil
}

func encodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	fmt.Println("decode Get response")
	fmt.Println(reflect.TypeOf(grpcReply))
	switch reply := grpcReply.(type) {
	case *db_endpoints.GetResponse:
		//reply, _ := grpcReply.(*database.GetResponse)
		var docs []internal.Document
		for _, d := range reply.Documents {
			doc := internal.Document{
				Id:        int(d.Id),
				Content:   d.Content,
				Title:     d.Title,
				Author:    d.Author,
				Topic:     d.Topic,
				Watermark: d.Watermark,
			}
			docs = append(docs, doc)
		}
		return db_endpoints.GetResponse{Documents: docs, Err: reply.Err}, nil
	}
	return db_endpoints.GetResponse{Documents: nil, Err: ""}, nil
}

func encodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*database.UpdateResponse)
	return db_endpoints.UpdateResponse{Err: reply.Err}, nil
}

func encodeGRPCAddResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*database.AddResponse)
	return db_endpoints.AddResponse{Err: reply.Err}, nil
}

func encodeGRPCRemoveResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*database.RemoveResponse)
	return db_endpoints.RemoveResponse{Err: reply.Err}, nil
}

func encodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*database.ServiceStatusResponse)
	return db_endpoints.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}
