package main

import (
	"fmt"
	"github.com/fwhezfwhez/fuse"
	"github.com/fwhezfwhez/rediy"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

func main() {
	pool := rediy.NewPool("localhost:6379", "", 0)
	rediy.Use(fusemiddleware)

	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		go func() {
			conn := pool.Get()
			defer conn.Close()
			if _, e := redis.String(conn.Do("set", "uname", "ft")); e != nil {
				fmt.Println(e.Error())
				return
			}
			fmt.Println("set ok")
		}()
	}
}

var fs = fuse.NewFuse(1, 5, 5, 16)

func fusemiddleware(c *rediy.Context) {
	if !strIn(c.Command, []string{"set", "setex", "setnx", "hset"}) {
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

	// todo, 增加超频预警/日志
	fmt.Printf("%s\n", c.Info())

	c.Abort("abort for too frequent", fmt.Errorf("too frequent for %s", keyinfo))
}

func strIn(str string, arr []string) bool {
	for _, v := range arr {
		if strings.ToUpper(v) == strings.ToUpper(str) {
			return true
		}
	}
	return false
}
