package logic

import (
	"context"

	"github.com/ibinarytree/koala/example/guestbook/dal"
	"github.com/ibinarytree/koala/example/guestbook/model"
)

type GetLeaveLogic struct {
	offset uint32
	limit  uint32
}

func NewGetLeaveLogic(offset, limit uint32) *GetLeaveLogic {
	return &GetLeaveLogic{
		offset: offset,
		limit:  limit,
	}
}

func (a *GetLeaveLogic) Execute(ctx context.Context) (result []*model.Leave, err error) {

	return dal.GetLeave(ctx, a.offset, a.limit)
}
