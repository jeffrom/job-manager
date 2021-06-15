package v1

import "github.com/jeffrom/job-manager/mjob/resource"

func AckFromProto(msg *Ack) *resource.Ack {
	return &resource.Ack{
		JobID:  msg.Id,
		Status: JobStatusFromProto(msg.Status),
		Data:   msg.Data,
		Error:  msg.Error,
	}
}
