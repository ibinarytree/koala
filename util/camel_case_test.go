package util

import (
	"testing"
)

func TestCamelCase(t *testing.T) {
	str := "guestbook"
	t.Logf("result:%v", CamelCase(str))
}
