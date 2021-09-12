package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

type sessionStats struct {
	totalKeys int
	keys      map[rune]int
	kps       float64
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

func startTimer(seconds int, timerChannel chan bool) {
	timer := time.NewTimer(time.Duration(seconds * int(time.Second)))

	go func() {
		<-timer.C
		timerChannel <- true
	}()
}

func startInterval(seconds time.Duration, timerChannel chan bool, tickerChannel chan bool, statsChannel chan sessionStats, doneChannel chan sessionStats) {
	sessionStats := newSessionStats()

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

	mergeStats(stats, sessionStats)
	sessionStats.kps = calculateKps(startTime, sessionStats.totalKeys)

	printIntervalStats(stats, sessionStats.kps)
}

func calculateKps(startTime time.Time, keys int) float64 {
	interval := -startTime.Sub(time.Now()).Seconds()
	kps := float64(keys) / interval

	return math.Round(kps*100) / 100
}

func findMostPolularKeys(keys map[rune]int, amount int) map[rune]int {
	if len(keys) <= amount {
		return keys
	}

	popularKeys := make(map[rune]int)

	for i := 0; i < amount; i++ {
		popularKey := findMostPopularKey(keys, popularKeys)
		popularKeys[popularKey] = keys[popularKey]
	}

	return popularKeys
}

func findMostPopularKey(keys map[rune]int, ignoreKeys map[rune]int) rune {
	popularKey := getMapKey(keys)
	for key, amount := range keys {
		_, contains := ignoreKeys[key]
		if amount > keys[popularKey] && !contains {
			popularKey = key
		}
	}

	return popularKey
}
