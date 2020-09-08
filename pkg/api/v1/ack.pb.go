// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: api/v1/ack.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	v1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
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

type AckJobsRequestArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status v1.Status      `protobuf:"varint,2,opt,name=status,proto3,enum=job.v1.Status" json:"status,omitempty"`
	Data   *_struct.Value `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Claims []string       `protobuf:"bytes,4,rep,name=claims,proto3" json:"claims,omitempty"`
}

func (x *AckJobsRequestArgs) Reset() {
	*x = AckJobsRequestArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_ack_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckJobsRequestArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckJobsRequestArgs) ProtoMessage() {}

func (x *AckJobsRequestArgs) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_ack_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckJobsRequestArgs.ProtoReflect.Descriptor instead.
func (*AckJobsRequestArgs) Descriptor() ([]byte, []int) {
	return file_api_v1_ack_proto_rawDescGZIP(), []int{0}
}

func (x *AckJobsRequestArgs) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AckJobsRequestArgs) GetStatus() v1.Status {
	if x != nil {
		return x.Status
	}
	return v1.Status_STATUS_UNSPECIFIED
}

func (x *AckJobsRequestArgs) GetData() *_struct.Value {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AckJobsRequestArgs) GetClaims() []string {
	if x != nil {
		return x.Claims
	}
	return nil
}

type AckJobsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Acks []*AckJobsRequestArgs `protobuf:"bytes,1,rep,name=acks,proto3" json:"acks,omitempty"`
}

func (x *AckJobsRequest) Reset() {
	*x = AckJobsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_ack_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckJobsRequest) ProtoMessage() {}

func (x *AckJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_ack_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckJobsRequest.ProtoReflect.Descriptor instead.
func (*AckJobsRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_ack_proto_rawDescGZIP(), []int{1}
}

func (x *AckJobsRequest) GetAcks() []*AckJobsRequestArgs {
	if x != nil {
		return x.Acks
	}
	return nil
}

type AckJobsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AckJobsResponse) Reset() {
	*x = AckJobsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_ack_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckJobsResponse) ProtoMessage() {}

func (x *AckJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_ack_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckJobsResponse.ProtoReflect.Descriptor instead.
func (*AckJobsResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_ack_proto_rawDescGZIP(), []int{2}
}

var File_api_v1_ack_proto protoreflect.FileDescriptor

var file_api_v1_ack_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x6a, 0x6f, 0x62, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01,
	0x0a, 0x12, 0x41, 0x63, 0x6b, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x41, 0x72, 0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x6a, 0x6f, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2a, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6c, 0x61, 0x69,
	0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x22, 0x40, 0x0a, 0x0e, 0x41, 0x63, 0x6b, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2e, 0x0a, 0x04, 0x61, 0x63, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x6b, 0x4a, 0x6f, 0x62,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x41, 0x72, 0x67, 0x73, 0x52, 0x04, 0x61, 0x63,
	0x6b, 0x73, 0x22, 0x11, 0x0a, 0x0f, 0x41, 0x63, 0x6b, 0x4a, 0x6f, 0x62, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x65, 0x66, 0x66, 0x72, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x62, 0x2d,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_ack_proto_rawDescOnce sync.Once
	file_api_v1_ack_proto_rawDescData = file_api_v1_ack_proto_rawDesc
)

func file_api_v1_ack_proto_rawDescGZIP() []byte {
	file_api_v1_ack_proto_rawDescOnce.Do(func() {
		file_api_v1_ack_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_ack_proto_rawDescData)
	})
	return file_api_v1_ack_proto_rawDescData
}

var file_api_v1_ack_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_v1_ack_proto_goTypes = []interface{}{
	(*AckJobsRequestArgs)(nil), // 0: api.v1.AckJobsRequestArgs
	(*AckJobsRequest)(nil),     // 1: api.v1.AckJobsRequest
	(*AckJobsResponse)(nil),    // 2: api.v1.AckJobsResponse
	(v1.Status)(0),             // 3: job.v1.Status
	(*_struct.Value)(nil),      // 4: google.protobuf.Value
}
var file_api_v1_ack_proto_depIdxs = []int32{
	3, // 0: api.v1.AckJobsRequestArgs.status:type_name -> job.v1.Status
	4, // 1: api.v1.AckJobsRequestArgs.data:type_name -> google.protobuf.Value
	0, // 2: api.v1.AckJobsRequest.acks:type_name -> api.v1.AckJobsRequestArgs
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_v1_ack_proto_init() }
func file_api_v1_ack_proto_init() {
	if File_api_v1_ack_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_ack_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckJobsRequestArgs); i {
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
		file_api_v1_ack_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckJobsRequest); i {
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
		file_api_v1_ack_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckJobsResponse); i {
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
			RawDescriptor: file_api_v1_ack_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_ack_proto_goTypes,
		DependencyIndexes: file_api_v1_ack_proto_depIdxs,
		MessageInfos:      file_api_v1_ack_proto_msgTypes,
	}.Build()
	File_api_v1_ack_proto = out.File
	file_api_v1_ack_proto_rawDesc = nil
	file_api_v1_ack_proto_goTypes = nil
	file_api_v1_ack_proto_depIdxs = nil
}
