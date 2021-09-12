package main

import (
	"fmt"
	"os"
)

type sessionStats struct {
	totalKeys int
	keys      []keyOccurence
	kps       float64
}

type keyOccurence struct {
	key       rune
	occurence int
}

func main() {
	keystrokes := 10

	intervalDuration, sessionDuration := validateArguments(os.Args[1:])

	timerChannel := make(chan bool)
	tickerChannel := make(chan bool)
	statsChannel := make(chan sessionStats)
	doneChannel := make(chan sessionStats)

	go startTimer(sessionDuration, timerChannel)
	go startInterval(intervalDuration, timerChannel, tickerChannel, statsChannel, doneChannel)
	go capture(keystrokes, doneChannel, tickerChannel, statsChannel)

	stats := <-doneChannel
	printSessionStats(stats)
}

func validateArguments(args []string) (int, int) {
	count := len(args)
	if count != 2 {
		message := fmt.Sprintf("\nThis program expects 2 arguments! You have provided %d", count)
		panic(message)
	}

	intervalDuration := atoi(args[0])
	sessionDuration := atoi(args[1])

	if intervalDuration > sessionDuration {
		panic("Interval duration cannot be longer than session duration!")
	}

	return intervalDuration, sessionDuration
}
