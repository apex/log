package firehose

import (
	"encoding/json"

	"github.com/apex/log"
	"github.com/erikreppel/go-firehose"
)

// Handler implementation
type Handler struct {
	appName  string
	Producer *firehose.Producer
}

// HandleLog implements log.Handler
func (h *Handler) HandleLog(e *log.Entry) error {
	j, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = h.Producer.Put(j)
	return err
}

// New returns a handler for streaming logs into a firehose Kinesis stream.
// Like the Kinesis handler, to configure producer options or pass our own AWS
// Kinesis client use NewConfig instead
func New(stream, region string) *Handler {
	return NewConfig(firehose.Config{
		FireHoseName: stream,
		Region:       region,
	})
}

// NewConfig handler for streaming logs into a firehose Kinesis stream.
// random value used as partition key
func NewConfig(config firehose.Config) *Handler {
	return &Handler{
		Producer: firehose.New(config),
	}
}
