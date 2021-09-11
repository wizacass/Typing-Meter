package main

import (
	"fmt"
	"math"
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
	keystrokes := 10
	intervalDuration, sessionDuration := validateArguments(os.Args[1:])

	timerChannel := make(chan bool)
	tickerChannel := make(chan bool)
	statsChannel := make(chan sessionStats)
	doneChannel := make(chan sessionStats)

	go startTimer(sessionDuration, timerChannel)
	go startInterval(time.Duration(intervalDuration), timerChannel, tickerChannel, statsChannel, doneChannel)
	go capture(keystrokes, doneChannel, tickerChannel, statsChannel)

	stats := <-doneChannel
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

func startTimer(seconds int, timerChannel chan bool) {
	timer := time.NewTimer(time.Duration(seconds * int(time.Second)))

	go func() {
		<-timer.C
		timerChannel <- true
	}()
}

func startInterval(seconds time.Duration, timerChannel chan bool, tickerChannel chan bool, statsChannel chan sessionStats, doneChannel chan sessionStats) {
	sessionStats := sessionStats{
		totalKeys: 0,
	}

	startTime := time.Now()
	ticker := time.NewTicker(seconds * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				analyzeInterval(startTime, &sessionStats, tickerChannel, statsChannel)
			case <-timerChannel:
				analyzeInterval(startTime, &sessionStats, tickerChannel, statsChannel)

				doneChannel <- sessionStats
				break
			}
		}
	}()
}

func analyzeInterval(startTime time.Time, sessionStats *sessionStats, tickerChannel chan bool, statsChannel chan sessionStats) {
	tickerChannel <- true
	stats := <-statsChannel

	sessionStats.totalKeys += stats.totalKeys
	kps := calculateKps(startTime, sessionStats.totalKeys)

	printIntervalStats(*sessionStats, kps)
}

func calculateKps(startTime time.Time, keys int) float64 {
	interval := -startTime.Sub(time.Now()).Seconds()
	kpm := float64(keys) / interval

	return math.Round(kpm*100) / 100
}

func printIntervalStats(stats sessionStats, kps float64) {
	fmt.Println("Keys pressed during last interval:", stats.totalKeys)
	fmt.Println("Keys per second:", kps)
}

func capture(keystrokes int, doneChannel chan sessionStats, tickerChannel chan bool, statsChannel chan sessionStats) {
	stats := sessionStats{
		totalKeys: 0,
	}

	startCapturing()
	defer closeCapturing()

	keysEvents, err := keyboard.GetKeys(keystrokes)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-keysEvents:
			if event.Err != nil {
				panic(event.Err)
			}

			stats.totalKeys += 1
		case <-tickerChannel:
			statsChannel <- stats
			stats = sessionStats{
				totalKeys: 0,
			}

		case <-doneChannel:
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
