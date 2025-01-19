package colors

import (
	"fmt"
	"log"
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
	fmt.Print(Red)

	if len(err) > 0 {
		log.Default().Print(msg, err)
	} else {
		log.Default().Print(msg)
	}

	fmt.Print(Reset)
}
