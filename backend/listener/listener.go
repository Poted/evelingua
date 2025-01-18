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

// This function is used to animate dots while the app is starting
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

// This function is used to animate text printing
func animateText(text string) {

	for _, t := range text {

		if t == '.' {
			time.Sleep(200 * time.Millisecond)
		} else {
			time.Sleep(20 * time.Millisecond)
		}
		fmt.Print(string(t))
	}

	fmt.Print("\n")
}

// Listener that runs functions using commands passed while the app is running
func Listen() {

	// Restart after panic
	defer func() {
		if r := recover(); r != nil {
			Listen()
		}
	}()

	// Start timer to measure app uptime
	startTime = time.Now()

	// Read commands from the console
	reader := bufio.NewReader(os.Stdin)

	dots()
	animateText("App is running. Type commands to control: ")

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

// These are functions that can be run while the app is running by a Listener
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
