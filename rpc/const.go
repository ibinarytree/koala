package rpc

import (
	"time"
)

const (
	DefaultConnTimeout  = 100 * time.Millisecond
	DefaultReadTimeout  = time.Second
	DefaultWriteTimeout = time.Second
)
