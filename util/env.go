package util

import (
	"os"
	"strings"
)

const (
	KOALA_ENV   = "KOALA_ENV"
	PRODUCT_ENV = "product"
	TEST_ENV    = "test"
)

var (
	cur_koala_env string = TEST_ENV
)

func init() {
	cur_koala_env = strings.ToLower(os.Getenv(KOALA_ENV))
	cur_koala_env = strings.TrimSpace(cur_koala_env)

	if len(cur_koala_env) == 0 {
		cur_koala_env = TEST_ENV
	}
}

func IsProduct() bool {
	return cur_koala_env == PRODUCT_ENV
}

func IsTest() bool {
	return cur_koala_env == TEST_ENV
}

func GetEnv() string {
	return cur_koala_env
}
