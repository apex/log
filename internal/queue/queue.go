package queue

import (
	"errors"
	"sync"
)

// Queue that will never block, if the queue
// is at capacity it will silently drop new
// requests until we have additional capacity.
type Queue struct {
	wg   *sync.WaitGroup
	jobs chan func()
}

// New queue
func New(capacity int, concurrency int) *Queue {
	jobs := make(chan func(), capacity)
	wg := &sync.WaitGroup{}

	// concurrent workers
	for i := 0; i <= concurrency; i++ {
		go worker(wg, jobs)
	}

	return &Queue{
		jobs: jobs,
		wg:   wg,
	}
}

func worker(wg *sync.WaitGroup, jobs chan func()) {
	for job := range jobs {
		job()
		wg.Done()
	}
}

// Push a function into the queue be called
// when ready. This will never block but it
// may drop functions if we're at capacity.
func (g *Queue) Push(fn func()) error {
	g.wg.Add(1)
	select {
	case g.jobs <- fn:
		// queued
		return nil
	default:
		// dropped
		g.wg.Done()
		return errors.New("queue at capacity, dropped fn")
	}
}

// Wait fn
func (g *Queue) Wait() {
	g.wg.Wait()
}
