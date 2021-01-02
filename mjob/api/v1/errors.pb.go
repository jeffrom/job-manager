// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: api/v1/errors.proto

package v1

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
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

type ErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to ErrorResponse:
	//	*ErrorResponse_Error
	//	*ErrorResponse_Invalid
	ErrorResponse isErrorResponse_ErrorResponse `protobuf_oneof:"error_response"`
}

func (x *ErrorResponse) Reset() {
	*x = ErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_errors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorResponse) ProtoMessage() {}

func (x *ErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_errors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorResponse.ProtoReflect.Descriptor instead.
func (*ErrorResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_errors_proto_rawDescGZIP(), []int{0}
}

func (m *ErrorResponse) GetErrorResponse() isErrorResponse_ErrorResponse {
	if m != nil {
		return m.ErrorResponse
	}
	return nil
}

func (x *ErrorResponse) GetError() *GenericError {
	if x, ok := x.GetErrorResponse().(*ErrorResponse_Error); ok {
		return x.Error
	}
	return nil
}

func (x *ErrorResponse) GetInvalid() *ValidationErrorResponse {
	if x, ok := x.GetErrorResponse().(*ErrorResponse_Invalid); ok {
		return x.Invalid
	}
	return nil
}

type isErrorResponse_ErrorResponse interface {
	isErrorResponse_ErrorResponse()
}

type ErrorResponse_Error struct {
	Error *GenericError `protobuf:"bytes,1,opt,name=error,proto3,oneof"`
}

type ErrorResponse_Invalid struct {
	Invalid *ValidationErrorResponse `protobuf:"bytes,2,opt,name=invalid,proto3,oneof"`
}

func (*ErrorResponse_Error) isErrorResponse_ErrorResponse() {}

func (*ErrorResponse_Invalid) isErrorResponse_ErrorResponse() {}

type GenericError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind       string             `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	Message    string             `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Resource   string             `protobuf:"bytes,3,opt,name=resource,proto3" json:"resource,omitempty"`
	ResourceId string             `protobuf:"bytes,4,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	Reason     string             `protobuf:"bytes,5,opt,name=reason,proto3" json:"reason,omitempty"`
	Invalid    []*ValidationError `protobuf:"bytes,6,rep,name=invalid,proto3" json:"invalid,omitempty"`
}

func (x *GenericError) Reset() {
	*x = GenericError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_errors_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenericError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenericError) ProtoMessage() {}

func (x *GenericError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_errors_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenericError.ProtoReflect.Descriptor instead.
func (*GenericError) Descriptor() ([]byte, []int) {
	return file_api_v1_errors_proto_rawDescGZIP(), []int{1}
}

func (x *GenericError) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *GenericError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GenericError) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *GenericError) GetResourceId() string {
	if x != nil {
		return x.ResourceId
	}
	return ""
}

func (x *GenericError) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *GenericError) GetInvalid() []*ValidationError {
	if x != nil {
		return x.Invalid
	}
	return nil
}

type ValidationError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path    string         `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Message string         `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Value   *_struct.Value `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ValidationError) Reset() {
	*x = ValidationError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_errors_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationError) ProtoMessage() {}

func (x *ValidationError) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_errors_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationError.ProtoReflect.Descriptor instead.
func (*ValidationError) Descriptor() ([]byte, []int) {
	return file_api_v1_errors_proto_rawDescGZIP(), []int{2}
}

