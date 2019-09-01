package main

import (
	"fmt"
	"math"
	"time"
)

type BucketLimit struct {
	rate       float64 //漏桶中水的漏出速率
	bucketSize float64 //漏桶最多能装的水大小
	unixNano   int64   //unix时间戳
	curWater   float64 //当前桶里面的水
}

func NewBucketLimit(rate float64, bucketSize int64) *BucketLimit {
	return &BucketLimit{
		bucketSize: float64(bucketSize),
		rate:       rate,
		unixNano:   time.Now().UnixNano(),
		curWater:   0,
	}
}

func (b *BucketLimit) reflesh() {
	now := time.Now().UnixNano()
	//时间差, 把纳秒换成秒
	diffSec := float64(now-b.unixNano) / 1000 / 1000 / 1000
	b.curWater = math.Max(0, b.curWater-diffSec*b.rate)
	b.unixNano = now
	return
}

func (b *BucketLimit) Allow() bool {
	b.reflesh()
	if b.curWater < b.bucketSize {
		b.curWater = b.curWater + 1
		return true
	}

	return false
}

func main() {

	//限速50qps, 桶大小100
	limit := NewBucketLimit(50, 100)
	m := make(map[int]bool)
	for i := 0; i < 1000; i++ {
		allow := limit.Allow()
		if allow {
			m[i] = true
			continue
		}
		m[i] = false
		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < 1000; i++ {
		fmt.Printf("i=%d allow=%v\n", i, m[i])
	}
}
