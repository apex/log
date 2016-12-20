package papertrail_test

import (
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/papertrail"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0).UTC()
	}
}

func Test(t *testing.T) {
	log.SetHandler(papertrail.New(&papertrail.Config{
		Host:     "logs4",
		Port:     28705,
		Hostname: "check_processor",
		Tag:      "v1",
	}))

	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.Info("world")
	log.Error("boom")
}
