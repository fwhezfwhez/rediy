package rediy

import (
	"fmt"
	"github.com/fwhezfwhez/cmap"
	"github.com/garyburd/redigo/redis"
	"runtime/debug"
	"time"
)

var rem = cmap.NewMapV2(nil, 128, 2*time.Minute)

type Reply struct {
	conn     redis.Conn
	hasValue bool
	reply    interface{}
	e        error
}

func newReply(conn redis.Conn) Reply {
	return Reply{
		conn: conn,
	}
}

func (r *Reply) setReply(reply interface{}, e error) {
	r.hasValue = true

	r.reply = reply
	r.e = e
}

// 100万/10秒 同key命令，则会在10秒内保持同key熔断
var AlertN int64 = 1000000

func AlertRedisHighFrequent(command string, key string, reply *Reply) func(c *Context) {
	return func(c *Context) {
		if !inCommand(command, []string{"set", "setex", "setnx", "hset"}) {
			c.Next()
			return
		}

		var keyinfo = fmt.Sprintf("%s:%s:%d", command, key, time.Now().Unix()/10)
		sum := rem.IncrByEx(keyinfo, 1, 10)
		if sum > AlertN {
			reply.setReply("command too frequency", fmt.Errorf("redis command too frequent '%s %s'", command, key))

			reply.conn.Close()

			HandleTooFrequentError(ErrorContext{
				Stack:     debug.Stack(),
				Key:       key,
				Command:   command,
				AlertInfo: keyinfo,
			})

			rem.Delete(keyinfo)
			c.abort()
			return
		}
		c.Next()
	}
}
