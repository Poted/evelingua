package listener

import (
	"bufio"
	"evelinqua/es"
	"evelinqua/handler"
	"fmt"
	"os"
	"strings"
	"time"
)

var startTime time.Time

func dots() {

	time.Sleep(100 * time.Millisecond)

	dots := []string{".  ", ".. ", "...\r"}
	for range dots {
		for _, dot := range dots {
			fmt.Printf("\r%s", dot)
			time.Sleep(140 * time.Millisecond)
		}
	}
}

func Listen() {

	defer func() {
		if r := recover(); r != nil {
			Listen()
		}
	}()

	startTime = time.Now()
	reader := bufio.NewReader(os.Stdin)
	text := "App is running. Type commands to control:\n"

	dots()

	for _, t := range text {

		if t == '.' {
			time.Sleep(200 * time.Millisecond)
		} else {
			time.Sleep(20 * time.Millisecond)
		}
		fmt.Print(string(t))
	}

	for {
		time.Sleep(300 * time.Millisecond)
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		f := commands[input]
		if f != nil {
			f()
		} else {
			fmt.Println("command not found")
		}
	}
}

var commands = map[string]func(){

	"status": func() {
		fmt.Println("App is running smoothly!")
		fmt.Println("Time up: ", time.Since(startTime))
	},

	"reload": func() {
		fmt.Println("Reloading listener...")
		panic("")
	},

	"exit": func() {
		os.Exit(0)
	},

	"restart-es": func() {
		es.ElasticSearchConnection()
	},

	"panic": func() {
		panic("Panic!")
	},

	"restart-handler": func() {
		handler.Restart()
	},

	"stop-handler": func() {
		handler.Stop()
	},
}
