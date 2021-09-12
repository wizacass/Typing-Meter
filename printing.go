package main

import "fmt"

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
		popularKeys := findMostPolularKeys(keys, 3)

		fmt.Println("Most popular key occurences:")
		for key, amount := range popularKeys {
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
