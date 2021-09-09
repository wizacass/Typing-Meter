package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type sessionStats struct {
	totalKeys       int
	mostPopularKeys [3]byte
}

func main() {

	intervalDuration, sessionDuration := validateArguments(os.Args[1:])

	fmt.Println("Interval duration:", intervalDuration, "seconds.")
	fmt.Println("Session duration:", sessionDuration, "seconds.")

	controlChannel := make(chan int)

	go startTimer(sessionDuration, controlChannel)

	<-controlChannel
	fmt.Println("Done!")
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

func atoi(value string) int {
	number, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	}

	return number
}

func startTimer(seconds int, timeChannel chan int) {
	timer := time.NewTimer(time.Duration(seconds * int(time.Second)))

	go func() {
		<-timer.C
		timeChannel <- 1
	}()
}
