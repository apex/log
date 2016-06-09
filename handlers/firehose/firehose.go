package firehose

import (
	// "encoding/json"

	// "github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	// "github.com/aws/aws-sdk-go/aws/session"
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

// New returns a handler for streaming logs into a firehose Kinesis stream.
// Like the Kinesis handler, to configure producer options or pass our own AWS
// Kinesis client use NewConfig instead
// func New(stream, region string) *Handler {

// }

// NewConfig handler for streaming logs into a firehose Kinesis stream.
// random value used as partition key
func NewConfig(streamName string, c client.ConfigProvider, cfgs ...*aws.Config) *Handler {
	fh := firehose.New(c, cfgs...)
	return &Handler{
		producer: fh,
		gen:      fastuuid.MustNewGenerator(),
	}
}
