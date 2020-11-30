// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: api/proto/processDetails.proto

package processes

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type SortOrderEnums int32

const (
	SortOrderEnums_ASC  SortOrderEnums = 0
	SortOrderEnums_DESC SortOrderEnums = 1
)

// Enum value maps for SortOrderEnums.
var (
	SortOrderEnums_name = map[int32]string{
		0: "ASC",
		1: "DESC",
	}
	SortOrderEnums_value = map[string]int32{
		"ASC":  0,
		"DESC": 1,
	}
)

func (x SortOrderEnums) Enum() *SortOrderEnums {
	p := new(SortOrderEnums)
	*p = x
	return p
}

func (x SortOrderEnums) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SortOrderEnums) Descriptor() protoreflect.EnumDescriptor {
	return file_api_proto_processDetails_proto_enumTypes[0].Descriptor()
}

func (SortOrderEnums) Type() protoreflect.EnumType {
	return &file_api_proto_processDetails_proto_enumTypes[0]
}

func (x SortOrderEnums) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortOrderEnums.Descriptor instead.
func (SortOrderEnums) EnumDescriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{0}
}

type FilterValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *FilterValue) Reset() {
	*x = FilterValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilterValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterValue) ProtoMessage() {}

func (x *FilterValue) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterValue.ProtoReflect.Descriptor instead.
func (*FilterValue) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{0}
}

func (x *FilterValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FilterValue) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type SortOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attribute string         `protobuf:"bytes,1,opt,name=attribute,proto3" json:"attribute,omitempty"`
	Order     SortOrderEnums `protobuf:"varint,2,opt,name=order,proto3,enum=processdetails.SortOrderEnums" json:"order,omitempty"`
}

func (x *SortOrder) Reset() {
	*x = SortOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SortOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SortOrder) ProtoMessage() {}

func (x *SortOrder) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SortOrder.ProtoReflect.Descriptor instead.
func (*SortOrder) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{1}
}

func (x *SortOrder) GetAttribute() string {
	if x != nil {
		return x.Attribute
	}
	return ""
}

func (x *SortOrder) GetOrder() SortOrderEnums {
	if x != nil {
		return x.Order
	}
	return SortOrderEnums_ASC
}

type Process struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcID          string `protobuf:"bytes,1,opt,name=procID,proto3" json:"procID,omitempty"`
	LastSeen        int64  `protobuf:"varint,2,opt,name=lastSeen,proto3" json:"lastSeen,omitempty"`
	FirstSeen       int64  `protobuf:"varint,3,opt,name=firstSeen,proto3" json:"firstSeen,omitempty"`
	Namespace       string `protobuf:"bytes,4,opt,name=namespace,proto3" json:"namespace,omitempty"`
	ApplicationName string `protobuf:"bytes,5,opt,name=applicationName,proto3" json:"applicationName,omitempty"`
}

func (x *Process) Reset() {
	*x = Process{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Process) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Process) ProtoMessage() {}

func (x *Process) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Process.ProtoReflect.Descriptor instead.
func (*Process) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{2}
}

func (x *Process) GetProcID() string {
	if x != nil {
		return x.ProcID
	}
	return ""
}

func (x *Process) GetLastSeen() int64 {
	if x != nil {
		return x.LastSeen
	}
	return 0
}

func (x *Process) GetFirstSeen() int64 {
	if x != nil {
		return x.FirstSeen
	}
	return 0
}

func (x *Process) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Process) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

type ProcessesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Processes []*Process `protobuf:"bytes,1,rep,name=processes,proto3" json:"processes,omitempty"`
}

func (x *ProcessesResponse) Reset() {
	*x = ProcessesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessesResponse) ProtoMessage() {}

func (x *ProcessesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessesResponse.ProtoReflect.Descriptor instead.
func (*ProcessesResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessesResponse) GetProcesses() []*Process {
	if x != nil {
		return x.Processes
	}
	return nil
}

type ProcessesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace       string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	ApplicationName string `protobuf:"bytes,2,opt,name=applicationName,proto3" json:"applicationName,omitempty"`
	Page            uint64 `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	Size            uint64 `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Destination     string `protobuf:"bytes,5,opt,name=destination,proto3" json:"destination,omitempty"`
}

