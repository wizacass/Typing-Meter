package main

import (
	"sort"
	"strconv"
)

func newSessionStats() sessionStats {
	return sessionStats{
		totalKeys: 0,
		keys:      make([]keyOccurence, 0),
		kps:       0,
	}
}

func newKeyOccurence(key rune) keyOccurence {
	return keyOccurence{
		key:       key,
		occurence: 1,
	}
}

func atoi(value string) int {
	number, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	}

	return number
}

func registerKey(key rune, stats *sessionStats) {
	if i := findKey(key, stats.keys); i >= 0 {
		stats.keys[i].occurence++
	} else {
		occurence := newKeyOccurence(key)
		stats.keys = append(stats.keys, occurence)
	}

	stats.totalKeys++
}

func mergeStats(from sessionStats, to *sessionStats) {
	for _, occurence := range from.keys {
		mergeKey(occurence, from, to)
	}
	to.totalKeys += from.totalKeys
}

func mergeKey(occurence keyOccurence, from sessionStats, to *sessionStats) {
	if i := findKey(occurence.key, to.keys); i >= 0 {
		to.keys[i].occurence += occurence.occurence
	} else {
		to.keys = append(to.keys, occurence)
	}
}

func findKey(key rune, keys []keyOccurence) int {
	for i, keyInfo := range keys {
		if key == keyInfo.key {
			return i
		}
	}

	return -1
}

func findMostPolularKeys(occurences []keyOccurence, amount int) []keyOccurence {
	sort.Slice(occurences, func(i, j int) bool {
		return occurences[i].occurence > occurences[j].occurence
	})

	if len(occurences) <= amount {
		return occurences
	}

	return occurences[:amount]
}
