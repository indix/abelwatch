package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
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
	SlackChannel string     `json:"slackChannel"`
	SlackWebhook string
	NextCheck    int64
	LastChecked  int64

	stopChannel  chan bool
	RunWaitGroup sync.WaitGroup
	// Shared instance of an AbelClient to talk to Abel server for querying for the given metric
	AbelClient *abel.Abel
}

// NewWatch creates a new watch from the input json
func NewWatch(ID string, rawJSON []byte, abelClient *abel.Abel, slackWebhook string) *Watch {
	now := time.Now().In(time.UTC)

	var watch Watch
	util.JsonDecode(string(rawJSON), &watch)
	watch.RawJSON = rawJSON
	watch.ID = ID
	watch.NextCheck = abel.NextAggregateWindow(watch.Duration, now)
	watch.LastChecked = 1000 * now.Unix()
	watch.AbelClient = abelClient
	watch.SlackWebhook = slackWebhook

	return &watch
}

// StartWatching starts watching the current metric
func (w *Watch) StartWatching() {
	w.stopChannel = make(chan bool)
	w.RunWaitGroup.Add(1)
	if w.Duration > int64(0) {
		go w.watch()
	} else {
		log.Printf("[WARN] Not watching %v since it is aggregating forever.", w.String())
	}
}

// Stop stops the watching the current metric
func (w *Watch) Stop() {
	w.RunWaitGroup.Done()
	w.stopChannel <- true
	close(w.stopChannel)
}

func (w *Watch) watch() {
	log.Printf("[INFO] Starting Watch for %s\n", w.String())
	running := true
	for running {
		now := time.Now().In(time.UTC)
		select {
		case tick := <-time.After(abel.TimeToNextAggregateWindow(w.Duration, now)):
			log.Printf("[INFO] Watching %v at %v\n", w, tick)
			start := abel.PreviousAggregateWindow(w.Duration, now)
			count, datatype, err := w.AbelClient.GetCount(w.Name, w.Tags, start, 0, w.Duration)
			if err != nil || datatype == jsonparser.NotExist {
				log.Printf("[ERROR] Encountered an error while checking the metric")
				log.Printf("%v\n", err)
			}
			if w.Condition.HasBreached(count) {
				attachment := slack.Attachment{}
				attachment.AddField(slack.Field{Title: "Metric", Value: w.Name, Short: true}).
					AddField(slack.Field{Title: "Tags", Value: strings.Join(w.Tags, ","), Short: true}).
					AddField(slack.Field{Title: "Agg. Window", Value: (time.Millisecond * time.Duration(w.Duration)).String(), Short: true}).
					AddField(slack.Field{Title: "Expected", Value: strconv.FormatInt(w.Condition.Value, 10), Short: true}).
					AddField(slack.Field{Title: "Actual", Value: strconv.FormatInt(count, 10), Short: true})
				payload := slack.Payload{
					Text:        fmt.Sprintf("%s has breached the expected condition %d %s %d", w.String(), count, w.Condition.Op, w.Condition.Value),
					Username:    "abelwatch",
					Channel:     w.SlackChannel,
					IconEmoji:   ":helmet_with_white_cross:",
					Attachments: []slack.Attachment{attachment},
				}
				errs := slack.Send(w.SlackWebhook, "", payload)
				if len(errs) > 0 {
					err := combineErrors(errs)
					log.Printf("[ERROR] Exception Found while posting message to Slack: %s\n", err.Error())
					log.Fatalf("%v\n", err)
				}
				log.Printf("Sent a notification for %s\n", w.String())
			} else {
				log.Printf("[TRACE] %s has not breached the expected threshold", w.String())
			}

		case <-w.stopChannel:
			running = false
			log.Printf("[INFO] Stopping to watch (ID=%s) %s [%v] windowed by %d\n", w.ID, w.Name, w.Tags, w.Duration)
		}
	}
}

// String returns the string representation of the metric
func (w *Watch) String() string {
	if len(w.Tags) > 0 {
		return fmt.Sprintf("%s with tags[%s] aggregating at %s", w.Name, strings.Join(w.Tags, ","), (time.Millisecond * time.Duration(w.Duration)).String())
	} else {
		return fmt.Sprintf("%s with no tags aggregating at %s", w.Name, (time.Millisecond * time.Duration(w.Duration)).String())
	}
}

func combineErrors(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		msg := "Error(s):"
		for _, err := range errs {
			msg += " " + err.Error()
		}
		return errors.New(msg)
	} else {
		return nil
	}
}
