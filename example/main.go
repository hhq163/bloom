package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/hhq163/bloom"
)

var numCount = 20000

func main() {

	// pool := &redis.Pool{
	// 	MaxIdle:     3,
	// 	IdleTimeout: 240 * time.Second,
	// 	Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
	// }
	redisGame := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1",
		DB:       0,
		PoolSize: 100,
	})

	m, k := bloom.EstimateParameters(10000000, .001) //存储10000000个key，错误率0.1%,返回的k是hash函数个数，m为位图长度
	bitSet := bloom.NewRedisBitSet("users", m, conn)
	b := bloom.New(m, k, bitSet)

	// m, k := bloom.EstimateParameters(10000000, .001)
	// b := bloom.New(m, k, bloom.NewBitSet(m))

	for i := 1; i <= 10000000; i++ {
		name := fmt.Sprintf("username_%d", i)
		b.Add([]byte(name))
	}
	nums := []int{}
	for i := 0; i < numCount; i++ {
		nums = append(nums, RandInt(1, 10000000))
	}
	var wg sync.WaitGroup
	qch, ch := make(chan struct{}), make(chan int, 200)
	count := 0

	go func() {
		for num := range ch {
			count += num
		}
		close(qch)
	}()

	for i := 0; i < numCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := time.Now()
			name := fmt.Sprintf("username_%d", nums[i])
			_, err := b.Exists([]byte(name))
			if err != nil {
				fmt.Println("Exists error:", err.Error())
			}
			cost := time.Since(start).Nanoseconds()
			//fmt.Println("result is ", exists, ", cost=", cost, " ns")
			ch <- int(cost)
		}(i)
	}
	wg.Wait()
	close(ch)
	<-qch
	fmt.Println("result is avg cost=", count/numCount, " ns")
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
