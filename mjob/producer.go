package mjob

import "github.com/jeffrom/job-manager/mjob/client"

type Producer struct {
	client client.Interface
}

func NewProducer() *Producer {
	return &Producer{}
}
