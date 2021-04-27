package rediy

import (
	"fmt"
	"github.com/fwhezfwhez/fuse"
	"github.com/garyburd/redigo/redis"
)

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

func (r *Reply) SetReply(reply interface{}, e error) {
	r.hasValue = true

	r.reply = reply
	r.e = e
}

// 10秒同key的set/setex/setnx/hset打到100万，则会熔断该key
var fs = fuse.NewFuse(1000000, 10, 10, 16)

func FuseHighFrequency(c *Context) {
	if !inCommand(c.Command, []string{"set", "setex", "setnx", "hset"}) {
		c.Next()
		return
	}

	var keyinfo = fmt.Sprintf("%s:%s", c.Command, c.Key)

	ok := fs.FuseOk(keyinfo)
	if ok {
		fs.Fail(keyinfo)
		c.Next()
		return
	}

	// Customized logger after copy yours
	fmt.Printf("%s\n", c.Info())

	c.Abort("abort for too frequent", fmt.Errorf("too frequent for %s", keyinfo))
}

func (r *Reply) Err() error {
	return r.e
}
