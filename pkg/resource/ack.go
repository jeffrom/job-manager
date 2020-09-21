package resource

import "encoding/json"

type Ack struct {
	ID     string      `json:"-" db:"id"`
	JobID  string      `json:"job_id" db:"job_id"`
	Status *Status     `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func (ack *Ack) String() string {
	b, _ := json.Marshal(ack)
	return string(b)
}

type Acks struct {
	Acks []*Ack `json:"acks"`
}

func (as *Acks) IDs() []string {
	ids := make([]string, len(as.Acks))
	for i, ack := range as.Acks {
		ids[i] = ack.ID
	}
	return ids
}
