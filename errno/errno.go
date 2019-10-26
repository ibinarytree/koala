package errno

import (
	"fmt"
)

type KoalaError struct {
	Code int
	Message string
}

func (k *KoalaError) Error() string{
	return fmt.Sprintf("koala error, code:%d message:%v", k.Code, k.Message)
}


var (
	NotHaveInstance = &KoalaError{
		Code: 1,
		Message: "not have instance",
	}
)
