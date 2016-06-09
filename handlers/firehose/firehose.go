package firehose

import (
	"encoding/json"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/rogpeppe/fastuuid"
)

// Handler implementation
type Handler struct {
	appName    string
	producer   *firehose.Firehose
	gen        *fastuuid.Generator
	streamName string
}

// HandleLog implements log.Handler
func (h *Handler) HandleLog(e *log.Entry) error {
	j, err := json.Marshal(e)
	if err != nil {
		return err
	}

	i := &firehose.PutRecordInput{
		DeliveryStreamName: aws.String(h.streamName),
		Record: &firehose.Record{
			Data: j,
		},
	}
	_, err = h.producer.PutRecord(i)
	return err
}

// New returns a handler for streaming logs into a firehose Kinesis stream.
// Like the Kinesis handler, to configure producer options or pass our own AWS
// Kinesis client use NewConfig instead
func New(stream, region string) *Handler {
	return NewConfig(stream, session.New(), &aws.Config{Region: aws.String(region)})
}

// NewConfig handler for streaming logs into a firehose Kinesis stream.
// random value used as partition key
func NewConfig(streamName string, c client.ConfigProvider, cfgs ...*aws.Config) *Handler {
	fh := firehose.New(c, cfgs...)
	return &Handler{
		producer: fh,
		gen:      fastuuid.MustNewGenerator(),
	}
}
