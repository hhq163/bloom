package bloom_test

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hhq163/bloom"
)

func TestRedisBitSet_New_Set_Test(t *testing.T) {

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
	}
	conn := pool.Get()
	defer conn.Close()

	bitSet := bloom.NewRedisBitSet("user_filter", 512, conn)
	isSetBefore, err := bitSet.Test([]uint{0})
	if err != nil {
		t.Error("Could not test bitset in redis")
	}
	if isSetBefore {
		t.Error("Bit should not be set")
	}
	err = bitSet.Set([]uint{512})
	if err != nil {
		t.Error("Could not set bitset in redis")
	}
	isSetAfter, err := bitSet.Test([]uint{512})
	if err != nil {
		t.Error("Could not test bitset in redis")
	}
	if !isSetAfter {
		t.Error("Bit should be set")
	}
	err = bitSet.Expire(3600)
	if err != nil {
		t.Errorf("Error adding expiration to bitset: %v", err)
	}
	err = bitSet.Delete()
	if err != nil {
		t.Errorf("Error cleaning up bitset: %v", err)
	}
}
