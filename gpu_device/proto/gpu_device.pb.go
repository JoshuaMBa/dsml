// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.20.3
// source: gpu_device/proto/gpu_device.proto

package proto

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

type StreamStatus int32

const (
	StreamStatus_IN_PROGRESS StreamStatus = 0
	StreamStatus_SUCCESS     StreamStatus = 1
	StreamStatus_FAILED      StreamStatus = 2
)

// Enum value maps for StreamStatus.
var (
	StreamStatus_name = map[int32]string{
		0: "IN_PROGRESS",
		1: "SUCCESS",
		2: "FAILED",
	}
	StreamStatus_value = map[string]int32{
		"IN_PROGRESS": 0,
		"SUCCESS":     1,
		"FAILED":      2,
	}
)

func (x StreamStatus) Enum() *StreamStatus {
	p := new(StreamStatus)
	*p = x
	return p
}

func (x StreamStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StreamStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_gpu_device_proto_gpu_device_proto_enumTypes[0].Descriptor()
}

func (StreamStatus) Type() protoreflect.EnumType {
	return &file_gpu_device_proto_gpu_device_proto_enumTypes[0]
}

func (x StreamStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StreamStatus.Descriptor instead.
func (StreamStatus) EnumDescriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{0}
}

type DeviceId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DeviceId) Reset() {
	*x = DeviceId{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeviceId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceId) ProtoMessage() {}

func (x *DeviceId) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceId.ProtoReflect.Descriptor instead.
func (*DeviceId) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceId) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Rank struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Rank) Reset() {
	*x = Rank{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Rank) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rank) ProtoMessage() {}

func (x *Rank) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rank.ProtoReflect.Descriptor instead.
func (*Rank) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{1}
}

func (x *Rank) GetValue() uint32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type MemAddr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MemAddr) Reset() {
	*x = MemAddr{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MemAddr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemAddr) ProtoMessage() {}

func (x *MemAddr) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemAddr.ProtoReflect.Descriptor instead.
func (*MemAddr) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{2}
}

func (x *MemAddr) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type StreamId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value uint64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StreamId) Reset() {
	*x = StreamId{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamId) ProtoMessage() {}

func (x *StreamId) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamId.ProtoReflect.Descriptor instead.
func (*StreamId) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{3}
}

func (x *StreamId) GetValue() uint64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type DeviceMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceId   *DeviceId `protobuf:"bytes,1,opt,name=deviceId,proto3" json:"deviceId,omitempty"`
	MinMemAddr *MemAddr  `protobuf:"bytes,2,opt,name=minMemAddr,proto3" json:"minMemAddr,omitempty"`
	MaxMemAddr *MemAddr  `protobuf:"bytes,3,opt,name=maxMemAddr,proto3" json:"maxMemAddr,omitempty"`
}

func (x *DeviceMetadata) Reset() {
	*x = DeviceMetadata{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeviceMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceMetadata) ProtoMessage() {}

func (x *DeviceMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceMetadata.ProtoReflect.Descriptor instead.
func (*DeviceMetadata) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{4}
}

func (x *DeviceMetadata) GetDeviceId() *DeviceId {
	if x != nil {
		return x.DeviceId
	}
	return nil
}

func (x *DeviceMetadata) GetMinMemAddr() *MemAddr {
	if x != nil {
		return x.MinMemAddr
	}
	return nil
}

func (x *DeviceMetadata) GetMaxMemAddr() *MemAddr {
	if x != nil {
		return x.MaxMemAddr
	}
	return nil
}

type GetDeviceMetadataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetDeviceMetadataRequest) Reset() {
	*x = GetDeviceMetadataRequest{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDeviceMetadataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceMetadataRequest) ProtoMessage() {}

func (x *GetDeviceMetadataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceMetadataRequest.ProtoReflect.Descriptor instead.
func (*GetDeviceMetadataRequest) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{5}
}

type GetDeviceMetadataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata      *DeviceMetadata   `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	RankToAddress map[uint32]string `protobuf:"bytes,2,rep,name=rank_to_address,json=rankToAddress,proto3" json:"rank_to_address,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // rank -> "IP:Port"
}

