package queue_test

import (
	"testing"

	"github.com/apex/log/internal/queue"
)

func TestNonBlockingWrites(t *testing.T) {
	q := queue.New(1, 1)
	block := make(chan bool)

	e := q.Push(func() {
		<-block
	})
	if e != nil {
		t.Fatal(e)
	}

	e = q.Push(func() {
		<-block
	})
	if e == nil {
		t.Fatal("should have error")
	}

	e = q.Push(func() {
		<-block
	})
	if e == nil {
		t.Fatal("should have error")
	}

	e = q.Push(func() {
		<-block
	})
	if e == nil {
		t.Fatal("should have error")
	}

	// unblock
	block <- false

	// should work again
	e = q.Push(func() {
		<-block
	})
	if e != nil {
		t.Fatal(e)
	}

	block <- false
	q.Wait()
}
