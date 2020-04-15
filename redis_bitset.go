package bloom

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const redisMaxLength = 8 * 512 * 1024 * 1024

type Connection interface {
	Do(cmd string, args ...interface{}) (reply interface{}, err error)
	Send(cmd string, args ...interface{}) error
	Flush() error
}

type RedisBitSet struct {
	keyPrefix string
	client    *redis.Client
	m         uint
}

func NewRedisBitSet(keyPrefix string, m uint, client *redis.Client) *RedisBitSet {
	return &RedisBitSet{keyPrefix, client, m}
}

func (r *RedisBitSet) Set(offsets []uint) error {
	pipe := r.client.Pipeline()
	for _, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		pipe.SetBit(key, int64(thisOffset), 1)
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisBitSet) Test(offsets []uint) (bool, error) {
	pipe := r.client.Pipeline()
	res := make([]*redis.IntCmd, len(offsets))
	for i, offset := range offsets {
		key, thisOffset := r.getKeyOffset(offset)
		cmd := r.client.GetBit(key, int64(thisOffset))
		res[i] = cmd
	}
	_, err := pipe.Exec()
	if err != nil {
		return false, err
	}
	for _, cmd := range res {
		mark, err := cmd.Result()
		if err != nil {
			return false, err
		}
		if mark == 0 {
			return false, nil
		}
	}

	return true, nil
}

func (r *RedisBitSet) Expire(seconds uint) error {
	n := uint(0)
	for n <= uint(r.m/redisMaxLength) {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		n = n + 1
		_, err := r.client.Expire(key, time.Duration(seconds)).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisBitSet) Delete() error {
	n := uint(0)
	keys := make([]string, 0)
	for n <= uint(r.m/redisMaxLength) {
		key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
		keys = append(keys, key)
		n = n + 1
	}
	_, err := r.client.Del(strings.Join(keys, " ")).Result()
	return err
}

func (r *RedisBitSet) getKeyOffset(offset uint) (string, uint) {
	n := uint(offset / redisMaxLength)
	thisOffset := offset - n*redisMaxLength
	key := fmt.Sprintf("%s:%d", r.keyPrefix, n)
	return key, thisOffset
}
