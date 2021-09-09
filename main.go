package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
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
	sessionChannel := make(chan sessionStats)

	go startTimer(sessionDuration, controlChannel)
	go capture(sessionChannel, controlChannel)

	stats := <-sessionChannel

	fmt.Println("Total keys pressed:", stats.totalKeys)
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

func startSession(seconds int) {
	timer := time.NewTimer(time.Duration(seconds * int(time.Second)))

	fmt.Println("Session started!")
	go func() {
		<-timer.C
		fmt.Println("Session ended!")
	}()
}

func capture(c chan sessionStats, timerChannel chan int) {
	stats := sessionStats{
		totalKeys: 0,
	}

	startCapturing()
	defer closeCapturing()

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}

			fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
			stats.totalKeys += 1

		case <-timerChannel:
			c <- stats
			break
		}
	}
}

func startCapturing() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
}

func closeCapturing() {
	_ = keyboard.Close()
}
