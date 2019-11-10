package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	hystrix.ConfigureCommand("koala_rpc", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	for {
		err := hystrix.Do("get_baidu", func() error {
			// talk to other services
			_, err := http.Get("https://www.baidu.com/")
			if err != nil {
				fmt.Println("get error")
				return err
			}
			return nil
		}, func(err error) error {
			fmt.Println("get an error, handle it, err:", err)
			return err
		})
		if err == nil {
			fmt.Printf("request success\n")
		}
		time.Sleep(time.Millisecond * 10)
	}

	time.Sleep(2 * time.Second) // 调用Go方法就是起了一个goroutine，这里要sleep一下，不然看不到效果
}
