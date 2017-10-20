package papertrail_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/papertrail"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func TestPapertrail(t *testing.T) {
	server, client := net.Pipe()
	res := make(chan []byte, 1)

	go func() {
		buf, e := ioutil.ReadAll(server)
		if e != nil {
			t.Fatal(e)
		}
		res <- buf
		server.Close()
	}()

	pt := papertrail.New(&papertrail.Config{
		Host: "host",
		Port: 8080,
		Conn: client,
	})

	log.SetHandler(pt)
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.Flush()

	client.Close()
	buf := <-res
	fmt.Printf("%s", buf)
}
