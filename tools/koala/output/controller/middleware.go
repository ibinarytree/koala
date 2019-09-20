package controller

/*
import (
	"context"
	"fmt"
	"time"

	"github.com/ibinarytree/koala/middleware"
)

func init() {
	middleware.Use(CostMiddleware)
}

func CostMiddleware(next middleware.MiddlewareFunc) middleware.MiddlewareFunc {

	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		startTimeNano := time.Now().UnixNano()
		resp, err = next(ctx, req)
		endTimeNano := time.Now().UnixNano()

		cost := (endTimeNano - startTimeNano) / 1000

		fmt.Printf("cost:%d ms\n", cost/1000)
		return
	}
}
*/
