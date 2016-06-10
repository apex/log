package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/firehose"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

func main() {
	l := log.Logger{
		Handler: multi.New(
			text.New(os.Stderr),
			firehose.New("testStream", "us-west-2"),
		),
		Level: log.DebugLevel,
	}

	ctx := l.WithFields(log.Fields{
		"file": "cat.png",
		"type": "image/png",
		"user": "Cassius Clay",
	})

	for range time.Tick(time.Millisecond * 100) {
		ctx.Info("upload")
		ctx.Info("upload complete")
	}
	// ctx.Info("Done!")

}
