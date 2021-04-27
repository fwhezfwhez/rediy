package rediy

import "fmt"

type ErrorContext struct {
	E         error  // 错误
	Stack     []byte // 链路
	Key       string // key
	Command   string // 命令
	AlertInfo string // 报警key
}

var HandleTooFrequentError = func(ec ErrorContext) {
	if ec.E == nil {
		return
	}
	fmt.Println(ec.E.Error())
}
