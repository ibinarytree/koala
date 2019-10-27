package errno

import (
	"fmt"
)

type KoalaError struct {
	Code    int
	Message string
}

func (k *KoalaError) Error() string {
	return fmt.Sprintf("koala error, code:%d message:%v", k.Code, k.Message)
}

var (
	NotHaveInstance = &KoalaError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &KoalaError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &KoalaError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &KoalaError{
		Code:    4,
		Message: "all node failed",
	}
)

func IsConnectError(err error) bool {

	koalaErr, ok := err.(*KoalaError)
	if !ok {
		return false
	}
	var result bool
	if koalaErr == ConnFailed {
		result = true
	}
	return result
}
