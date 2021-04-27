## rediy

rediy 对redis包提供切面操作。

- rediy内置了高频同key熔断(100w/10s)。该机制仅对`set` `setex` `setnx` `hset`命令有效。如不需要,设置rediy.StopFrequencyProtecting=true。

## 接入
- 对新项目
```go
pool := rediy.NewPool("localhost:6379", "", 0)
```

- 对已有项目
```git
// 1. 将pool的类型声明从 redis.Pool 修改为 rediy.RedisPoolI
- var pool *redis.Pool
+ var pool rediy.RedisPoolI
```

```go
// 2. 使用rediy.NewRedisPooly方法，来对原实例包装
- pool = xxxxnewPool("x.x.x.x:6379","pw",0)
+ pool := rediy.NewRedisPooly(xxxxnewPool("x.x.x.x:6379","pw",0))
```

## 1.Usage
### 1.1 基本操作
```go
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
```

#### 1.2 切面操作
- 切面仅对conn.Do方法生效

```go
package main

import (
	"fmt"
	"github.com/fwhezfwhez/rediy"
	"github.com/garyburd/redigo/redis"
)

func main() {
	rediy.Use(middleware1)
	rediy.Use(middleware2)
	// rediy.Use(abort)

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

```
