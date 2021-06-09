package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/D-0000000000/autoloader/v2/common"
	"github.com/D-0000000000/autoloader/v2/watcher"
)

var Version string = "v2.3.4"

func consume(ch chan common.NotifyPayload) {
	msg := <-ch
	msgfile, err := os.Create("msgTitle.txt")
	if err != nil {
		panic(err)
	}
	msgfile.WriteString(msg.Title + "\n")
	msgfile.Close()
	msgfile, err = os.Create("msgBody.txt")
	if err != nil {
		panic(err)
	}
	msgfile.WriteString(msg.Body + "\n")
	msgfile.Close()
	msgfile, err = os.Create("msgURL.txt")
	if err != nil {
		panic(err)
	}
	msgfile.WriteString(msg.URL + "\n")
	msgfile.Close()
	msgfile, err = os.Create("msgPicURL.txt")
	if err != nil {
		panic(err)
	}
	msgfile.WriteString(msg.PicURL + "\n")
	msgfile.Close()
	cmd := exec.Command("./qqmessagesender")
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(buf))
}

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
	pathPtr := flag.String("c", ".", "Configuration and data directory")
	flag.Parse()

	if *printVersion {
		fmt.Printf("dr-feeder %s\n", Version)
		return
	}

	if *debugMode {
		println("Running on debug mode...")
	}

	config, err := LoadConfig(path.Join(*pathPtr, "config.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	watchers, err := watcher.ParseWatchers(config.Watchers, *pathPtr, *debugMode)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan common.NotifyPayload)

	for _, watcher := range watchers {
		go watch(watcher, ch)
	}

	go consume(ch)

	select {}
}
