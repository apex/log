package main

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/kinesis"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsk "github.com/aws/aws-sdk-go/service/kinesis"
	k "github.com/tj/go-kinesis"
)

func main() {
	log.SetHandler(multi.New(
		text.New(os.Stderr),
		kinesis.NewConfig(k.Config{
			StreamName: "testSteam",
			Client:     awsk.New(session.New(), &aws.Config{Region: aws.String("us-west-2")}),
		}),
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
