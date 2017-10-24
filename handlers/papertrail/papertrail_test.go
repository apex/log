package papertrail_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/papertrail"
)

var host string
var port int

func init() {
	host = os.Getenv("PAPERTRAIL_HOST")
	portenv := os.Getenv("PAPERTRAIL_PORT")
	if host == "" || portenv == "" {
		panic("PAPERTRAIL_HOST & PAPERTRAIL_PORT required to run the papertrail tests")
	}
	p, e := strconv.Atoi(portenv)
	if e != nil {
		panic(e)
	}
	port = p
}

func TestPapertrailFlush(t *testing.T) {
	paper := papertrail.New(&papertrail.Config{
		Hostname: "apex",
		Tag:      "log",
		Host:     host,
		Port:     port,
	})

	l := log.NewEntry(&log.Logger{
		Level:   log.InfoLevel,
		Handler: paper,
	})

	defer l.Flush()
	l.Infof("1")
	l.Infof("2")
	l.Infof("3")
	l.Infof("4")
	l.Infof("5")
}
