package main

import (
	"dburriss/impilo_gh/domain"
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

func channelCmd(wg *sync.WaitGroup, reportFeed chan []domain.Report, cmd domain.Command) {
	result := cmd.Execute()
	reportFeed <- result
	wg.Done()
}

func main() {
	// ARGS to INPUT
	actionInput := domain.NewActionInput(os.Args)
	fmt.Printf("Input: %+v\n", actionInput)
	// INPUT to CONFIG
	config := BuildConfig(actionInput)
	fmt.Printf("Config: %+v\n", config)
	// CONFIG to COMMANDS
	commands := BuildCommands(config)
	fmt.Printf("Commands: %+v\n", commands)

	// execute COMMANDs
	var wg sync.WaitGroup
	// create a channel that sends type `[]Report`
	ch := make(chan []domain.Report)
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
			report.Run()
		}
	}
}