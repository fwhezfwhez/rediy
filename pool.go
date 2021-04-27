package rediy

import (
	"github.com/garyburd/redigo/redis"
)

type RedisPoolI interface {
	Get() redis.Conn
}

type RedisPooly struct {
	Pool *redis.Pool
}

func NewRedisPooly(pool *redis.Pool) *RedisPooly {
	return &RedisPooly{
		Pool: pool,
	}
}

func NewPool(server string, password string, db int) RedisPoolI {
	p := newPool(server, password, db)
	return NewRedisPooly(p)
}

func NewOfficialPool(server string, password string, db int) *redis.Pool {
	p := newPool(server, password, db)
	return p
}

func (rp *RedisPooly) Get() redis.Conn {
	conn := rp.Pool.Get()

	return newConny(conn)
}