func (x *ValidationError) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ValidationError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ValidationError) GetValue() *_struct.Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type ValidationErrorArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path    string         `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Message string         `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Value   *_struct.Value `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ValidationErrorArg) Reset() {
	*x = ValidationErrorArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_errors_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationErrorArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationErrorArg) ProtoMessage() {}

func (x *ValidationErrorArg) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_errors_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationErrorArg.ProtoReflect.Descriptor instead.
func (*ValidationErrorArg) Descriptor() ([]byte, []int) {
	return file_api_v1_errors_proto_rawDescGZIP(), []int{3}
}

func (x *ValidationErrorArg) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ValidationErrorArg) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ValidationErrorArg) GetValue() *_struct.Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type ValidationErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errs []*ValidationErrorArg `protobuf:"bytes,1,rep,name=errs,proto3" json:"errs,omitempty"`
}

func (x *ValidationErrorResponse) Reset() {
	*x = ValidationErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_errors_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidationErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidationErrorResponse) ProtoMessage() {}

func (x *ValidationErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_errors_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidationErrorResponse.ProtoReflect.Descriptor instead.
func (*ValidationErrorResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_errors_proto_rawDescGZIP(), []int{4}
}

func (x *ValidationErrorResponse) GetErrs() []*ValidationErrorArg {
	if x != nil {
		return x.Errs
	}
	return nil
}

var File_api_v1_errors_proto protoreflect.FileDescriptor

var file_api_v1_errors_proto_rawDesc = []byte{
	0x0a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c, 0x01, 0x0a, 0x0d,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x3b, 0x0a, 0x07, 0x69,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52,
	0x07, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x42, 0x10, 0x0a, 0x0e, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xc4, 0x01, 0x0a, 0x0c, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6b,
	0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x31,
	0x0a, 0x07, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x07, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x22, 0x6d, 0x0a, 0x0f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x70, 0x0a, 0x12, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x41, 0x72, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x49, 0x0a, 0x17, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a,
	0x04, 0x65, 0x72, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x41, 0x72, 0x67, 0x52, 0x04, 0x65, 0x72, 0x72, 0x73, 0x42, 0x2c, 0x5a,
	0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x65, 0x66, 0x66,
	0x72, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x62, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f,
	0x6d, 0x6a, 0x6f, 0x62, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_v1_errors_proto_rawDescOnce sync.Once
	file_api_v1_errors_proto_rawDescData = file_api_v1_errors_proto_rawDesc
)

func file_api_v1_errors_proto_rawDescGZIP() []byte {
	file_api_v1_errors_proto_rawDescOnce.Do(func() {
		file_api_v1_errors_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_errors_proto_rawDescData)
	})
	return file_api_v1_errors_proto_rawDescData
}

var file_api_v1_errors_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_v1_errors_proto_goTypes = []interface{}{
	(*ErrorResponse)(nil),           // 0: api.v1.ErrorResponse
	(*GenericError)(nil),            // 1: api.v1.GenericError
	(*ValidationError)(nil),         // 2: api.v1.ValidationError
	(*ValidationErrorArg)(nil),      // 3: api.v1.ValidationErrorArg
	(*ValidationErrorResponse)(nil), // 4: api.v1.ValidationErrorResponse
	(*_struct.Value)(nil),           // 5: google.protobuf.Value
}
var file_api_v1_errors_proto_depIdxs = []int32{
	1, // 0: api.v1.ErrorResponse.error:type_name -> api.v1.GenericError
	4, // 1: api.v1.ErrorResponse.invalid:type_name -> api.v1.ValidationErrorResponse
	2, // 2: api.v1.GenericError.invalid:type_name -> api.v1.ValidationError
	5, // 3: api.v1.ValidationError.value:type_name -> google.protobuf.Value
	5, // 4: api.v1.ValidationErrorArg.value:type_name -> google.protobuf.Value
	3, // 5: api.v1.ValidationErrorResponse.errs:type_name -> api.v1.ValidationErrorArg
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_api_v1_errors_proto_init() }
func file_api_v1_errors_proto_init() {
	if File_api_v1_errors_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_errors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorResponse); i {
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
		file_api_v1_errors_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenericError); i {
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
		file_api_v1_errors_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationError); i {
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
		file_api_v1_errors_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationErrorArg); i {
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
		file_api_v1_errors_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidationErrorResponse); i {
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
	file_api_v1_errors_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ErrorResponse_Error)(nil),
		(*ErrorResponse_Invalid)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_errors_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_errors_proto_goTypes,
		DependencyIndexes: file_api_v1_errors_proto_depIdxs,
		MessageInfos:      file_api_v1_errors_proto_msgTypes,
	}.Build()
	File_api_v1_errors_proto = out.File
	file_api_v1_errors_proto_rawDesc = nil
	file_api_v1_errors_proto_goTypes = nil
	file_api_v1_errors_proto_depIdxs = nil
}