func (x *GetDeviceMetadataResponse) Reset() {
	*x = GetDeviceMetadataResponse{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDeviceMetadataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDeviceMetadataResponse) ProtoMessage() {}

func (x *GetDeviceMetadataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDeviceMetadataResponse.ProtoReflect.Descriptor instead.
func (*GetDeviceMetadataResponse) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{6}
}

func (x *GetDeviceMetadataResponse) GetMetadata() *DeviceMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *GetDeviceMetadataResponse) GetRankToAddress() map[uint32]string {
	if x != nil {
		return x.RankToAddress
	}
	return nil
}

type BeginSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SendBuffAddr *MemAddr `protobuf:"bytes,1,opt,name=sendBuffAddr,proto3" json:"sendBuffAddr,omitempty"`
	NumBytes     uint64   `protobuf:"varint,2,opt,name=numBytes,proto3" json:"numBytes,omitempty"`
	DstRank      *Rank    `protobuf:"bytes,3,opt,name=dstRank,proto3" json:"dstRank,omitempty"`
}

func (x *BeginSendRequest) Reset() {
	*x = BeginSendRequest{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BeginSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BeginSendRequest) ProtoMessage() {}

func (x *BeginSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BeginSendRequest.ProtoReflect.Descriptor instead.
func (*BeginSendRequest) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{7}
}

func (x *BeginSendRequest) GetSendBuffAddr() *MemAddr {
	if x != nil {
		return x.SendBuffAddr
	}
	return nil
}

func (x *BeginSendRequest) GetNumBytes() uint64 {
	if x != nil {
		return x.NumBytes
	}
	return 0
}

func (x *BeginSendRequest) GetDstRank() *Rank {
	if x != nil {
		return x.DstRank
	}
	return nil
}

type BeginSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Initiated bool      `protobuf:"varint,1,opt,name=initiated,proto3" json:"initiated,omitempty"`
	StreamId  *StreamId `protobuf:"bytes,2,opt,name=streamId,proto3" json:"streamId,omitempty"`
}

func (x *BeginSendResponse) Reset() {
	*x = BeginSendResponse{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BeginSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BeginSendResponse) ProtoMessage() {}

func (x *BeginSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BeginSendResponse.ProtoReflect.Descriptor instead.
func (*BeginSendResponse) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{8}
}

func (x *BeginSendResponse) GetInitiated() bool {
	if x != nil {
		return x.Initiated
	}
	return false
}

func (x *BeginSendResponse) GetStreamId() *StreamId {
	if x != nil {
		return x.StreamId
	}
	return nil
}

type BeginReceiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamId     *StreamId `protobuf:"bytes,1,opt,name=streamId,proto3" json:"streamId,omitempty"`
	RecvBuffAddr *MemAddr  `protobuf:"bytes,2,opt,name=recvBuffAddr,proto3" json:"recvBuffAddr,omitempty"`
	NumBytes     uint64    `protobuf:"varint,3,opt,name=numBytes,proto3" json:"numBytes,omitempty"`
	SrcRank      *Rank     `protobuf:"bytes,4,opt,name=srcRank,proto3" json:"srcRank,omitempty"`
}

func (x *BeginReceiveRequest) Reset() {
	*x = BeginReceiveRequest{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BeginReceiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BeginReceiveRequest) ProtoMessage() {}

func (x *BeginReceiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BeginReceiveRequest.ProtoReflect.Descriptor instead.
func (*BeginReceiveRequest) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{9}
}

func (x *BeginReceiveRequest) GetStreamId() *StreamId {
	if x != nil {
		return x.StreamId
	}
	return nil
}

func (x *BeginReceiveRequest) GetRecvBuffAddr() *MemAddr {
	if x != nil {
		return x.RecvBuffAddr
	}
	return nil
}

func (x *BeginReceiveRequest) GetNumBytes() uint64 {
	if x != nil {
		return x.NumBytes
	}
	return 0
}

func (x *BeginReceiveRequest) GetSrcRank() *Rank {
	if x != nil {
		return x.SrcRank
	}
	return nil
}

type BeginReceiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Initiated bool `protobuf:"varint,1,opt,name=initiated,proto3" json:"initiated,omitempty"`
}

func (x *BeginReceiveResponse) Reset() {
	*x = BeginReceiveResponse{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BeginReceiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BeginReceiveResponse) ProtoMessage() {}

func (x *BeginReceiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BeginReceiveResponse.ProtoReflect.Descriptor instead.
func (*BeginReceiveResponse) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{10}
}

func (x *BeginReceiveResponse) GetInitiated() bool {
	if x != nil {
		return x.Initiated
	}
	return false
}

type DataChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"` // You may add more fields here
}

