// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.1
// source: api/api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Number struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	X             float64                `protobuf:"fixed64,10,opt,name=x,proto3" json:"x,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Number) Reset() {
	*x = Number{}
	mi := &file_api_api_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Number) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Number) ProtoMessage() {}

func (x *Number) ProtoReflect() protoreflect.Message {
	mi := &file_api_api_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Number.ProtoReflect.Descriptor instead.
func (*Number) Descriptor() ([]byte, []int) {
	return file_api_api_proto_rawDescGZIP(), []int{0}
}

func (x *Number) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

var File_api_api_proto protoreflect.FileDescriptor

const file_api_api_proto_rawDesc = "" +
	"\n" +
	"\rapi/api.proto\x12\x06api.v1\"\x16\n" +
	"\x06Number\x12\f\n" +
	"\x01x\x18\n" +
	" \x01(\x01R\x01x2\xc6\x01\n" +
	"\vCalsService\x12*\n" +
	"\x06Square\x12\x0e.api.v1.Number\x1a\x0e.api.v1.Number\"\x00\x12)\n" +
	"\x03Sum\x12\x0e.api.v1.Number\x1a\x0e.api.v1.Number\"\x00(\x01\x12,\n" +
	"\x06Repeat\x12\x0e.api.v1.Number\x1a\x0e.api.v1.Number\"\x000\x01\x122\n" +
	"\n" +
	"PipeSquare\x12\x0e.api.v1.Number\x1a\x0e.api.v1.Number\"\x00(\x010\x01B\aZ\x05./apib\x06proto3"

var (
	file_api_api_proto_rawDescOnce sync.Once
	file_api_api_proto_rawDescData []byte
)

func file_api_api_proto_rawDescGZIP() []byte {
	file_api_api_proto_rawDescOnce.Do(func() {
		file_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_api_proto_rawDesc), len(file_api_api_proto_rawDesc)))
	})
	return file_api_api_proto_rawDescData
}

var file_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_api_proto_goTypes = []any{
	(*Number)(nil), // 0: api.v1.Number
}
var file_api_api_proto_depIdxs = []int32{
	0, // 0: api.v1.CalsService.Square:input_type -> api.v1.Number
	0, // 1: api.v1.CalsService.Sum:input_type -> api.v1.Number
	0, // 2: api.v1.CalsService.Repeat:input_type -> api.v1.Number
	0, // 3: api.v1.CalsService.PipeSquare:input_type -> api.v1.Number
	0, // 4: api.v1.CalsService.Square:output_type -> api.v1.Number
	0, // 5: api.v1.CalsService.Sum:output_type -> api.v1.Number
	0, // 6: api.v1.CalsService.Repeat:output_type -> api.v1.Number
	0, // 7: api.v1.CalsService.PipeSquare:output_type -> api.v1.Number
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_api_proto_init() }
func file_api_api_proto_init() {
	if File_api_api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_api_proto_rawDesc), len(file_api_api_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_api_proto_goTypes,
		DependencyIndexes: file_api_api_proto_depIdxs,
		MessageInfos:      file_api_api_proto_msgTypes,
	}.Build()
	File_api_api_proto = out.File
	file_api_api_proto_goTypes = nil
	file_api_api_proto_depIdxs = nil
}
