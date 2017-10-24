package firehose_test

import (
	"os"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/firehose"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var accessKey, secretKey, region string

func init() {
	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	region = os.Getenv("AWS_REGION")
	if accessKey == "" || secretKey == "" || region == "" {
		panic("AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY & AWS_REGION needed to run the firehose tests")
	}
}

func TestFirehose(t *testing.T) {
	sess, e := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
	})
	if e != nil {
		t.Fatal(e)
	}

	fh := firehose.New(sess, "testhose")

	l := log.NewEntry(&log.Logger{
		Level:   log.InfoLevel,
		Handler: fh,
	})

	defer l.Flush()
	l.Infof("1")
	l.Infof("2")
	l.Infof("3")
	l.Infof("4")
	l.WithField("five", 5).Infof("5")
}
