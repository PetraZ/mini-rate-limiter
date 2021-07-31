package main

import (
	"fmt"
	"limiter/leakybucket"
	"limiter/tokenbucket"
	"time"
)

func testTokenBucket() {
	limiter := tokenbucket.NewLimiter(2.0, 2.0)
	fmt.Println(limiter.AllowRequest(2.0))
	fmt.Println(limiter.CurrrentCapacity)
	fmt.Println(limiter.LastRefillUnixTimestamp)
	time.Sleep(time.Second)
	fmt.Println(limiter.AllowRequest(1.0))
}

func testLeakyBucket() {
	limiter := leakybucket.NewLimiter(100)
	prev := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		now := limiter.Take()
		fmt.Println(i, now-prev)
		prev = now
	}
}

func main() {
	testLeakyBucket()
}

// input request --> (generate a unique key)
// there'll be a map associate keys to corresponding Bucket
// if bucket is available then go ahead and avaialbe -1
// otherwise block request return back
