package colors

import (
	"fmt"
	"log"
	"runtime"
)

var Reset Colors = "\033[0m"
var Red Colors = "\033[31m"
var Green Colors = "\033[32m"
var Yellow Colors = "\033[33m"
var Blue Colors = "\033[34m"
var Magenta Colors = "\033[35m"
var Cyan Colors = "\033[36m"
var Gray Colors = "\033[37m"
var White Colors = "\033[97m"

type Colors string

func FuncInColors(color Colors, fn func()) {
	fmt.Print(color)
	fn()
	fmt.Print(Reset)
}

func LogInColors(color Colors, msg string) {
	fmt.Print(color)
	log.Default().Print(msg)
	fmt.Print(Reset)
}

func ErrInColors(msg string, err ...error) {

	_, filename, line, _ := runtime.Caller(1)

	fmt.Print(Red)
	log.Default().Print(filename, ":", line)
	fmt.Print(Reset)

	if len(err) > 0 {
		fmt.Printf("%s %v\n", msg, err)
	} else {
		fmt.Printf("%s\n", msg)
	}

}
