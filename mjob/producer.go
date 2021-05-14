package mjob

import "github.com/jeffrom/job-manager/mjob/client"

type Producer struct {
	c client.Interface
}

func NewProducer(c client.Interface) *Producer {
	return &Producer{c: c}
}
