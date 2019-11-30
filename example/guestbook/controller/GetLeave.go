package controller

import (
	"context"

	"github.com/ibinarytree/koala/example/guestbook/generate/guestbook"
	"github.com/ibinarytree/koala/example/guestbook/logic"
	"github.com/ibinarytree/koala/logs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetLeaveController struct {
}

//检查请求参数，如果该函数返回错误，则Run函数不会执行
func (s *GetLeaveController) CheckParams(ctx context.Context, r *guestbook.GetLeaveRequest) (err error) {

	if r.GetOffset() < 0 || r.GetLimit() <= 0 {
		err = status.Errorf(codes.InvalidArgument, "add leave failed")
		return
	}
	return
}

//SayHello函数的实现
func (s *GetLeaveController) Run(ctx context.Context, r *guestbook.GetLeaveRequest) (
	resp *guestbook.GetLeaveResponse, err error) {

	resp = &guestbook.GetLeaveResponse{}
	getLeave := logic.NewGetLeaveLogic(r.GetOffset(), r.GetLimit())
	result, err := getLeave.Execute(ctx)
	if err != nil {
		logs.Error(ctx, "get leave failed, err:%v", err)
		return
	}
	for _, one := range result {
		leave := &guestbook.Leave{
			Email:   one.Email,
			Content: one.Content,
		}
		resp.Leaves = append(resp.Leaves, leave)
	}
	return
}
