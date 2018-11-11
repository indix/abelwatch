package main

import (
	"github.com/ashwanthkumar/wasp-cli/client"

	"github.com/indix/abelwatch/abel"
)

func main() {
	// Rough Idea
	// 1. Manually add a new watch in WASP (just a single watch would do)
	// 2. Peridoically wake up to see if the count of the variable is within acceptable limits
	// 3. If yes, sleep again
	// 4. If not, send a slack alert

	waspClient := &client.WASP{
		Url: "http://wasp.indix.tv:9000",
	}
	abelClient := &abel.Abel{
		URL: "http://abel.prod.indix.tv:3330",
	}

	watcher := NewWatchManager(waspClient, abelClient)
	watcher.StartAndWait()
}
