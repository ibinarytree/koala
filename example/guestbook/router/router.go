
package router
import(
	"context"

	"github.com/ibinarytree/koala/server"
	"github.com/ibinarytree/koala/meta"
	
		"github.com/ibinarytree/koala/example/guestbook/generate/guestbook"
	

	
	"github.com/ibinarytree/koala/example/guestbook/controller"

)
	
type RouterServer struct{}


func (s *RouterServer) AddLeave(ctx context.Context, r*guestbook.AddLeaveRequest)(resp*guestbook.AddLeaveResponse, err error){
	
	ctx = meta.InitServerMeta(ctx,"guestbook", "AddLeave")
	mwFunc := server.BuildServerMiddleware(mwAddLeave)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	
	resp = mwResp.(*guestbook.AddLeaveResponse)
	return
}


func mwAddLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
	
		r := request.(*guestbook.AddLeaveRequest)
		ctrl := &controller.AddLeaveController{}
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
}

func (s *RouterServer) GetLeave(ctx context.Context, r*guestbook.GetLeaveRequest)(resp*guestbook.GetLeaveResponse, err error){
	
	ctx = meta.InitServerMeta(ctx,"guestbook", "GetLeave")
	mwFunc := server.BuildServerMiddleware(mwGetLeave)
	mwResp, err := mwFunc(ctx, r)
	if err != nil {
		return
	}
	
	resp = mwResp.(*guestbook.GetLeaveResponse)
	return
}


func mwGetLeave(ctx context.Context, request interface{}) (resp interface{}, err error) {
	
		r := request.(*guestbook.GetLeaveRequest)
		ctrl := &controller.GetLeaveController{}
		err = ctrl.CheckParams(ctx, r)
		if err != nil {
			return
		}
	
		resp, err = ctrl.Run(ctx, r)
		return
}

