package watcher

import "github.com/D-0000000000/autoloader/common"

type Watcher interface {
	Produce(chan common.NotifyPayload)
}