func (x *DataChunk) Reset() {
	*x = DataChunk{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataChunk) ProtoMessage() {}

func (x *DataChunk) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataChunk.ProtoReflect.Descriptor instead.
func (*DataChunk) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{11}
}

func (x *DataChunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type StreamSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *StreamSendResponse) Reset() {
	*x = StreamSendResponse{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamSendResponse) ProtoMessage() {}

func (x *StreamSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamSendResponse.ProtoReflect.Descriptor instead.
func (*StreamSendResponse) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{12}
}

func (x *StreamSendResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type GetStreamStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StreamId *StreamId `protobuf:"bytes,1,opt,name=streamId,proto3" json:"streamId,omitempty"`
}

func (x *GetStreamStatusRequest) Reset() {
	*x = GetStreamStatusRequest{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStreamStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamStatusRequest) ProtoMessage() {}

func (x *GetStreamStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamStatusRequest.ProtoReflect.Descriptor instead.
func (*GetStreamStatusRequest) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{13}
}

func (x *GetStreamStatusRequest) GetStreamId() *StreamId {
	if x != nil {
		return x.StreamId
	}
	return nil
}

type GetStreamStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status StreamStatus `protobuf:"varint,1,opt,name=status,proto3,enum=gpu_sim.StreamStatus" json:"status,omitempty"`
}

func (x *GetStreamStatusResponse) Reset() {
	*x = GetStreamStatusResponse{}
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStreamStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamStatusResponse) ProtoMessage() {}

func (x *GetStreamStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_device_proto_gpu_device_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamStatusResponse.ProtoReflect.Descriptor instead.
func (*GetStreamStatusResponse) Descriptor() ([]byte, []int) {
	return file_gpu_device_proto_gpu_device_proto_rawDescGZIP(), []int{14}
}

func (x *GetStreamStatusResponse) GetStatus() StreamStatus {
	if x != nil {
		return x.Status
	}
	return StreamStatus_IN_PROGRESS
}

var File_gpu_device_proto_gpu_device_proto protoreflect.FileDescriptor

var file_gpu_device_proto_gpu_device_proto_rawDesc = []byte{
	0x0a, 0x21, 0x67, 0x70, 0x75, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x70, 0x75, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x22, 0x20, 0x0a, 0x08,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x1c,
	0x0a, 0x04, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x1f, 0x0a, 0x07,
	0x4d, 0x65, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a,
	0x08, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0xa3, 0x01, 0x0a, 0x0e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x2d, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x52, 0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x30, 0x0a, 0x0a, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e,
	0x4d, 0x65, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x52, 0x0a, 0x6d, 0x69, 0x6e, 0x4d, 0x65, 0x6d, 0x41,
	0x64, 0x64, 0x72, 0x12, 0x30, 0x0a, 0x0a, 0x6d, 0x61, 0x78, 0x4d, 0x65, 0x6d, 0x41, 0x64, 0x64,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69,
	0x6d, 0x2e, 0x4d, 0x65, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x52, 0x0a, 0x6d, 0x61, 0x78, 0x4d, 0x65,
	0x6d, 0x41, 0x64, 0x64, 0x72, 0x22, 0x1a, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0xf1, 0x01, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x33, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x5d, 0x0a, 0x0f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x74, 0x6f, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e,
	0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x52, 0x61, 0x6e, 0x6b, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0d, 0x72, 0x61, 0x6e, 0x6b, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x1a, 0x40, 0x0a, 0x12, 0x52, 0x61, 0x6e, 0x6b, 0x54, 0x6f, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8d, 0x01, 0x0a, 0x10, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0c, 0x73, 0x65,
	0x6e, 0x64, 0x42, 0x75, 0x66, 0x66, 0x41, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x4d, 0x65, 0x6d, 0x41, 0x64,
	0x64, 0x72, 0x52, 0x0c, 0x73, 0x65, 0x6e, 0x64, 0x42, 0x75, 0x66, 0x66, 0x41, 0x64, 0x64, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x08, 0x6e, 0x75, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x27, 0x0a, 0x07,
	0x64, 0x73, 0x74, 0x52, 0x61, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x52, 0x61, 0x6e, 0x6b, 0x52, 0x07, 0x64, 0x73,
	0x74, 0x52, 0x61, 0x6e, 0x6b, 0x22, 0x60, 0x0a, 0x11, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x53, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x64, 0x12, 0x2d, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x70, 0x75,
	0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x52, 0x08, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x22, 0xbf, 0x01, 0x0a, 0x13, 0x42, 0x65, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2d, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x49, 0x64, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x34,
	0x0a, 0x0c, 0x72, 0x65, 0x63, 0x76, 0x42, 0x75, 0x66, 0x66, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x4d,
	0x65, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x76, 0x42, 0x75, 0x66, 0x66,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6e, 0x75, 0x6d, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x12, 0x27, 0x0a, 0x07, 0x73, 0x72, 0x63, 0x52, 0x61, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x52, 0x61, 0x6e, 0x6b,
	0x52, 0x07, 0x73, 0x72, 0x63, 0x52, 0x61, 0x6e, 0x6b, 0x22, 0x34, 0x0a, 0x14, 0x42, 0x65, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65, 0x64, 0x22,
	0x1f, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x2e, 0x0a, 0x12, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x47, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x08, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67,
	0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x52,
	0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x17, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2a, 0x38, 0x0a, 0x0c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45,
	0x53, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x32, 0x99, 0x03,
	0x0a, 0x09, 0x47, 0x50, 0x55, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x21, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x47, 0x65,
	0x74, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x09, 0x42, 0x65, 0x67,
	0x69, 0x6e, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x19, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d,
	0x2e, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1a, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x42, 0x65, 0x67, 0x69,
	0x6e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4d, 0x0a, 0x0c, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x12,
	0x1c, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x41,
	0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x12, 0x2e, 0x67,
	0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x1a, 0x1b, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28,
	0x01, 0x12, 0x56, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x67, 0x70, 0x75, 0x5f, 0x73, 0x69, 0x6d, 0x2e,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4a, 0x6f, 0x73, 0x68, 0x75, 0x61, 0x4d, 0x42,
	0x61, 0x2f, 0x64, 0x73, 0x6d, 0x6c, 0x2f, 0x67, 0x70, 0x75, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gpu_device_proto_gpu_device_proto_rawDescOnce sync.Once
	file_gpu_device_proto_gpu_device_proto_rawDescData = file_gpu_device_proto_gpu_device_proto_rawDesc
)

func file_gpu_device_proto_gpu_device_proto_rawDescGZIP() []byte {
	file_gpu_device_proto_gpu_device_proto_rawDescOnce.Do(func() {
		file_gpu_device_proto_gpu_device_proto_rawDescData = protoimpl.X.CompressGZIP(file_gpu_device_proto_gpu_device_proto_rawDescData)
	})
	return file_gpu_device_proto_gpu_device_proto_rawDescData
}

var file_gpu_device_proto_gpu_device_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gpu_device_proto_gpu_device_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_gpu_device_proto_gpu_device_proto_goTypes = []any{
	(StreamStatus)(0),                 // 0: gpu_sim.StreamStatus
	(*DeviceId)(nil),                  // 1: gpu_sim.DeviceId
	(*Rank)(nil),                      // 2: gpu_sim.Rank
	(*MemAddr)(nil),                   // 3: gpu_sim.MemAddr
	(*StreamId)(nil),                  // 4: gpu_sim.StreamId
	(*DeviceMetadata)(nil),            // 5: gpu_sim.DeviceMetadata
	(*GetDeviceMetadataRequest)(nil),  // 6: gpu_sim.GetDeviceMetadataRequest
	(*GetDeviceMetadataResponse)(nil), // 7: gpu_sim.GetDeviceMetadataResponse
	(*BeginSendRequest)(nil),          // 8: gpu_sim.BeginSendRequest
	(*BeginSendResponse)(nil),         // 9: gpu_sim.BeginSendResponse
	(*BeginReceiveRequest)(nil),       // 10: gpu_sim.BeginReceiveRequest
	(*BeginReceiveResponse)(nil),      // 11: gpu_sim.BeginReceiveResponse
	(*DataChunk)(nil),                 // 12: gpu_sim.DataChunk
	(*StreamSendResponse)(nil),        // 13: gpu_sim.StreamSendResponse
	(*GetStreamStatusRequest)(nil),    // 14: gpu_sim.GetStreamStatusRequest
	(*GetStreamStatusResponse)(nil),   // 15: gpu_sim.GetStreamStatusResponse
	nil,                               // 16: gpu_sim.GetDeviceMetadataResponse.RankToAddressEntry
}
var file_gpu_device_proto_gpu_device_proto_depIdxs = []int32{
	1,  // 0: gpu_sim.DeviceMetadata.deviceId:type_name -> gpu_sim.DeviceId
	3,  // 1: gpu_sim.DeviceMetadata.minMemAddr:type_name -> gpu_sim.MemAddr
	3,  // 2: gpu_sim.DeviceMetadata.maxMemAddr:type_name -> gpu_sim.MemAddr
	5,  // 3: gpu_sim.GetDeviceMetadataResponse.metadata:type_name -> gpu_sim.DeviceMetadata
	16, // 4: gpu_sim.GetDeviceMetadataResponse.rank_to_address:type_name -> gpu_sim.GetDeviceMetadataResponse.RankToAddressEntry
	3,  // 5: gpu_sim.BeginSendRequest.sendBuffAddr:type_name -> gpu_sim.MemAddr
	2,  // 6: gpu_sim.BeginSendRequest.dstRank:type_name -> gpu_sim.Rank
	4,  // 7: gpu_sim.BeginSendResponse.streamId:type_name -> gpu_sim.StreamId
	4,  // 8: gpu_sim.BeginReceiveRequest.streamId:type_name -> gpu_sim.StreamId
	3,  // 9: gpu_sim.BeginReceiveRequest.recvBuffAddr:type_name -> gpu_sim.MemAddr
	2,  // 10: gpu_sim.BeginReceiveRequest.srcRank:type_name -> gpu_sim.Rank
	4,  // 11: gpu_sim.GetStreamStatusRequest.streamId:type_name -> gpu_sim.StreamId
	0,  // 12: gpu_sim.GetStreamStatusResponse.status:type_name -> gpu_sim.StreamStatus
	6,  // 13: gpu_sim.GPUDevice.GetDeviceMetadata:input_type -> gpu_sim.GetDeviceMetadataRequest
	8,  // 14: gpu_sim.GPUDevice.BeginSend:input_type -> gpu_sim.BeginSendRequest
	10, // 15: gpu_sim.GPUDevice.BeginReceive:input_type -> gpu_sim.BeginReceiveRequest
	12, // 16: gpu_sim.GPUDevice.StreamSend:input_type -> gpu_sim.DataChunk
	14, // 17: gpu_sim.GPUDevice.GetStreamStatus:input_type -> gpu_sim.GetStreamStatusRequest
	7,  // 18: gpu_sim.GPUDevice.GetDeviceMetadata:output_type -> gpu_sim.GetDeviceMetadataResponse
	9,  // 19: gpu_sim.GPUDevice.BeginSend:output_type -> gpu_sim.BeginSendResponse
	11, // 20: gpu_sim.GPUDevice.BeginReceive:output_type -> gpu_sim.BeginReceiveResponse
	13, // 21: gpu_sim.GPUDevice.StreamSend:output_type -> gpu_sim.StreamSendResponse
	15, // 22: gpu_sim.GPUDevice.GetStreamStatus:output_type -> gpu_sim.GetStreamStatusResponse
	18, // [18:23] is the sub-list for method output_type
	13, // [13:18] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_gpu_device_proto_gpu_device_proto_init() }
func file_gpu_device_proto_gpu_device_proto_init() {
	if File_gpu_device_proto_gpu_device_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gpu_device_proto_gpu_device_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gpu_device_proto_gpu_device_proto_goTypes,
		DependencyIndexes: file_gpu_device_proto_gpu_device_proto_depIdxs,
		EnumInfos:         file_gpu_device_proto_gpu_device_proto_enumTypes,
		MessageInfos:      file_gpu_device_proto_gpu_device_proto_msgTypes,
	}.Build()
	File_gpu_device_proto_gpu_device_proto = out.File
	file_gpu_device_proto_gpu_device_proto_rawDesc = nil
	file_gpu_device_proto_gpu_device_proto_goTypes = nil
	file_gpu_device_proto_gpu_device_proto_depIdxs = nil
}
