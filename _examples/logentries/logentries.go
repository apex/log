package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logentries"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

func main() {
	handler, _ := logentries.New("REPLACE_WITH_YOUR_TOKEN")

	log.SetHandler(multi.New(
		text.New(os.Stderr),
		handler,
	))

	ctx := log.WithFields(log.Fields{
		"file": "something.png",
		"type": "image/png",
		"user": "tobi",
	})

	for range time.Tick(time.Millisecond * 100) {
		ctx.Info("upload")
		ctx.Info("upload complete")
	}
}
