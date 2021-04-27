package main

import (
	"fmt"
	"github.com/fwhezfwhez/rediy"
)

func main() {
	pool := rediy.NewPool("localhost:6379", "", 0)
	rediy.Use(reportErr)

	conn := pool.Get()
	defer conn.Close()
	conn.Do("set", "uname", "ft", "wrong-arg")   // writing style will ignore error,but will captured in reportErr
}

func reportErr(c *rediy.Context) {
	c.Next()

	if c.Reply.Err() != nil {
		fmt.Println("recv err", c.Reply.Err().Error())
		return
	}
}
