package queue

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

/****************************基于chan的队列,可批量取出数据,有超时设定**************************************************************************/
type BatchQueue struct {
	C              chan interface{}
	TimeoutSeconds int
}

func NewBatchQueue(len, seconds int) *BatchQueue {
	b := &BatchQueue{}
	b.C = make(chan interface{}, len)
	b.TimeoutSeconds = seconds
	return b
}
func (b *BatchQueue) Len() int {
	return len(b.C)
}

func (b *BatchQueue) Put(data interface{}) {
	b.C <- data
}

func (b *BatchQueue) Get(num int) ([]interface{}, int) {
	result := make([]interface{}, num)
	timeoutSeconds := b.TimeoutSeconds
	timer := time.After(time.Second * time.Duration(timeoutSeconds))
	counter := 0
	for i := 0; i < num; i++ {
		select {
		case d := <-b.C:
			result[i] = d
			counter++
		case <-timer:
			if counter > 0 {
				return result[:counter], counter
			}
			return nil, 0
		}
	}
	return result, counter

}

/****************************基于redis list 实现的队列**************************************************************/
type RedisQueue struct {
	RedisPool *redis.Pool
	Key       string
}

func NewRedisQueue(redisPool *redis.Pool, key string) *RedisQueue {
	r := &RedisQueue{}
	r.RedisPool = redisPool
	r.Key = key
	return r
}

func (r *RedisQueue) Put(data string) error {
	redisClient := r.RedisPool.Get()
	defer redisClient.Close()
	_, errRedis := redisClient.Do("lpush", r.Key, data)
	if errRedis != nil {
		return errRedis
	}
	return nil
}

func (r *RedisQueue) Get() ([]byte, error) {
	redisClient := r.RedisPool.Get()
	defer redisClient.Close()
	redisValue, err := redis.Bytes(redisClient.Do("rpop", r.Key))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			return nil, errors.New("无任务数据")
		}
		return nil, err
	}
	return redisValue, nil
}
