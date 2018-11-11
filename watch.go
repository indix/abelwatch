package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ashwanthkumar/wasp-cli/util"
	"github.com/buger/jsonparser"
	"github.com/indix/abelwatch/abel"
)

// Condition represents the condition that we must apply on the
type Condition struct {
	Op    string `json:"op"`
	Value int64  `json:"value"`
}

// HasBreached checks if a particular condition has been met
func (q *Condition) HasBreached(value int64) bool {
	switch q.Op {
	case "<":
		return value < q.Value
	case ">":
		return value > q.Value
	case "=":
		return value == q.Value
	case ">=":
		return value >= q.Value
	case "<=":
		return value <= q.Value
	default:
		return false
	}
}

// NewCondition creates a new condition
func NewCondition(op string, value int64) *Condition {
	return &Condition{
		Op:    op,
		Value: value,
	}
}

// Watch represents a single watch to monitor and run
type Watch struct {
	RawJSON []byte
	ID      string
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	// Duration can't be zero, it could mean we're trying to run a watch over something that's aggregating forever
	Duration     int64      `json:"duration"`
	Condition    *Condition `json:"condition"`
	NextCheck    int64
	LastChecked  int64
	SlackChannel string `json:"slackChannel"`

	stopChannel  chan bool
	RunWaitGroup sync.WaitGroup
	// Shared instance of an AbelClient to talk to Abel server for querying for the given metric
	AbelClient *abel.Abel
}

// NewWatch creates a new watch from the input json
func NewWatch(ID string, RawJson []byte, abelClient *abel.Abel) *Watch {
	now := time.Now().In(time.UTC)

	var watch Watch
	util.JsonDecode(string(RawJson), &watch)
	watch.RawJSON = RawJson
	watch.ID = ID
	watch.NextCheck = abel.NextAggregateWindow(watch.Duration, now)
	watch.LastChecked = 1000 * now.Unix()
	watch.AbelClient = abelClient

	return &watch
}

// StartWatching starts watching the current metric
func (w *Watch) StartWatching() {
	w.stopChannel = make(chan bool)
	w.RunWaitGroup.Add(1)
	go w.watch()
}

// Stop stops the watching the current metric
func (w *Watch) Stop() {
	w.RunWaitGroup.Done()
	w.stopChannel <- true
	close(w.stopChannel)
}

func (w *Watch) watch() {
	running := true
	for running {
		now := time.Now().In(time.UTC)
		select {
		case tick := <-time.After(abel.TimeToNextAggregateWindow(w.Duration, now)):
			fmt.Printf("Watching %v at %v\n", w, tick)
			start := abel.PreviousAggregateWindow(w.Duration, now)
			count, datatype, err := w.AbelClient.GetCount(w.Name, w.Tags, start, 0, w.Duration)
			if err != nil || datatype == jsonparser.NotExist {
				log.Printf("Encountered an error while checking the metric")
				log.Printf("%v\n", err)
			}
			if w.Condition.HasBreached(count) {
				// TODO: Send a slack message
			} else {
				fmt.Printf("%s has not breached the expected threshold", w.String())
			}

		case <-w.stopChannel:
			running = false
			log.Printf("Stopping to watch (ID=%s) %s [%v] windowed by %d\n", w.ID, w.Name, w.Tags, w.Duration)
		}
	}
}

// String returns the string representation of the metric
func (w *Watch) String() string {
	return fmt.Sprintf("%s [%v] at %s", w.Name, w.Tags, time.Duration(w.Duration).String())
}
