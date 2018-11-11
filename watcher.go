package main

import (
	"log"
	"sync"
	"time"

	waspclient "github.com/ashwanthkumar/wasp-cli/client"
	"github.com/buger/jsonparser"
	"github.com/indix/abelwatch/abel"
)

// WatchManager is responsible for managing the life-cycle of all the active Watchs
type WatchManager struct {
	WASP         *waspclient.WASP
	Abel         *abel.Abel
	RunWaitGroup sync.WaitGroup
	stopChannel  chan bool
	idToWatch    map[string]*Watch
}

// StartAndWait starts and waits indefinitely for the WatchManager to complete
func (w *WatchManager) StartAndWait() {
	w.stopChannel = make(chan bool)
	w.RunWaitGroup.Add(1)

	go w.run()
	w.RunWaitGroup.Wait()
}

// Run starts the WatchManager
func (w *WatchManager) run() {
	w.pollUpdatesInWasp() // initial pull from WASP
	running := true
	for running {
		select {
		case <-time.After(1 * time.Minute):
			w.pollUpdatesInWasp()

		case <-w.stopChannel:
			running = false
		}

		time.Sleep(1 * time.Second)
	}
}

func (w *WatchManager) pollUpdatesInWasp() {
	log.Printf("[INFO] Polling WASP for new updates")
	config, err := w.WASP.Get("dev.abel.watchers.rules")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	jsonparser.ObjectEach([]byte(config), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		id := string(key)
		_, present := w.idToWatch[id]
		if !present {
			watch := NewWatch(id, value, w.Abel)
			w.idToWatch[id] = watch
			watch.StartWatching()
		}
		return nil
	})
}

// Stop stops the WatchManager
func (w *WatchManager) Stop() {
	log.Println("Stopping Watcher...")
	close(w.stopChannel)
	w.RunWaitGroup.Done()
}

// NewWatchManager creates a new instance of WatchManager
func NewWatchManager(waspClient *waspclient.WASP, abelClient *abel.Abel) *WatchManager {
	return &WatchManager{
		WASP:      waspClient,
		Abel:      abelClient,
		idToWatch: make(map[string]*Watch),
	}
}
