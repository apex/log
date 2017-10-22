package firehose

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/apex/log"
	"github.com/apex/log/buffer"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/cenkalti/backoff"
)

// Config for firehose
type Config struct {
	// AWS Firehose client
	Client *firehose.Firehose

	// Name of the firehose stream
	Stream string

	// Optionally provide your own buffer
	Buffer *buffer.Buffer

	// Provide your own writer (useful for testing)
	Writer io.Writer
}

// New firehose logger with the firehose client and the name of the stream
func New(client *firehose.Firehose, stream string) *Handler {
	return NewConfig(&Config{
		Client: client,
	})
}

// NewConfig fn
func NewConfig(config *Config) *Handler {
	if config.Writer == nil {
		config.Writer = &writer{
			client: config.Client,
			stream: config.Stream,
		}
	}

	if config.Buffer == nil {
		config.Buffer = buffer.New(config.Writer)
	}

	return &Handler{config}
}

// Handler for firehose
type Handler struct {
	config *Config
}

// HandleLog buffers the logs then sends them to firehose
func (h *Handler) HandleLog(e *log.Entry) error {
	return nil
}

// Flush manually and block until we've drained
func (h *Handler) Flush() {

}

type writer struct {
	client *firehose.Firehose
	stream string
}

func (w *writer) Write(b []byte) (int, error) {
	r := bytes.NewBuffer(b)
	d := json.NewDecoder(r)

	var records []*firehose.Record
	for {
		var raw json.RawMessage
		if e := d.Decode(&raw); e == io.EOF {
			break // done decoding file
		} else if e != nil {
			// probably a permanent issue, don't retry
			return len(b), e
		}
		records = append(records, &firehose.Record{
			Data: raw,
		})
	}

	maxBatchSize := 400
	lrecords := len(records)
	for i := 0; i <= lrecords; i += maxBatchSize {
		end := i + maxBatchSize
		if end > lrecords {
			end = lrecords
		}
		recs := records[i:end]

		// write records with an exponential backoff
		err := backoff.Retry(func() error {
			return w.write(recs)
		}, backoff.NewExponentialBackOff())

		// don't retry for now
		// TODO: we could be smarter here,
		// and retry the remaining records
		if err != nil {
			return len(b), err
		}
	}

	// success!
	return len(b), nil
}

func (w *writer) write(records []*firehose.Record) error {
	_, e := w.client.PutRecordBatch(&firehose.PutRecordBatchInput{
		DeliveryStreamName: aws.String(w.stream),
		Records:            records,
	})
	return e
}
