// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: api/v1/job.proto

package v1

import (
	proto "github.com/golang/protobuf/proto"
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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_v1_job_proto_rawDescGZIP(), []int{0}
}

var File_api_v1_job_proto protoreflect.FileDescriptor

var file_api_v1_job_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6a, 0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x32, 0xab, 0x01, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x39, 0x0a, 0x07, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x15, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x1a, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x07, 0x44,
	0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x15, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x17, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x11, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x6b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x65,
	0x66, 0x66, 0x72, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x62, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
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

var file_api_v1_job_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_v1_job_proto_goTypes = []interface{}{
	(*Empty)(nil),           // 0: api.v1.Empty
	(*EnqueueParams)(nil),   // 1: api.v1.EnqueueParams
	(*DequeueParams)(nil),   // 2: api.v1.DequeueParams
	(*AckParams)(nil),       // 3: api.v1.AckParams
	(*EnqueueResponse)(nil), // 4: api.v1.EnqueueResponse
	(*DequeueResponse)(nil), // 5: api.v1.DequeueResponse
}
var file_api_v1_job_proto_depIdxs = []int32{
	1, // 0: api.v1.JobService.Enqueue:input_type -> api.v1.EnqueueParams
	2, // 1: api.v1.JobService.Dequeue:input_type -> api.v1.DequeueParams
	3, // 2: api.v1.JobService.Ack:input_type -> api.v1.AckParams
	4, // 3: api.v1.JobService.Enqueue:output_type -> api.v1.EnqueueResponse
	5, // 4: api.v1.JobService.Dequeue:output_type -> api.v1.DequeueResponse
	0, // 5: api.v1.JobService.Ack:output_type -> api.v1.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_v1_job_proto_init() }
func file_api_v1_job_proto_init() {
	if File_api_v1_job_proto != nil {
		return
	}
	file_api_v1_enqueue_proto_init()
	file_api_v1_dequeue_proto_init()
	file_api_v1_ack_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_job_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
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