func (x *ProcessesRequest) Reset() {
	*x = ProcessesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessesRequest) ProtoMessage() {}

func (x *ProcessesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessesRequest.ProtoReflect.Descriptor instead.
func (*ProcessesRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{4}
}

func (x *ProcessesRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ProcessesRequest) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *ProcessesRequest) GetPage() uint64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ProcessesRequest) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ProcessesRequest) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

type ProcessesQueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace       string         `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	ApplicationName string         `protobuf:"bytes,2,opt,name=applicationName,proto3" json:"applicationName,omitempty"`
	Page            uint64         `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	Size            uint64         `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Sort            *SortOrder     `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Search          string         `protobuf:"bytes,6,opt,name=search,proto3" json:"search,omitempty"`
	TimestampStart  uint64         `protobuf:"varint,7,opt,name=timestampStart,proto3" json:"timestampStart,omitempty"`
	TimestampEnd    uint64         `protobuf:"varint,8,opt,name=timestampEnd,proto3" json:"timestampEnd,omitempty"`
	Filters         []*FilterValue `protobuf:"bytes,9,rep,name=filters,proto3" json:"filters,omitempty"`
}

func (x *ProcessesQueryRequest) Reset() {
	*x = ProcessesQueryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessesQueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessesQueryRequest) ProtoMessage() {}

func (x *ProcessesQueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessesQueryRequest.ProtoReflect.Descriptor instead.
func (*ProcessesQueryRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{5}
}

func (x *ProcessesQueryRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ProcessesQueryRequest) GetApplicationName() string {
	if x != nil {
		return x.ApplicationName
	}
	return ""
}

func (x *ProcessesQueryRequest) GetPage() uint64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ProcessesQueryRequest) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ProcessesQueryRequest) GetSort() *SortOrder {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ProcessesQueryRequest) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *ProcessesQueryRequest) GetTimestampStart() uint64 {
	if x != nil {
		return x.TimestampStart
	}
	return 0
}

func (x *ProcessesQueryRequest) GetTimestampEnd() uint64 {
	if x != nil {
		return x.TimestampEnd
	}
	return 0
}

func (x *ProcessesQueryRequest) GetFilters() []*FilterValue {
	if x != nil {
		return x.Filters
	}
	return nil
}

type ProcessQueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset     uint64     `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Count      uint64     `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	TotalCount uint64     `protobuf:"varint,3,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	Processes  []*Process `protobuf:"bytes,4,rep,name=processes,proto3" json:"processes,omitempty"`
}

func (x *ProcessQueryResponse) Reset() {
	*x = ProcessQueryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proto_processDetails_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessQueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessQueryResponse) ProtoMessage() {}

func (x *ProcessQueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_processDetails_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessQueryResponse.ProtoReflect.Descriptor instead.
func (*ProcessQueryResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_processDetails_proto_rawDescGZIP(), []int{6}
}

func (x *ProcessQueryResponse) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ProcessQueryResponse) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ProcessQueryResponse) GetTotalCount() uint64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *ProcessQueryResponse) GetProcesses() []*Process {
	if x != nil {
		return x.Processes
	}
	return nil
}

var File_api_proto_processDetails_proto protoreflect.FileDescriptor

var file_api_proto_processDetails_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2c,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x73, 0x77, 0x61, 0x67, 0x67,
	0x65, 0x72, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x0b,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x5f, 0x0a, 0x09, 0x53, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x12, 0x34, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x2e, 0x53, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x75, 0x6d, 0x73, 0x52,
	0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0xa3, 0x01, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x63, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x63, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6c, 0x61,
	0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x53,
	0x65, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74,
	0x53, 0x65, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4a, 0x0a, 0x11,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x35, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0xa4, 0x01, 0x0a, 0x10, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x61,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0xd1, 0x02, 0x0a, 0x15, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x12, 0x26, 0x0a, 0x0e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x45, 0x6e, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x45, 0x6e, 0x64, 0x12, 0x35, 0x0a, 0x07,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x73, 0x22, 0x9c, 0x01, 0x0a, 0x14, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x35, 0x0a, 0x09, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x2a, 0x23, 0x0a, 0x0e, 0x53, 0x6f, 0x72, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x45,
	0x6e, 0x75, 0x6d, 0x73, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x53, 0x43, 0x10, 0x00, 0x12, 0x08, 0x0a,
	0x04, 0x44, 0x45, 0x53, 0x43, 0x10, 0x01, 0x32, 0xc2, 0x02, 0x0a, 0x15, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x88, 0x01, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x12, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x33, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12,
	0x0d, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x92, 0x41,
	0x1b, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x1a, 0x0e, 0x4c, 0x69,
	0x73, 0x74, 0x20, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x9d, 0x01, 0x0a,
	0x0f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3d, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x1c,
	0x0a, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x1a, 0x0f, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x20, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x42, 0xdf, 0x01, 0x5a,
	0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6f, 0x67, 0x69,
	0x71, 0x61, 0x69, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x71, 0x63, 0x74, 0x6c, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x3b, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x92, 0x41, 0xa3, 0x01, 0x12, 0x10, 0x0a, 0x09, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x02, 0x01,
	0x02, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a,
	0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x52, 0x19, 0x0a, 0x03, 0x34, 0x30, 0x30, 0x12, 0x12, 0x0a, 0x10,
	0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x52, 0x17, 0x0a, 0x03, 0x34, 0x30, 0x33, 0x12, 0x10, 0x0a, 0x0e, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x20, 0x45, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x2e, 0x52, 0x33, 0x0a, 0x03, 0x34, 0x30, 0x34,
	0x12, 0x2c, 0x0a, 0x2a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77, 0x68, 0x65,
	0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x20, 0x64,
	0x6f, 0x65, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x2e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_processDetails_proto_rawDescOnce sync.Once
	file_api_proto_processDetails_proto_rawDescData = file_api_proto_processDetails_proto_rawDesc
)

func file_api_proto_processDetails_proto_rawDescGZIP() []byte {
	file_api_proto_processDetails_proto_rawDescOnce.Do(func() {
		file_api_proto_processDetails_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_processDetails_proto_rawDescData)
	})
	return file_api_proto_processDetails_proto_rawDescData
}

var file_api_proto_processDetails_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_proto_processDetails_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_api_proto_processDetails_proto_goTypes = []interface{}{
	(SortOrderEnums)(0),           // 0: processdetails.SortOrderEnums
	(*FilterValue)(nil),           // 1: processdetails.FilterValue
	(*SortOrder)(nil),             // 2: processdetails.SortOrder
	(*Process)(nil),               // 3: processdetails.Process
	(*ProcessesResponse)(nil),     // 4: processdetails.ProcessesResponse
	(*ProcessesRequest)(nil),      // 5: processdetails.ProcessesRequest
	(*ProcessesQueryRequest)(nil), // 6: processdetails.ProcessesQueryRequest
	(*ProcessQueryResponse)(nil),  // 7: processdetails.ProcessQueryResponse
}
var file_api_proto_processDetails_proto_depIdxs = []int32{
	0, // 0: processdetails.SortOrder.order:type_name -> processdetails.SortOrderEnums
	3, // 1: processdetails.ProcessesResponse.processes:type_name -> processdetails.Process
	2, // 2: processdetails.ProcessesQueryRequest.sort:type_name -> processdetails.SortOrder
	1, // 3: processdetails.ProcessesQueryRequest.filters:type_name -> processdetails.FilterValue
	3, // 4: processdetails.ProcessQueryResponse.processes:type_name -> processdetails.Process
	5, // 5: processdetails.ProcessDetailsService.GetProcesses:input_type -> processdetails.ProcessesRequest
	6, // 6: processdetails.ProcessDetailsService.GetProcessQuery:input_type -> processdetails.ProcessesQueryRequest
	4, // 7: processdetails.ProcessDetailsService.GetProcesses:output_type -> processdetails.ProcessesResponse
	7, // 8: processdetails.ProcessDetailsService.GetProcessQuery:output_type -> processdetails.ProcessQueryResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_api_proto_processDetails_proto_init() }
func file_api_proto_processDetails_proto_init() {
	if File_api_proto_processDetails_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proto_processDetails_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilterValue); i {
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
		file_api_proto_processDetails_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SortOrder); i {
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
		file_api_proto_processDetails_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Process); i {
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
		file_api_proto_processDetails_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessesResponse); i {
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
		file_api_proto_processDetails_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessesRequest); i {
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
		file_api_proto_processDetails_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessesQueryRequest); i {
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
		file_api_proto_processDetails_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessQueryResponse); i {
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
			RawDescriptor: file_api_proto_processDetails_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_processDetails_proto_goTypes,
		DependencyIndexes: file_api_proto_processDetails_proto_depIdxs,
		EnumInfos:         file_api_proto_processDetails_proto_enumTypes,
		MessageInfos:      file_api_proto_processDetails_proto_msgTypes,
	}.Build()
	File_api_proto_processDetails_proto = out.File
	file_api_proto_processDetails_proto_rawDesc = nil
	file_api_proto_processDetails_proto_goTypes = nil
	file_api_proto_processDetails_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProcessDetailsServiceClient is the client API for ProcessDetailsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProcessDetailsServiceClient interface {
	GetProcesses(ctx context.Context, in *ProcessesRequest, opts ...grpc.CallOption) (*ProcessesResponse, error)
	GetProcessQuery(ctx context.Context, in *ProcessesQueryRequest, opts ...grpc.CallOption) (*ProcessQueryResponse, error)
}

type processDetailsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessDetailsServiceClient(cc grpc.ClientConnInterface) ProcessDetailsServiceClient {
	return &processDetailsServiceClient{cc}
}

func (c *processDetailsServiceClient) GetProcesses(ctx context.Context, in *ProcessesRequest, opts ...grpc.CallOption) (*ProcessesResponse, error) {
	out := new(ProcessesResponse)
	err := c.cc.Invoke(ctx, "/processdetails.ProcessDetailsService/GetProcesses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processDetailsServiceClient) GetProcessQuery(ctx context.Context, in *ProcessesQueryRequest, opts ...grpc.CallOption) (*ProcessQueryResponse, error) {
	out := new(ProcessQueryResponse)
	err := c.cc.Invoke(ctx, "/processdetails.ProcessDetailsService/GetProcessQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessDetailsServiceServer is the server API for ProcessDetailsService service.
type ProcessDetailsServiceServer interface {
	GetProcesses(context.Context, *ProcessesRequest) (*ProcessesResponse, error)
	GetProcessQuery(context.Context, *ProcessesQueryRequest) (*ProcessQueryResponse, error)
}

// UnimplementedProcessDetailsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProcessDetailsServiceServer struct {
}

func (*UnimplementedProcessDetailsServiceServer) GetProcesses(context.Context, *ProcessesRequest) (*ProcessesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProcesses not implemented")
}
func (*UnimplementedProcessDetailsServiceServer) GetProcessQuery(context.Context, *ProcessesQueryRequest) (*ProcessQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProcessQuery not implemented")
}

func RegisterProcessDetailsServiceServer(s *grpc.Server, srv ProcessDetailsServiceServer) {
	s.RegisterService(&_ProcessDetailsService_serviceDesc, srv)
}

func _ProcessDetailsService_GetProcesses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessDetailsServiceServer).GetProcesses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processdetails.ProcessDetailsService/GetProcesses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessDetailsServiceServer).GetProcesses(ctx, req.(*ProcessesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProcessDetailsService_GetProcessQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessesQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessDetailsServiceServer).GetProcessQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processdetails.ProcessDetailsService/GetProcessQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessDetailsServiceServer).GetProcessQuery(ctx, req.(*ProcessesQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProcessDetailsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "processdetails.ProcessDetailsService",
	HandlerType: (*ProcessDetailsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProcesses",
			Handler:    _ProcessDetailsService_GetProcesses_Handler,
		},
		{
			MethodName: "GetProcessQuery",
			Handler:    _ProcessDetailsService_GetProcessQuery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/processDetails.proto",
}
