// Package logger provides an interface for customizable client logging.
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type Logger interface {
	Log(ctx context.Context, e *Event)
}

type Event struct {
	Level   string      `json:"level"`
	Message string      `json:"msg"`
	JobID   string      `json:"job_id,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type NilLogger struct{}

func (l *NilLogger) Log(ctx context.Context, e *Event) {
}

var Nil = &NilLogger{}

type Default struct{}

func (l *Default) Log(ctx context.Context, e *Event) {
	var msg string
	if e != nil {
		msg = e.Message
		e.Message = ""
	}
	b, err := json.Marshal(e)
	if err != nil || e == nil {
		log.Print("error marshalling log event!", e, err)
	} else {
		s := fmt.Sprintf("%s  | %s", msg, string(b))
		if e.Error != nil {
			s += "  | err: " + e.Error.Error()
		}
		log.Print(s)
	}
}
