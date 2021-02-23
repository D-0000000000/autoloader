package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/D-0000000000/autoloader/common"
	"github.com/gocolly/colly/v2"
)

const iOSClientUA = "arknights/385" +
	" CFNetwork/1220.1" +
	" Darwin/20.3.0'"

type announce struct {
	AnnounceID string `json:"announceId"`
	Title      string `json:"title"`
	IsWebURL   bool   `json:"isWebUrl"`
	WebURL     string `json:"webUrl"`
	Day        int    `json:"day"`
	Month      int    `json:"month"`
	Group      string `json:"group"`
}

type announceMeta struct {
	FocusAnnounceID string     `json:"focusAnnounceId"`
	AnnounceList    []announce `json:"announceList"`
}

type akAnnounceWatcher struct {
	name       string
	focusID    string
	latestAnno announce
	existedID  []string
}

func NewAkAnnounceWatcher() (Watcher, error) {
	watcher := new(akAnnounceWatcher)
	watcher.name = "明日方舟客户端公告"
	err := watcher.setup()
	return watcher, err
}

func (watcher akAnnounceWatcher) fetchAPI() (announceMeta, error) {
	const apiURL = "https://ak-fs.hypergryph.com/announce/IOS/announcement.meta.json?sign=1145141919"
	var err error = nil
	var data announceMeta
	c := colly.NewCollector(
		colly.UserAgent(iOSClientUA),
	)

	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})

	c.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &data)
	})

	c.Visit(apiURL)
	c.Wait()

	return data, err
}

func (watcher *akAnnounceWatcher) setup() error {
	data, err := watcher.fetchAPI()
	if err != nil {
		return err
	}

	watcher.focusID = data.FocusAnnounceID
	watcher.existedID = flushIDList(data.AnnounceList)

	return nil
}

func flushIDList(announceList []announce) []string {
	ret := make([]string, len(announceList))
	for i, anno := range announceList {
		ret[i] = anno.AnnounceID
	}

	return ret
}

func (watcher *akAnnounceWatcher) update() bool {
	data, err := watcher.fetchAPI()
	if err != nil {
		log.Println(err)
		return false
	}

	if watcher.focusID != data.FocusAnnounceID {
		watcher.focusID = data.FocusAnnounceID
		existed := false
		for _, anno := range data.AnnounceList {
			if anno.AnnounceID == data.FocusAnnounceID {
				existed = true
				break
			}
		}
		if existed == false {
			watcher.latestAnno = announce{
				Title:  "出现公告弹窗，可能会有新饼",
				WebURL: "about:blank",
			}
			return true
		}
	}

	for _, anno := range data.AnnounceList {
		newID := anno.AnnounceID
		existed := false
		for _, oldID := range watcher.existedID {
			if newID == oldID {
				existed = true
				break
			}
		}
		if existed == false {
			watcher.existedID = flushIDList(data.AnnounceList)
			if strings.Contains(anno.Title, "制作组通讯") {
				watcher.latestAnno = anno
				return true
			}
		}
	}

	return false
}

func (watcher akAnnounceWatcher) parseContent() common.NotifyPayload {
	anno := watcher.latestAnno

	return common.NotifyPayload{
		Title: watcher.name,
		Body:  anno.Title,
		URL:   anno.WebURL,
	}
}
func (watcher *akAnnounceWatcher) Produce(ch chan common.NotifyPayload) {
	if watcher.update() {
		log.Printf("New post from \"%s\"...\n", watcher.name)
		parseMessage := watcher.parseContent()
		msg := "\"" + parseMessage.Body + "\n" + parseMessage.Title + "\n" + parseMessage.URL + "\n" + "\""
		log.Printf(msg)
		cmd := exec.Command("./qqmessagesender", msg)
		// Specific message sender here
		buf, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(buf))
	} else {
		log.Printf("Waiting for post \"%s\"...\n", watcher.name)
	}
}
