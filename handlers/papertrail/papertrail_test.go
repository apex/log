package papertrail_test

import (
	"fmt"
	"net"
	"testing"
)

func newPipe(t *testing.T) (client, server net.Conn) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		defer ln.Close()
		s, err := ln.Accept()
		if err != nil {
			t.Fatal(err)
		}
		server = s
	}()

	client, err = net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	return client, server
}

func TestFlakyPapertrail(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	// mock papertrail server
	go func() {
		for {
			s, e := ln.AcceptTCP()
			if e != nil {
				t.Fatal(e)
			}

			buf := make([]byte, 1)
			if n, e := s.Read(buf); e != nil {
				t.Fatal(e)
			} else {
				fmt.Printf("read %d bytes\n", n)
			}

			s.Close()
		}
	}()

	// make the connection to our little server
	raddr, err := net.ResolveTCPAddr("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	client, e := net.DialTCP("tcp", nil, raddr)
	if e != nil {
		t.Fatal(e)
	}

	if n, e := client.Write([]byte("hi world!")); e != nil {
		fmt.Printf("error writing %s\n", e)
	} else {
		fmt.Printf("wrote %d bytes\n", n)
	}
}
