package v1

import "github.com/jeffrom/job-manager/pkg/resource"

func AckFromProto(msg *Ack) *resource.Ack {
	return &resource.Ack{
		JobID:  msg.Id,
		Status: jobStatusFromProto(msg.Status),
		Data:   msg.Data.AsInterface(),
		Error:  msg.Error,
	}
}
