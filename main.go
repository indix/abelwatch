package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ashwanthkumar/wasp-cli/client"

	"github.com/indix/abelwatch/abel"
)

var pidFile string
var slackWebhook string
var waspUrl string
var abelUrl string
var waspNamespace string
var APP_VERSION = "dev-release"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	log.SetOutput(os.Stdout)
	log.Printf("App Version: %s", APP_VERSION)

	defineFlags()
	flag.Parse()
	pid := []byte(fmt.Sprintf("%d\n", os.Getpid()))
	err := ioutil.WriteFile(pidFile, pid, 0644)
	if err != nil {
		fmt.Println("Unable to write pid file. ")
		log.Fatalf("Error - %v\n", err)
	}

	// Rough Idea
	// 1. Manually add a new watch in WASP (just a single watch would do)
	// 2. Peridoically wake up to see if the count of the variable is within acceptable limits
	// 3. If yes, sleep again
	// 4. If not, send a slack alert

	waspClient := &client.WASP{
		Url: waspUrl,
	}
	abelClient := &abel.Abel{
		Url: abelUrl,
	}

	watcher := NewWatchManager(waspClient, abelClient, slackWebhook, waspNamespace)
	watcher.StartAndWait()
}

func defineFlags() {
	flag.StringVar(&pidFile, "pid", "PID", "File to write PID file")
	flag.StringVar(&slackWebhook, "slack-webhook", "", "Slack webhook to post the alert")
	flag.StringVar(&waspUrl, "wasp-url", "", "WASP URL (Eg. http://wasp.domain.tld:9000) without the trailing slash")
	flag.StringVar(&abelUrl, "abel-url", "", "Abel URL (Eg. http://abel.domain.tld:3330) without the trailing slash")
	flag.StringVar(&waspNamespace, "wasp-namespace", "dev.abel.watchers.rules", "Namespace in WASP to get the AbelWatch rules")
}
