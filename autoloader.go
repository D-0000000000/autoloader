package main

import (
	// "flag"
	"log"
	"math/rand"
	"time"

	"github.com/D-0000000000/autoloader/common"
	"github.com/D-0000000000/autoloader/watcher"
)


func watch(watcher watcher.Watcher, ch chan common.NotifyPayload) {
	for {
		waitSec := rand.Intn(6) + 12
		watcher.Produce(ch)
		time.Sleep(time.Duration(waitSec) * time.Second)
	}
}

func main() {
	// pathPtr := flag.String("c", "config.yaml", "Configuration file")
	// flag.Parse()

	// notifiers, err := ParseConfig(*pathPtr)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	gfweibo, err := watcher.NewWeiboWatcher(5611537367)
	if err != nil {
		log.Fatal(err)
	}

	arkweibo, err := watcher.NewWeiboWatcher(6279793937)
	if err != nil {
		log.Fatal(err)
	}

	anAnno, err := watcher.NewAkAnnounceWatcher()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan common.NotifyPayload)

	go watch(gfweibo, ch)
	go watch(arkweibo, ch)
	go watch(anAnno, ch)

	// go consume(ch, notifiers)

	select {}
}
