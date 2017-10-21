package papertrail_test

import (
	"fmt"
	"io"
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
	res := make(chan []byte, 20)

	// little server that reads 16 bytes at a time
	go func() {
		for {
			buf := make([]byte, 16)
			n, e := server.Read(buf)
			fmt.Printf("read %d bytes\n", n)
			time.Sleep(600 * time.Millisecond)
			if e != nil {
				fmt.Printf("error %s with buf %s (%d bytes)\n", e.Error(), buf, n)
				if e == io.EOF {
					close(res)
					break
				}
				t.Fatal(e)
			}
			res <- buf
		}
	}()

	pt := papertrail.New(&papertrail.Config{
		Host:   "host",
		Port:   8080,
		Writer: client,
	})

	log.SetHandler(pt)
	log.Infof("a")
	log.Infof("b")
	log.Flush()
	fmt.Println("Flushed...")
	client.Close()
	var full []byte
	for b := range res {
		full = append(full, b...)
	}

	fmt.Printf("%s", full)
}
