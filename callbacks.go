package rediy

var Callbacks = make([]func(c *Context), 0, 10)

func Use(f func(c *Context)) {
	Callbacks = append(Callbacks, f)
}
