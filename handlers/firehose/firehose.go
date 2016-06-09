package firehose

import (
	"encoding/base64"
	"encoding/json"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/rogpeppe/fastuuid"
)

type Handler struct {
	appName  string
	producer *firehose.Firehose
	gen      *fastuuid.Generator
}
