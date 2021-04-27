package rediy

import (
	"github.com/garyburd/redigo/redis"
)

type Conny struct {
	Conn redis.Conn
}

func newConny(conn redis.Conn) *Conny {
	return &Conny{
		Conn: conn,
	}
}

func (c *Conny) Close() error {
	return c.Conn.Close()
}

func (c *Conny) Err() error {
	return c.Conn.Err()
}

func (c *Conny) Do(commandName string, args ...interface{}) (interface{}, error) {
	// 对args无值时，取官包的执行逻辑
	if len(args) == 0 {
		return c.Conn.Do(commandName, args...)
	}

	// 对args[0]，也就是key为非字符串时，取官包的执行逻辑
	key, isString := args[0].(string)
	if !isString {
		return c.Conn.Do(commandName, args...)
	}
	ctx := newContext()
	ctx.Caller = caller(10)
	ctx.Command = commandName
	ctx.Args = args
	ctx.Key = key
	ctx.Conn = c.Conn

	ctx.Reply = newReply(c.Conn)

	wrapF := WrapFuncWithContext(func() {
		Debugf("exec conn.Do(%s, %v)", ctx.Command, ctx.Args)
		ctx.Reply.reply, ctx.Reply.e = c.Conn.Do(commandName, args ...)
	}, ctx)

	for i, _ := range Callbacks {
		wrapF.Use(Callbacks[i])
	}

	// 	wrapF.Use(AlertRedisHighFrequent)

	wrapF.Handle()

	return ctx.Reply.reply, ctx.Reply.e
}

func (c *Conny) Send(commandName string, args ...interface{}) error {
	return c.Conn.Send(commandName, args ...)
}

func (c *Conny) Flush() error {
	return c.Conn.Flush()
}

func (c *Conny) Receive() (reply interface{}, err error) {
	return c.Conn.Receive()
}
