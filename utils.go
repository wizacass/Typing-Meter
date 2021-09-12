package main

import "strconv"

func newSessionStats() sessionStats {
	return sessionStats{
		totalKeys: 0,
		keys:      make(map[rune]int),
		kps:       0,
	}
}

func atoi(value string) int {
	number, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	}

	return number
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

func getMapKey(keys map[rune]int) rune {
	for key := range keys {
		return key
	}

	return rune(0)
}

func registerKey(key rune, stats *sessionStats) {
	if _, ok := stats.keys[key]; ok {
		stats.keys[key]++
	} else {
		stats.keys[key] = 1
	}

	stats.totalKeys++
}
