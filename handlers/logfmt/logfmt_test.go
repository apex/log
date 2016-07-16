package logfmt_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func Test(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(logfmt.New(&buf))
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.Info("world")
	log.Error("boom")

	expected := `timestamp=1969-12-31T16:00:00-08:00 level=info message=hello user=tj id=123
timestamp=1969-12-31T16:00:00-08:00 level=info message=world
timestamp=1969-12-31T16:00:00-08:00 level=error message=boom
`

	assert.Equal(t, expected, buf.String())
}
