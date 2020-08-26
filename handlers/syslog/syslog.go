// Package syslog implements output to local or remote hosts via the syslog protocol.
package syslog

import (
	"fmt"
	"log/syslog"

	"github.com/apex/log"
	"github.com/go-logfmt/logfmt"
)

// Handler implementation.
type Handler struct {
	log *syslog.Writer
	// map apex levels to syslog levels.
	levelsMap map[log.Level](func(m string) error)
}

// New syslog handler.
//  addr: network address of syslog host. (default: localhost)
//  facility: syslog facility, as defined in "log/syslog". (default: LOG_USER)
//  tag: message identifier, usually the application name. (default: os.Args[0])
func New(addr string, facility syslog.Priority, tag string) (h *Handler, err error) {
	h = new(Handler)
	if facility == 0 {
		facility = syslog.LOG_USER
	}

	if addr == "" {
		if h.log, err = syslog.New(facility, tag); err != nil {
			// sometimes local socket connections don't work, so try localhost before
			// giving up.
			if h.log, err = syslog.Dial("tcp", "localhost:514", facility, tag); err != nil {
				return nil, err
			}
		}
	} else {
		if h.log, err = syslog.Dial("tcp", addr, facility, tag); err != nil {
			return nil, err
		}
	}

	// create apex -> syslog level mapping dynamically with h.log methods.
	// wish "log/syslog" had a Writer.Log(p Priority, m string) method.
	h.levelsMap = map[log.Level](func(m string) error){
		log.DebugLevel: h.log.Debug,
		log.InfoLevel:  h.log.Info,
		log.WarnLevel:  h.log.Warning,
		log.ErrorLevel: h.log.Err,
		log.FatalLevel: h.log.Emerg,
	}

	return h, nil
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	keyvals := make([]interface{}, 0, len(e.Fields.Names()) * 2)
	for _, name := range e.Fields.Names() {
		if name == "source" {
			continue
		}
		keyvals = append(keyvals, name)
		keyvals = append(keyvals, e.Fields.Get(name))
	}
	fmt.Println("keyvals:", keyvals)

	out, err := logfmt.MarshalKeyvals(keyvals...)
	fmt.Println(string(out), err)
	if err != nil {
		return err
	}
	return h.levelsMap[e.Level](string(out))
}
