// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: api/v1/queue.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
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

type SaveQueueParamArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	V               int32              `protobuf:"varint,2,opt,name=v,proto3" json:"v,omitempty"`
	Concurrency     int32              `protobuf:"varint,3,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
	MaxRetries      int32              `protobuf:"varint,4,opt,name=max_retries,json=maxRetries,proto3" json:"max_retries,omitempty"`
	Duration        *duration.Duration `protobuf:"bytes,5,opt,name=duration,proto3" json:"duration,omitempty"`
	CheckinDuration *duration.Duration `protobuf:"bytes,6,opt,name=checkin_duration,json=checkinDuration,proto3" json:"checkin_duration,omitempty"`
	Labels          map[string]string  `protobuf:"bytes,7,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ArgSchema       []byte             `protobuf:"bytes,8,opt,name=arg_schema,json=argSchema,proto3" json:"arg_schema,omitempty"`
	DataSchema      []byte             `protobuf:"bytes,9,opt,name=data_schema,json=dataSchema,proto3" json:"data_schema,omitempty"`
	ResultSchema    []byte             `protobuf:"bytes,10,opt,name=result_schema,json=resultSchema,proto3" json:"result_schema,omitempty"`
	Unique          bool               `protobuf:"varint,11,opt,name=unique,proto3" json:"unique,omitempty"`
}

func (x *SaveQueueParamArgs) Reset() {
	*x = SaveQueueParamArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveQueueParamArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveQueueParamArgs) ProtoMessage() {}

func (x *SaveQueueParamArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveQueueParamArgs.ProtoReflect.Descriptor instead.
func (*SaveQueueParamArgs) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{0}
}

func (x *SaveQueueParamArgs) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SaveQueueParamArgs) GetV() int32 {
	if x != nil {
		return x.V
	}
	return 0
}

func (x *SaveQueueParamArgs) GetConcurrency() int32 {
	if x != nil {
		return x.Concurrency
	}
	return 0
}

func (x *SaveQueueParamArgs) GetMaxRetries() int32 {
	if x != nil {
		return x.MaxRetries
	}
	return 0
}

func (x *SaveQueueParamArgs) GetDuration() *duration.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *SaveQueueParamArgs) GetCheckinDuration() *duration.Duration {
	if x != nil {
		return x.CheckinDuration
	}
	return nil
}

func (x *SaveQueueParamArgs) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *SaveQueueParamArgs) GetArgSchema() []byte {
	if x != nil {
		return x.ArgSchema
	}
	return nil
}

func (x *SaveQueueParamArgs) GetDataSchema() []byte {
	if x != nil {
		return x.DataSchema
	}
	return nil
}

func (x *SaveQueueParamArgs) GetResultSchema() []byte {
	if x != nil {
		return x.ResultSchema
	}
	return nil
}

func (x *SaveQueueParamArgs) GetUnique() bool {
	if x != nil {
		return x.Unique
	}
	return false
}

type SaveQueueParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queues []*SaveQueueParamArgs `protobuf:"bytes,1,rep,name=queues,proto3" json:"queues,omitempty"`
}

func (x *SaveQueueParams) Reset() {
	*x = SaveQueueParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveQueueParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveQueueParams) ProtoMessage() {}

func (x *SaveQueueParams) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveQueueParams.ProtoReflect.Descriptor instead.
func (*SaveQueueParams) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{1}
}

func (x *SaveQueueParams) GetQueues() []*SaveQueueParamArgs {
	if x != nil {
		return x.Queues
	}
	return nil
}

type SaveQueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue *job.Queue `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
}

func (x *SaveQueueResponse) Reset() {
	*x = SaveQueueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveQueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveQueueResponse) ProtoMessage() {}

func (x *SaveQueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveQueueResponse.ProtoReflect.Descriptor instead.
func (*SaveQueueResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{2}
}

func (x *SaveQueueResponse) GetQueue() *job.Queue {
	if x != nil {
		return x.Queue
	}
	return nil
}

type ListQueuesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names     []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
	Selectors []string `protobuf:"bytes,2,rep,name=selectors,proto3" json:"selectors,omitempty"`
}

func (x *ListQueuesRequest) Reset() {
	*x = ListQueuesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListQueuesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueuesRequest) ProtoMessage() {}

func (x *ListQueuesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueuesRequest.ProtoReflect.Descriptor instead.
func (*ListQueuesRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{3}
}

func (x *ListQueuesRequest) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

func (x *ListQueuesRequest) GetSelectors() []string {
	if x != nil {
		return x.Selectors
	}
	return nil
}

type ListQueuesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *job.Queues `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ListQueuesResponse) Reset() {
	*x = ListQueuesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListQueuesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueuesResponse) ProtoMessage() {}

