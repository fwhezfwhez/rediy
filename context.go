package rediy

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"runtime/debug"
)

const ABORT = math.MaxInt32 - 10000

type Context struct {
	Conn    redis.Conn
	Command string
	Key     string
	Args    []interface{}
	Reply   Reply

	offset   int
	handlers []func(*Context)
}

func newContext() *Context {
	return &Context{
		offset:   -1,
		handlers: make([]func(*Context), 0, 10),
	}
}
func (ctx *Context) Next() {
	ctx.offset ++
	s := len(ctx.handlers)
	for ; ctx.offset < s; ctx.offset++ {
		if !ctx.isAbort() {
			func() {
				defer func() {
					if e := recover(); e != nil {
						fmt.Printf("recover from \n %s\n", debug.Stack())
					}
				}()
				ctx.handlers[ctx.offset](ctx)
			}()
		} else {
			return
		}
	}
}
func (ctx *Context) Reset() {
	//ctx.PerRequestContext = &sync.Map{}
	ctx.offset = -1
	ctx.handlers = ctx.handlers[:0]
}

// stop middleware chain
func (ctx *Context) Abort(reply interface{}, e error) {
	ctx.Reply.reply = reply
	ctx.Reply.e = e
	ctx.offset = math.MaxInt32 - 10000
}

func (ctx *Context) abort() {
	ctx.offset = math.MaxInt32 - 10000
}

func (ctx *Context) isAbort() bool {
	if ctx.offset >= ABORT {
		return true
	}
	return false
}

func (ctx *Context) addHandler(f func(ctx *Context)) {
	ctx.handlers = append(ctx.handlers, f)
}

func (ctx Context) Info() string {
	b, _ := json.MarshalIndent(map[string]interface{}{
		"command": ctx.Command,
		"key":     ctx.Key,
		"args":    ctx.Args,
	}, "  ", "  ")
	return string(b)
}
