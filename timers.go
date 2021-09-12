package main

import (
	"math"
	"time"
)

func startTimer(seconds int, timerChannel chan bool) {
	duration := time.Duration(seconds * int(time.Second))
	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		timerChannel <- true
	}()
}

func startInterval(seconds int, timerChannel chan bool, tickerChannel chan bool, statsChannel chan sessionStats, doneChannel chan sessionStats) {
	sessionStats := newSessionStats()
	startTime := time.Now()
	duration := time.Duration(seconds * int(time.Second))
	ticker := time.NewTicker(duration)

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
