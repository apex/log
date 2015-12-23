package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/kinesis"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

func main() {
	handlers := multi.New(
		text.New(os.Stderr),
		kinesis.New("logs"),
	)

	log.SetHandler(handlers)

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
