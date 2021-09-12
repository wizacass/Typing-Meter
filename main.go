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
	totalKeys int
	keys      map[rune]int
	kps       float64
}

func newSessionStats() sessionStats {
	return sessionStats{
		totalKeys: 0,
		keys:      make(map[rune]int),
		kps:       0,
	}
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

func mergeStats(from sessionStats, to *sessionStats) {
	for key := range from.keys {
		mergeKey(key, from, to)
	}
	to.totalKeys += from.totalKeys
}

func mergeKey(key rune, from sessionStats, to *sessionStats) {
	if _, ok := to.keys[key]; ok {
		to.keys[key] += from.keys[key]
	} else {
		to.keys[key] = from.keys[key]
	}
}

func calculateKps(startTime time.Time, keys int) float64 {
	interval := -startTime.Sub(time.Now()).Seconds()
	kps := float64(keys) / interval

	return math.Round(kps*100) / 100
}

func printIntervalStats(stats sessionStats, kps float64) {
	fmt.Println("\n--- Interval stats ---")
	printStats(stats, kps)
	fmt.Println()
}

func printStats(stats sessionStats, kps float64) {
	fmt.Println("\nKeys pressed during last interval:", stats.totalKeys)
	printKeys(stats.keys)
	fmt.Println("\nCurrent Session typing speed:", kps, "kps")
}

func printKeys(keys map[rune]int) {
	if len(keys) == 0 {
		fmt.Println("No keys were pressed!")
	} else {
		fmt.Println("Key occurences:")
		for key, amount := range keys {
			fmt.Printf("Key: %q: %d\n", key, amount)
		}
	}
}

func printSessionStats(stats sessionStats) {
	fmt.Println("\n--- Session stats ---")
	fmt.Println("\nKeys pressed during Session:", stats.totalKeys)
	printKeys(stats.keys)
	fmt.Println("\nSession typing speed:", stats.kps, "kps")
}

func capture(keystrokes int, doneChannel chan sessionStats, tickerChannel chan bool, statsChannel chan sessionStats) {
	stats := newSessionStats()

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

			registerKey(event.Rune, &stats)

		case <-tickerChannel:
			statsChannel <- stats
			stats = newSessionStats()

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

func registerKey(key rune, stats *sessionStats) {
	if _, ok := stats.keys[key]; ok {
		stats.keys[key]++
	} else {
		stats.keys[key] = 1
	}

	stats.totalKeys++
}
