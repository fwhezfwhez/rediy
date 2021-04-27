package rediy

import (
	"fmt"
	"runtime"
)

var Mode = "pro"
var StopFrequencyProtecting = false

var (
	greenBg      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	whiteBg      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellowBg     = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	redBg        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blueBg       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magentaBg    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyanBg       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green        = string([]byte{27, 91, 51, 50, 109})
	white        = string([]byte{27, 91, 51, 55, 109})
	yellow       = string([]byte{27, 91, 51, 51, 109})
	red          = string([]byte{27, 91, 51, 49, 109})
	blue         = string([]byte{27, 91, 51, 52, 109})
	magenta      = string([]byte{27, 91, 51, 53, 109})
	cyan         = string([]byte{27, 91, 51, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

// Debugf will print content where called when MODE="debug"
func Debugf(f string, v ...interface{}) {
	if Mode == "dev" || Mode == "debug" {
		_, file, l, _ := runtime.Caller(1)
		fmt.Println(fmt.Sprintf("|%s|%s| %s:%d ",
			fmt.Sprintf("%s%s%s", magenta, "rediy", reset),
			fmt.Sprintf("%s%s%s", magenta, Mode, reset), file, l) + fmt.Sprintf(f, v...))
	}
}
