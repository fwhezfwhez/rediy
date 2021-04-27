package main

import (
	"fmt"
	"github.com/fwhezfwhez/rediy"
	"github.com/garyburd/redigo/redis"
)

func main() {
	rediy.Use(middleware1)
	rediy.Use(middleware2)
	rediy.Use(abort)

	rediy.Mode = "debug"

	pool := rediy.NewPool("localhost:6379", "", 0)
	conn := pool.Get()
	defer conn.Close()

	rs, e := redis.String(conn.Do("set", "uname", "ft"))
	if e != nil {
		fmt.Println(e)
		return
	}

	rs, e = redis.String(conn.Do("get", "uname"))

	if e != nil {
		fmt.Println(e)
		return
	}

	fmt.Println(rs)
}

func middleware1(c *rediy.Context) {
	fmt.Printf("start middleware1 %v %v \n", c.Command, c.Args)
	c.Next()
	fmt.Printf("end middleware1 %v %v \n", c.Command, c.Args)
}

func middleware2(c *rediy.Context) {
	fmt.Printf("start middleware2 %v %v \n", c.Command, c.Args)
	c.Next()
	fmt.Printf("end middleware2 %v %v \n", c.Command, c.Args)
}

func abort(c *rediy.Context) {
	c.Abort("command stops", fmt.Errorf("aborts! %v %v", c.Command, c.Args))
	return
}
