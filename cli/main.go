package main

import (
	"fmt"
	"os"
	"sync"
)

//go:generate go run ./gen/main.go

/*
ARGS -> INPUT -> CONFIG -> COMMANDs ->
							| SCAN VULNERABILITIES 	-> REPORT |
							| SCAN LICENSES			-> REPORT |
*/

func channelCmd(wg *sync.WaitGroup, reportFeed chan []Report, cmd Command) {
	result := cmd.Execute()
	reportFeed <- result
	wg.Done()
}

func main() {
	// ARGS to INPUT
	actionInput := NewActionInput(os.Args[1:])

	// INPUT to CONFIG
	config := BuildConfig(actionInput)
	fmt.Println("Config: ", config)
	// CONFIG to COMMANDS
	commands := BuildCommands(config)
	fmt.Println("Commands: ", commands)

	// execute COMMANDs
	var wg sync.WaitGroup
	// create a channel that sends type `[]Report`
	ch := make(chan []Report)
	for _, cmd := range commands {
		//cmd.Execute()
		wg.Add(1)
		go channelCmd(&wg, ch, cmd)
	}

	go func() {
		wg.Wait()
		// after all goroutines finished sending on channel, close it
		close(ch)
	}()
	// // run REPORTs
	// todo: see if can run reports concurrently
	for reports := range ch {
		for _, report := range reports {
			fmt.Println(report)
		}
	}
}
