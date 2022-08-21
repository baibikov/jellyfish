// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: api/proto/consumer.proto

package messages

import (
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

type ConsumerPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *ConsumerPayload) Reset() {
	*x = ConsumerPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_consumer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumerPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumerPayload) ProtoMessage() {}

func (x *ConsumerPayload) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_consumer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumerPayload.ProtoReflect.Descriptor instead.
func (*ConsumerPayload) Descriptor() ([]byte, []int) {
	return file_api_proto_consumer_proto_rawDescGZIP(), []int{0}
}

func (x *ConsumerPayload) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type ConsumerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsEmpty bool   `protobuf:"varint,1,opt,name=isEmpty,proto3" json:"isEmpty,omitempty"`
	Message []byte `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ConsumerResponse) Reset() {
	*x = ConsumerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_consumer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumerResponse) ProtoMessage() {}

func (x *ConsumerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_consumer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumerResponse.ProtoReflect.Descriptor instead.
func (*ConsumerResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_consumer_proto_rawDescGZIP(), []int{1}
}

func (x *ConsumerResponse) GetIsEmpty() bool {
	if x != nil {
		return x.IsEmpty
	}
	return false
}

func (x *ConsumerResponse) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_api_proto_consumer_proto protoreflect.FileDescriptor

var file_api_proto_consumer_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x22, 0x27, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65,
	0x72, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x46,
	0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x19, 0x5a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_consumer_proto_rawDescOnce sync.Once
	file_api_proto_consumer_proto_rawDescData = file_api_proto_consumer_proto_rawDesc
)

func file_api_proto_consumer_proto_rawDescGZIP() []byte {
	file_api_proto_consumer_proto_rawDescOnce.Do(func() {
		file_api_proto_consumer_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_consumer_proto_rawDescData)
	})
	return file_api_proto_consumer_proto_rawDescData
}

var file_api_proto_consumer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_consumer_proto_goTypes = []interface{}{
	(*ConsumerPayload)(nil),  // 0: generated.ConsumerPayload
	(*ConsumerResponse)(nil), // 1: generated.ConsumerResponse
}
var file_api_proto_consumer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_consumer_proto_init() }
func file_api_proto_consumer_proto_init() {
	if File_api_proto_consumer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_consumer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumerPayload); i {
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
		file_api_proto_consumer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumerResponse); i {
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
			RawDescriptor: file_api_proto_consumer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proto_consumer_proto_goTypes,
		DependencyIndexes: file_api_proto_consumer_proto_depIdxs,
		MessageInfos:      file_api_proto_consumer_proto_msgTypes,
	}.Build()
	File_api_proto_consumer_proto = out.File
	file_api_proto_consumer_proto_rawDesc = nil
	file_api_proto_consumer_proto_goTypes = nil
	file_api_proto_consumer_proto_depIdxs = nil
}
