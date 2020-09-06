package resource

import "github.com/jeffrom/job-manager/pkg/label"

type Ack struct {
	ID     string       `json:"id"`
	Status Status       `json:"status"`
	Data   interface{}  `json:"data,omitempty"`
	Claims label.Claims `json:"claims,omitempty"`
}

type Acks struct {
	Acks []*Ack `json:"acks"`
}
