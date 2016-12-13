// Package LogEntries implements a LogEntries JSON handler.
package logentries

import (
	"encoding/json"

	"github.com/apex/log"
	"github.com/bsphere/le_go"
)

type (
	LogHandler struct {
		logger *le_go.Logger
	}
)

func New(token string) (*LogHandler, error) {
	logger, err := le_go.Connect(token)

	if err != nil {
		return nil, err
	}

	return &LogHandler{
		logger: logger,
	}, nil
}

func (handler *LogHandler) HandleLog(entry *log.Entry) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, writer_error := handler.logger.Write(b)
	return writer_error
}
