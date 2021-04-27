package main

import (
	"fmt"
	"github.com/fwhezfwhez/rediy"
	"github.com/garyburd/redigo/redis"
)

func main() {
	pool := rediy.NewPool("localhost:6379", "", 0)

	conn := pool.Get()
	defer conn.Close()

	conn.Do("set", "uname", "ft")

	rs, e := redis.String(conn.Do("get", "uname"))

	if e != nil {
		fmt.Println(e)
		return
	}

	fmt.Println(rs)
}
