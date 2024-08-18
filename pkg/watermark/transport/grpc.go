package transport

import (
	"context"

	wm "github.com/ethancox127/WatermarkService/api/v1/pb/watermark"
	"github.com/ethancox127/WatermarkService/internal"
	"github.com/ethancox127/WatermarkService/pkg/watermark/endpoints"
	"github.com/ethancox127/WatermarkService/pkg/watermark"

)

type grpcServer struct {
	wm.UnimplementedWatermarkServer
	svc		watermark.Service
}

func NewGRPCServer(svc watermark.Service) wm.WatermarkServer {
	return &grpcServer{svc: svc}
}

func (g *grpcServer) Get(ctx context.Context, r *wm.GetRequest) (*wm.GetReply, error) {
	req, err := decodeGRPCGetRequest(ctx, r)
	if err != nil {
		return &wm.GetReply{Documents: nil, Err: err.Error()}, err
	}

	docs, err := g.svc.Get(req.Filters...)
	if err != nil {
		return &wm.GetReply{Documents: nil, Err: err.Error()}, err
	}

	myDocs := make([]*wm.Document, len(docs))
	for i, doc := range docs {
		myDocs[i] = &wm.Document{Id: int32(doc.Id), Content: doc.Content, Title: doc.Title, Author: doc.Author, Topic: doc.Topic, Watermark: doc.Watermark}
	}

	return &wm.GetReply{Documents: myDocs, Err: ""}, nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *wm.ServiceStatusRequest) (*wm.ServiceStatusReply, error) {
	_, err := decodeGRPCServiceStatusRequest(ctx, r)
	if err != nil {
		return &wm.ServiceStatusReply{Code: -1, Err: err.Error()}, err
	}

	code, err := g.svc.ServiceStatus()
	if err != nil {
		return &wm.ServiceStatusReply{Code: -1, Err: err.Error()}, err
	}

	return &wm.ServiceStatusReply{Code: int64(code), Err: ""}, nil
}

func (g *grpcServer) AddDocument(ctx context.Context, r *wm.AddDocumentRequest) (*wm.AddDocumentReply, error) {
	req, err := decodeGRPCAddDocumentRequest(ctx, r)
	if err != nil {
		return &wm.AddDocumentReply{TicketID: "", Err: err.Error()}, err
	}

	ticket, err := g.svc.AddDocument(req.Document)
	if err != nil {
		return &wm.AddDocumentReply{TicketID: "", Err: err.Error()}, err
	}

	return &wm.AddDocumentReply{TicketID: ticket, Err: ""}, nil
}

func (g *grpcServer) Status(ctx context.Context, r *wm.StatusRequest) (*wm.StatusReply, error) {
	req, err := decodeGRPCStatusRequest(ctx, r)
	if err != nil {
		return &wm.StatusReply{Status: wm.StatusReply_FAILED, Err: err.Error()}, err
	}

	_, err = g.svc.Status(req.TicketID)
	if err != nil {
		return &wm.StatusReply{Status: wm.StatusReply_FAILED, Err: err.Error()}, err
	}

	return &wm.StatusReply{Status: wm.StatusReply_FINISHED, Err: ""}, nil
}

func (g *grpcServer) Watermark(ctx context.Context, r *wm.WatermarkRequest) (*wm.WatermarkReply, error) {
	req, err := decodeGRPCWatermarkRequest(ctx, r)
	if err != nil {
		return &wm.WatermarkReply{Code: -1, Err: err.Error()}, err
	}

	code, err := g.svc.Watermark(req.TicketID, req.Mark)
	if err != nil {
		return &wm.WatermarkReply{Code: -1, Err: err.Error()}, err
	}

	return &wm.WatermarkReply{Code: int64(code), Err: ""}, nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (endpoints.GetRequest, error) {
	req := grpcReq.(*wm.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.Key, Value: f.Value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (endpoints.StatusRequest, error) {
	req := grpcReq.(*wm.StatusRequest)
	return endpoints.StatusRequest{TicketID: req.TicketID}, nil
}

func decodeGRPCWatermarkRequest(_ context.Context, grpcReq interface{}) (endpoints.WatermarkRequest, error) {
	req := grpcReq.(*wm.WatermarkRequest)
	return endpoints.WatermarkRequest{TicketID: req.TicketID, Mark: req.Mark}, nil
}

func decodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (endpoints.AddDocumentRequest, error) {
	req := grpcReq.(*wm.AddDocumentRequest)
	doc := &internal.Document{
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoints.AddDocumentRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (endpoints.ServiceStatusRequest, error) {
	return endpoints.ServiceStatusRequest{}, nil
}