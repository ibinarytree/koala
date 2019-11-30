package dal

import (
	"context"

	"github.com/ibinarytree/koala/example/guestbook/model"
)

var leaveStoreMgr = &LeaveStoreMgr{}

type LeaveStoreMgr struct {
	leaveList []*model.Leave
}

func (l *LeaveStoreMgr) AddLeave(ctx context.Context, leave *model.Leave) (err error) {
	l.leaveList = append(l.leaveList, leave)
	return
}

func (l *LeaveStoreMgr) GetLeave(ctx context.Context, offset, limit uint32) (result []*model.Leave, err error) {

	if offset < 0 || limit <= 0 {
		return
	}

	if offset >= uint32(len(l.leaveList)) {
		return
	}

	result = l.leaveList[offset : offset+limit]
	return
}

func AddLeave(ctx context.Context, leave *model.Leave) error {
	return leaveStoreMgr.AddLeave(ctx, leave)
}

func GetLeave(ctx context.Context, offset, limit uint32) (result []*model.Leave, err error) {
	return leaveStoreMgr.GetLeave(ctx, offset, limit)
}
