package firehose

import (
	"encoding/json"
	stdlog "log"

	"github.com/apex/log"
	"github.com/apex/log/internal/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

// Config for firehose
type Config struct {
	// AWS Firehose client
	Firehose *firehose.Firehose

	// Name of the firehose stream
	Stream string

	// Advanced: tweak the queue
	// Capacity: number of records before we start dropping
	Capacity int
	// Concurrency: number of records to write simulaneously
	Concurrency int
}

// New firehose logger with the firehose client and the name of the stream
func New(session *session.Session, stream string) *Handler {
	return NewConfig(&Config{
		Firehose:    firehose.New(session),
		Stream:      stream,
		Capacity:    20,
		Concurrency: 10,
	})
}

// NewConfig fn
func NewConfig(config *Config) *Handler {
	queue := queue.New(config.Capacity, config.Concurrency)
	return &Handler{
		c: config,
		q: queue,
	}
}

// Handler for firehose
type Handler struct {
	q *queue.Queue
	c *Config
}

// HandleLog buffers the logs then sends them to firehose
func (h *Handler) HandleLog(e *log.Entry) error {
	return h.q.Push(func() {
		if e := h.send(e); e != nil {
			stdlog.Printf("log/firehose: %s", e)
		}
	})
}

// TODO: consider batching later, though I'm
// not convinced it's necessary since Firehose
// should handle the writes and I don't think
// you're paying more for it.
func (h *Handler) send(e *log.Entry) error {
	streamName := h.c.Stream
	fh := h.c.Firehose

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	record := &firehose.Record{Data: data}
	_, err = fh.PutRecord(&firehose.PutRecordInput{
		DeliveryStreamName: aws.String(streamName),
		Record:             record,
	})

	return err
}

// Flush manually and block until we've drained
func (h *Handler) Flush() {
	h.q.Wait()
}
