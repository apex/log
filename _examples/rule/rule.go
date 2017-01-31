package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/rule"
)

func main() {
	log.SetHandler(rule.New(cli.Default, func(entry *log.Entry) bool {
		if entry.Fields.Get("internal") != nil {
			return false
		}

		return true
	}))

	log.SetLevel(log.DebugLevel)

	log.WithField("internal", "yup").Error("this won't go to log")

	log.WithField("alice", "bob").Info("but this will")
}
