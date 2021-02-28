package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/D-0000000000/autoloader/v2/common"
	"github.com/D-0000000000/autoloader/v2/watcher"
)

var Version string = "v2.1.0"

func watch(watcher watcher.Watcher, ch chan common.NotifyPayload) {
	for {
		waitSec := rand.Intn(6) + 6
		watcher.Produce(ch)
		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}

func main() {
	printVersion := flag.Bool("V", false, "Print current version")
	debugMode := flag.Bool("d", false, "Debug with fake server")
	pathPtr := flag.String("c", "config.yaml", "Configuration file")
	flag.Parse()

	if *printVersion {
		fmt.Printf("autoloader %s\n", Version)
		return
	}

	if *debugMode {
		println("Running on debug mode...")
	}

	config, err := LoadConfig(*pathPtr)
	if err != nil {
		log.Fatal(err)
	}

	watchers, err := watcher.ParseWatchers(config.Watchers, *debugMode)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan common.NotifyPayload)

	for _, watcher := range watchers {
		go watch(watcher, ch)
	}

	select {}
}
