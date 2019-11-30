
package guestbookc


import (
	"context"
	"fmt"

	"github.com/ibinarytree/koala/example/guestbook/client/generate/guestbook"
	
	"github.com/ibinarytree/koala/rpc"
	"github.com/ibinarytree/koala/errno"
	"github.com/ibinarytree/koala/meta"
	
)

type GuestbookClient struct {
	serviceName string
	client *rpc.KoalaClient
}

func NewGuestbookClient(serviceName string, opts...rpc.RpcOptionFunc) *GuestbookClient {
	c :=  &GuestbookClient{
		serviceName: serviceName,
	}
	c.client = rpc.NewKoalaClient(serviceName, opts...)
	return c
}


func (s *GuestbookClient) AddLeave(ctx context.Context, r*guestbook.AddLeaveRequest)(resp*guestbook.AddLeaveResponse, err error){
	/*
	middlewareFunc := rpc.BuildClientMiddleware(mwClientAddLeave)
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
*/
	mkResp, err := s.client.Call(ctx, "AddLeave", r, mwClientAddLeave)
	if err != nil {
		return nil, err
	}
	resp, ok := mkResp.(*guestbook.AddLeaveResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *guestbook.AddLeaveResponse")
		return nil, err
	}
	
	return resp, err
}


func mwClientAddLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}*/
	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*guestbook.AddLeaveRequest)
	client := guestbook.NewGuestBookServiceClient(rpcMeta.Conn)
	return client.AddLeave(ctx, req)
}

func (s *GuestbookClient) GetLeave(ctx context.Context, r*guestbook.GetLeaveRequest)(resp*guestbook.GetLeaveResponse, err error){
	/*
	middlewareFunc := rpc.BuildClientMiddleware(mwClientGetLeave)
	mkResp, err := middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}
*/
	mkResp, err := s.client.Call(ctx, "GetLeave", r, mwClientGetLeave)
	if err != nil {
		return nil, err
	}
	resp, ok := mkResp.(*guestbook.GetLeaveResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *guestbook.GetLeaveResponse")
		return nil, err
	}
	
	return resp, err
}


func mwClientGetLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
	/*
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect: %v", err)
		return nil, err
	}*/
	rpcMeta := meta.GetRpcMeta(ctx)
	if rpcMeta.Conn == nil {
		return nil, errno.ConnFailed
	}

	req := request.(*guestbook.GetLeaveRequest)
	client := guestbook.NewGuestBookServiceClient(rpcMeta.Conn)
	return client.GetLeave(ctx, req)
}