func (x *ListQueuesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueuesResponse.ProtoReflect.Descriptor instead.
func (*ListQueuesResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{4}
}

func (x *ListQueuesResponse) GetData() *job.Queues {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetQueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *job.Queue `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GetQueueResponse) Reset() {
	*x = GetQueueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_queue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQueueResponse) ProtoMessage() {}

func (x *GetQueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_queue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQueueResponse.ProtoReflect.Descriptor instead.
func (*GetQueueResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_queue_proto_rawDescGZIP(), []int{5}
}

func (x *GetQueueResponse) GetData() *job.Queue {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_api_v1_queue_proto protoreflect.FileDescriptor

var file_api_v1_queue_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x6a, 0x6f,
	0x62, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xee, 0x03, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x76,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x76, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b,
	0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x6d,
	0x61, 0x78, 0x5f, 0x72, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x08,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x44, 0x0a, 0x10, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69, 0x6e, 0x5f, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x69,
	0x6e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3e, 0x0a, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x41, 0x72, 0x67, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x72, 0x67,
	0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x61,
	0x72, 0x67, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x64,
	0x61, 0x74, 0x61, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0c, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x45, 0x0a, 0x0f, 0x53, 0x61, 0x76, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x12, 0x32, 0x0a, 0x06, 0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61,
	0x76, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x41, 0x72, 0x67, 0x73,
	0x52, 0x06, 0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x22, 0x38, 0x0a, 0x11, 0x53, 0x61, 0x76, 0x65,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a,
	0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6a,
	0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x22, 0x47, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x22, 0x38, 0x0a, 0x12, 0x4c,
	0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x22, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x35, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x42, 0x2b, 0x5a, 0x29,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x65, 0x66, 0x66, 0x72,
	0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x62, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_v1_queue_proto_rawDescOnce sync.Once
	file_api_v1_queue_proto_rawDescData = file_api_v1_queue_proto_rawDesc
)

func file_api_v1_queue_proto_rawDescGZIP() []byte {
	file_api_v1_queue_proto_rawDescOnce.Do(func() {
		file_api_v1_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_queue_proto_rawDescData)
	})
	return file_api_v1_queue_proto_rawDescData
}

var file_api_v1_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_v1_queue_proto_goTypes = []interface{}{
	(*SaveQueueParamArgs)(nil), // 0: api.v1.SaveQueueParamArgs
	(*SaveQueueParams)(nil),    // 1: api.v1.SaveQueueParams
	(*SaveQueueResponse)(nil),  // 2: api.v1.SaveQueueResponse
	(*ListQueuesRequest)(nil),  // 3: api.v1.ListQueuesRequest
	(*ListQueuesResponse)(nil), // 4: api.v1.ListQueuesResponse
	(*GetQueueResponse)(nil),   // 5: api.v1.GetQueueResponse
	nil,                        // 6: api.v1.SaveQueueParamArgs.LabelsEntry
	(*duration.Duration)(nil),  // 7: google.protobuf.Duration
	(*job.Queue)(nil),          // 8: job.v1.Queue
	(*job.Queues)(nil),         // 9: job.v1.Queues
}
var file_api_v1_queue_proto_depIdxs = []int32{
	7, // 0: api.v1.SaveQueueParamArgs.duration:type_name -> google.protobuf.Duration
	7, // 1: api.v1.SaveQueueParamArgs.checkin_duration:type_name -> google.protobuf.Duration
	6, // 2: api.v1.SaveQueueParamArgs.labels:type_name -> api.v1.SaveQueueParamArgs.LabelsEntry
	0, // 3: api.v1.SaveQueueParams.queues:type_name -> api.v1.SaveQueueParamArgs
	8, // 4: api.v1.SaveQueueResponse.queue:type_name -> job.v1.Queue
	9, // 5: api.v1.ListQueuesResponse.data:type_name -> job.v1.Queues
	8, // 6: api.v1.GetQueueResponse.data:type_name -> job.v1.Queue
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_api_v1_queue_proto_init() }
func file_api_v1_queue_proto_init() {
	if File_api_v1_queue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveQueueParamArgs); i {
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
		file_api_v1_queue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveQueueParams); i {
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
		file_api_v1_queue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveQueueResponse); i {
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
		file_api_v1_queue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListQueuesRequest); i {
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
		file_api_v1_queue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListQueuesResponse); i {
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
		file_api_v1_queue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQueueResponse); i {
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
			RawDescriptor: file_api_v1_queue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_queue_proto_goTypes,
		DependencyIndexes: file_api_v1_queue_proto_depIdxs,
		MessageInfos:      file_api_v1_queue_proto_msgTypes,
	}.Build()
	File_api_v1_queue_proto = out.File
	file_api_v1_queue_proto_rawDesc = nil
	file_api_v1_queue_proto_goTypes = nil
	file_api_v1_queue_proto_depIdxs = nil
}
