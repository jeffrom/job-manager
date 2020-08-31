// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: api/v1/job.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	_struct "github.com/golang/protobuf/ptypes/struct"
	job "github.com/jeffrom/job-manager/pkg/job"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type EnqueueParamArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Job      string             `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	Args     []*_struct.Value   `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
	Data     map[string]string  `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Retries  int32              `protobuf:"varint,4,opt,name=retries,proto3" json:"retries,omitempty"`
	Duration *duration.Duration `protobuf:"bytes,5,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *EnqueueParamArgs) Reset() {
	*x = EnqueueParamArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueParamArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueParamArgs) ProtoMessage() {}

func (x *EnqueueParamArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueParamArgs.ProtoReflect.Descriptor instead.
func (*EnqueueParamArgs) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{0}
}

func (x *EnqueueParamArgs) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *EnqueueParamArgs) GetArgs() []*_struct.Value {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *EnqueueParamArgs) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *EnqueueParamArgs) GetRetries() int32 {
	if x != nil {
		return x.Retries
	}
	return 0
}

func (x *EnqueueParamArgs) GetDuration() *duration.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

type EnqueueParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobs []*EnqueueParamArgs `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
}

func (x *EnqueueParams) Reset() {
	*x = EnqueueParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnqueueParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueParams) ProtoMessage() {}

func (x *EnqueueParams) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueParams.ProtoReflect.Descriptor instead.
func (*EnqueueParams) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{1}
}

func (x *EnqueueParams) GetJobs() []*EnqueueParamArgs {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type DequeueParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Num       int32    `protobuf:"varint,1,opt,name=num,proto3" json:"num,omitempty"`
	Job       string   `protobuf:"bytes,2,opt,name=job,proto3" json:"job,omitempty"`
	Selectors []string `protobuf:"bytes,3,rep,name=selectors,proto3" json:"selectors,omitempty"`
}

func (x *DequeueParams) Reset() {
	*x = DequeueParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DequeueParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DequeueParams) ProtoMessage() {}

func (x *DequeueParams) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_job_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DequeueParams.ProtoReflect.Descriptor instead.
func (*DequeueParams) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{2}
}

func (x *DequeueParams) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

func (x *DequeueParams) GetJob() string {
	if x != nil {
		return x.Job
	}
	return ""
}

func (x *DequeueParams) GetSelectors() []string {
	if x != nil {
		return x.Selectors
	}
	return nil
}

type AckParamArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string                    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status job.Status                `protobuf:"varint,2,opt,name=status,proto3,enum=job.v1.Status" json:"status,omitempty"`
	Data   map[string]*_struct.Value `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AckParamArgs) Reset() {
	*x = AckParamArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckParamArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckParamArgs) ProtoMessage() {}

func (x *AckParamArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_job_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckParamArgs.ProtoReflect.Descriptor instead.
func (*AckParamArgs) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{3}
}

func (x *AckParamArgs) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AckParamArgs) GetStatus() job.Status {
	if x != nil {
		return x.Status
	}
	return job.Status_STATUS_UNSPECIFIED
}

func (x *AckParamArgs) GetData() map[string]*_struct.Value {
	if x != nil {
		return x.Data
	}
	return nil
}

type AckParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Acks []*AckParamArgs `protobuf:"bytes,1,rep,name=acks,proto3" json:"acks,omitempty"`
}

func (x *AckParams) Reset() {
	*x = AckParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckParams) ProtoMessage() {}

func (x *AckParams) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_job_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckParams.ProtoReflect.Descriptor instead.
func (*AckParams) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{4}
}

func (x *AckParams) GetAcks() []*AckParamArgs {
	if x != nil {
		return x.Acks
	}
	return nil
}

var File_api_v1_job_proto protoreflect.FileDescriptor

var file_api_v1_job_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x6a, 0x6f, 0x62, 0x2f, 0x76, 0x31,
	0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x92, 0x02, 0x0a, 0x10, 0x45,
	0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6a, 0x6f,
	0x62, 0x12, 0x2a, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x12, 0x36, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x41, 0x72, 0x67, 0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12,
	0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x3d, 0x0a, 0x0d, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x12, 0x2c, 0x0a, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73, 0x52, 0x04, 0x6a, 0x6f, 0x62, 0x73, 0x22, 0x51,
	0x0a, 0x0d, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6e, 0x75,
	0x6d, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x6f, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6a, 0x6f, 0x62, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x73, 0x22, 0xcb, 0x01, 0x0a, 0x0c, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x72,
	0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x4f,
	0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x35, 0x0a, 0x09, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x28, 0x0a, 0x04,
	0x61, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73,
	0x52, 0x04, 0x61, 0x63, 0x6b, 0x73, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x65, 0x66, 0x66, 0x72, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x62,
	0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_job_proto_rawDescOnce sync.Once
	file_api_v1_job_proto_rawDescData = file_api_v1_job_proto_rawDesc
)

func file_api_v1_job_proto_rawDescGZIP() []byte {
	file_api_v1_job_proto_rawDescOnce.Do(func() {
		file_api_v1_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_job_proto_rawDescData)
	})
	return file_api_v1_job_proto_rawDescData
}

var file_api_v1_job_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_v1_job_proto_goTypes = []interface{}{
	(*EnqueueParamArgs)(nil),  // 0: api.v1.EnqueueParamArgs
	(*EnqueueParams)(nil),     // 1: api.v1.EnqueueParams
	(*DequeueParams)(nil),     // 2: api.v1.DequeueParams
	(*AckParamArgs)(nil),      // 3: api.v1.AckParamArgs
	(*AckParams)(nil),         // 4: api.v1.AckParams
	nil,                       // 5: api.v1.EnqueueParamArgs.DataEntry
	nil,                       // 6: api.v1.AckParamArgs.DataEntry
	(*_struct.Value)(nil),     // 7: google.protobuf.Value
	(*duration.Duration)(nil), // 8: google.protobuf.Duration
	(job.Status)(0),           // 9: job.v1.Status
}
var file_api_v1_job_proto_depIdxs = []int32{
	7, // 0: api.v1.EnqueueParamArgs.args:type_name -> google.protobuf.Value
	5, // 1: api.v1.EnqueueParamArgs.data:type_name -> api.v1.EnqueueParamArgs.DataEntry
	8, // 2: api.v1.EnqueueParamArgs.duration:type_name -> google.protobuf.Duration
	0, // 3: api.v1.EnqueueParams.jobs:type_name -> api.v1.EnqueueParamArgs
	9, // 4: api.v1.AckParamArgs.status:type_name -> job.v1.Status
	6, // 5: api.v1.AckParamArgs.data:type_name -> api.v1.AckParamArgs.DataEntry
	3, // 6: api.v1.AckParams.acks:type_name -> api.v1.AckParamArgs
	7, // 7: api.v1.AckParamArgs.DataEntry.value:type_name -> google.protobuf.Value
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_api_v1_job_proto_init() }
func file_api_v1_job_proto_init() {
	if File_api_v1_job_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueParamArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_job_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnqueueParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_job_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DequeueParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_job_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckParamArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_v1_job_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_job_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_job_proto_goTypes,
		DependencyIndexes: file_api_v1_job_proto_depIdxs,
		MessageInfos:      file_api_v1_job_proto_msgTypes,
	}.Build()
	File_api_v1_job_proto = out.File
	file_api_v1_job_proto_rawDesc = nil
	file_api_v1_job_proto_goTypes = nil
	file_api_v1_job_proto_depIdxs = nil
}
