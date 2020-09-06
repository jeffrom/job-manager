package resource

type Ack struct {
	ID     string      `json:"id"`
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type Acks struct {
	Acks []*Ack
}
