package json_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
)

func init() {
	log.Now = func() time.Time {
		return time.Unix(0, 0)
	}
}

func Test(t *testing.T) {
	var buf bytes.Buffer

	log.SetHandler(json.New(&buf))
	log.WithField("user", "tj").WithField("id", "123").Info("hello")
	log.Info("world")
	log.Error("boom")

	expected := `{"fields":{"id":"123","user":"tj"},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"hello"}
{"fields":{},"level":"info","timestamp":"1969-12-31T16:00:00-08:00","message":"world"}
{"fields":{},"level":"error","timestamp":"1969-12-31T16:00:00-08:00","message":"boom"}
`

	assert.Equal(t, expected, buf.String())
}
