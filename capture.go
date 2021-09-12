package main

import "github.com/eiannone/keyboard"

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
