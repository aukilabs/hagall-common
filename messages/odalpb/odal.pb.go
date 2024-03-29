// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: messages/odalpb/odal.proto

package odalpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MsgType int32

const (
	MsgType_MSG_TYPE_ERROR_RESPONSE                    MsgType = 0
	MsgType_MSG_TYPE_ODAL_STATE                        MsgType = 200
	MsgType_MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_REQUEST   MsgType = 201
	MsgType_MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_RESPONSE  MsgType = 202
	MsgType_MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_BROADCAST MsgType = 203
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0:   "MSG_TYPE_ERROR_RESPONSE",
		200: "MSG_TYPE_ODAL_STATE",
		201: "MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_REQUEST",
		202: "MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_RESPONSE",
		203: "MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_BROADCAST",
	}
	MsgType_value = map[string]int32{
		"MSG_TYPE_ERROR_RESPONSE":                    0,
		"MSG_TYPE_ODAL_STATE":                        200,
		"MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_REQUEST":   201,
		"MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_RESPONSE":  202,
		"MSG_TYPE_ODAL_ASSET_INSTANCE_ADD_BROADCAST": 203,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_messages_odalpb_odal_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_messages_odalpb_odal_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{0}
}

// An instance of an asset bound by a participant and tied to an entity created
// together with the asset instance.
type AssetInstance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The asset instance id.
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The asset id.
	AssetId string `protobuf:"bytes,2,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	// The id of the participant that owns the asset instance.
	ParticipantId uint32 `protobuf:"varint,3,opt,name=participant_id,json=participantId,proto3" json:"participant_id,omitempty"`
	// The entity id.
	EntityId uint32 `protobuf:"varint,4,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
}

func (x *AssetInstance) Reset() {
	*x = AssetInstance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_odalpb_odal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInstance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInstance) ProtoMessage() {}

func (x *AssetInstance) ProtoReflect() protoreflect.Message {
	mi := &file_messages_odalpb_odal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetInstance.ProtoReflect.Descriptor instead.
func (*AssetInstance) Descriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{0}
}

func (x *AssetInstance) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AssetInstance) GetAssetId() string {
	if x != nil {
		return x.AssetId
	}
	return ""
}

func (x *AssetInstance) GetParticipantId() uint32 {
	if x != nil {
		return x.ParticipantId
	}
	return 0
}

func (x *AssetInstance) GetEntityId() uint32 {
	if x != nil {
		return x.EntityId
	}
	return 0
}

// State represents a odal state within a session.
type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the message.
	Type MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=odal.MsgType" json:"type,omitempty"`
	// The time the message is sent.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The assetInstances present in the session.
	AssetInstances []*AssetInstance `protobuf:"bytes,3,rep,name=asset_instances,json=assetInstances,proto3" json:"asset_instances,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_odalpb_odal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_messages_odalpb_odal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{1}
}

func (x *State) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_MSG_TYPE_ERROR_RESPONSE
}

func (x *State) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *State) GetAssetInstances() []*AssetInstance {
	if x != nil {
		return x.AssetInstances
	}
	return nil
}

// AssetInstanceAddRequest represents a message to add an asset instance to the
// current session.
type AssetInstanceAddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the message.
	Type MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=odal.MsgType" json:"type,omitempty"`
	// The time the message is sent.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The id that identifies a request.
	RequestId uint32 `protobuf:"varint,1337,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// The id of the entity to associate with the asset.
	EntityId uint32 `protobuf:"varint,3,opt,name=entity_id,json=entityId,proto3" json:"entity_id,omitempty"`
	// The asset id.
	AssetId string `protobuf:"bytes,4,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
}

func (x *AssetInstanceAddRequest) Reset() {
	*x = AssetInstanceAddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_odalpb_odal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInstanceAddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInstanceAddRequest) ProtoMessage() {}

func (x *AssetInstanceAddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_odalpb_odal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetInstanceAddRequest.ProtoReflect.Descriptor instead.
func (*AssetInstanceAddRequest) Descriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{2}
}

func (x *AssetInstanceAddRequest) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_MSG_TYPE_ERROR_RESPONSE
}

func (x *AssetInstanceAddRequest) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *AssetInstanceAddRequest) GetRequestId() uint32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *AssetInstanceAddRequest) GetEntityId() uint32 {
	if x != nil {
		return x.EntityId
	}
	return 0
}

func (x *AssetInstanceAddRequest) GetAssetId() string {
	if x != nil {
		return x.AssetId
	}
	return ""
}

// AssetInstanceAddResponse is a message returned in response to an add asset
// instance request.
type AssetInstanceAddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the message.
	Type MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=odal.MsgType" json:"type,omitempty"`
	// The time the message is sent.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The id that identifies a request.
	RequestId uint32 `protobuf:"varint,1337,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// The attributed id to the added asset instance.
	AssetInstanceId uint32 `protobuf:"varint,3,opt,name=asset_instance_id,json=assetInstanceId,proto3" json:"asset_instance_id,omitempty"`
}

func (x *AssetInstanceAddResponse) Reset() {
	*x = AssetInstanceAddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_odalpb_odal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInstanceAddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInstanceAddResponse) ProtoMessage() {}

func (x *AssetInstanceAddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messages_odalpb_odal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetInstanceAddResponse.ProtoReflect.Descriptor instead.
func (*AssetInstanceAddResponse) Descriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{3}
}

func (x *AssetInstanceAddResponse) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_MSG_TYPE_ERROR_RESPONSE
}

func (x *AssetInstanceAddResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *AssetInstanceAddResponse) GetRequestId() uint32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *AssetInstanceAddResponse) GetAssetInstanceId() uint32 {
	if x != nil {
		return x.AssetInstanceId
	}
	return 0
}

// AssetInstanceAddBroadcast is a message that notifies that an asset instance
// has been added.
type AssetInstanceAddBroadcast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the message.
	Type MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=odal.MsgType" json:"type,omitempty"`
	// The time the message is sent.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// The timestamp of the request that triggered the broadcast.
	OriginTimestamp *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=origin_timestamp,json=originTimestamp,proto3" json:"origin_timestamp,omitempty"`
	// The id of the added asset instance.
	AssetInstance *AssetInstance `protobuf:"bytes,4,opt,name=asset_instance,json=assetInstance,proto3" json:"asset_instance,omitempty"`
}

func (x *AssetInstanceAddBroadcast) Reset() {
	*x = AssetInstanceAddBroadcast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_odalpb_odal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetInstanceAddBroadcast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetInstanceAddBroadcast) ProtoMessage() {}

func (x *AssetInstanceAddBroadcast) ProtoReflect() protoreflect.Message {
	mi := &file_messages_odalpb_odal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetInstanceAddBroadcast.ProtoReflect.Descriptor instead.
func (*AssetInstanceAddBroadcast) Descriptor() ([]byte, []int) {
	return file_messages_odalpb_odal_proto_rawDescGZIP(), []int{4}
}

func (x *AssetInstanceAddBroadcast) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_MSG_TYPE_ERROR_RESPONSE
}

func (x *AssetInstanceAddBroadcast) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *AssetInstanceAddBroadcast) GetOriginTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.OriginTimestamp
	}
	return nil
}

func (x *AssetInstanceAddBroadcast) GetAssetInstance() *AssetInstance {
	if x != nil {
		return x.AssetInstance
	}
	return nil
}

var File_messages_odalpb_odal_proto protoreflect.FileDescriptor

var file_messages_odalpb_odal_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6f, 0x64, 0x61, 0x6c, 0x70,
	0x62, 0x2f, 0x6f, 0x64, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6f, 0x64,
	0x61, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x7e, 0x0a, 0x0d, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x64, 0x12,
	0x25, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69,
	0x70, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x49, 0x64, 0x22, 0xa2, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x21, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x6f, 0x64,
	0x61, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3c, 0x0a, 0x0f, 0x61, 0x73,
	0x73, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6f, 0x64, 0x61, 0x6c, 0x2e, 0x41, 0x73, 0x73, 0x65, 0x74,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x0e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22, 0xce, 0x01, 0x0a, 0x17, 0x41, 0x73, 0x73,
	0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x6f, 0x64, 0x61, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0xb9, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x19,
	0x0a, 0x08, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x64, 0x22, 0xc3, 0x01, 0x0a, 0x18, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x41, 0x64, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x6f, 0x64, 0x61, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0xb9, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0f,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x22,
	0xfb, 0x01, 0x0a, 0x19, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x41, 0x64, 0x64, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x21, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x6f, 0x64,
	0x61, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x45, 0x0a, 0x10, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x3a, 0x0a, 0x0e, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6f, 0x64, 0x61, 0x6c,
	0x2e, 0x41, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x0d,
	0x61, 0x73, 0x73, 0x65, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x2a, 0xe2, 0x01,
	0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x4d, 0x53, 0x47,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x50,
	0x4f, 0x4e, 0x53, 0x45, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x13, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4f, 0x44, 0x41, 0x4c, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x10, 0xc8, 0x01,
	0x12, 0x2d, 0x0a, 0x28, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x44, 0x41,
	0x4c, 0x5f, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4e, 0x43, 0x45,
	0x5f, 0x41, 0x44, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x10, 0xc9, 0x01, 0x12,
	0x2e, 0x0a, 0x29, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x44, 0x41, 0x4c,
	0x5f, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4e, 0x43, 0x45, 0x5f,
	0x41, 0x44, 0x44, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0xca, 0x01, 0x12,
	0x2f, 0x0a, 0x2a, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4f, 0x44, 0x41, 0x4c,
	0x5f, 0x41, 0x53, 0x53, 0x45, 0x54, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4e, 0x43, 0x45, 0x5f,
	0x41, 0x44, 0x44, 0x5f, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x10, 0xcb, 0x01,
	0x22, 0x05, 0x08, 0x01, 0x10, 0xc7, 0x01, 0x22, 0x09, 0x08, 0xac, 0x02, 0x10, 0xff, 0xff, 0xff,
	0xff, 0x07, 0x42, 0x3c, 0x5a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6f,
	0x64, 0x61, 0x6c, 0x70, 0x62, 0xa2, 0x02, 0x04, 0x4f, 0x64, 0x61, 0x6c, 0xaa, 0x02, 0x21, 0x41,
	0x75, 0x6b, 0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x6a, 0x75, 0x72, 0x65, 0x4b, 0x69, 0x74, 0x2e, 0x4f,
	0x64, 0x61, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x6e,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_odalpb_odal_proto_rawDescOnce sync.Once
	file_messages_odalpb_odal_proto_rawDescData = file_messages_odalpb_odal_proto_rawDesc
)

func file_messages_odalpb_odal_proto_rawDescGZIP() []byte {
	file_messages_odalpb_odal_proto_rawDescOnce.Do(func() {
		file_messages_odalpb_odal_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_odalpb_odal_proto_rawDescData)
	})
	return file_messages_odalpb_odal_proto_rawDescData
}

var file_messages_odalpb_odal_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messages_odalpb_odal_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_messages_odalpb_odal_proto_goTypes = []interface{}{
	(MsgType)(0),                      // 0: odal.MsgType
	(*AssetInstance)(nil),             // 1: odal.AssetInstance
	(*State)(nil),                     // 2: odal.State
	(*AssetInstanceAddRequest)(nil),   // 3: odal.AssetInstanceAddRequest
	(*AssetInstanceAddResponse)(nil),  // 4: odal.AssetInstanceAddResponse
	(*AssetInstanceAddBroadcast)(nil), // 5: odal.AssetInstanceAddBroadcast
	(*timestamppb.Timestamp)(nil),     // 6: google.protobuf.Timestamp
}
var file_messages_odalpb_odal_proto_depIdxs = []int32{
	0,  // 0: odal.State.type:type_name -> odal.MsgType
	6,  // 1: odal.State.timestamp:type_name -> google.protobuf.Timestamp
	1,  // 2: odal.State.asset_instances:type_name -> odal.AssetInstance
	0,  // 3: odal.AssetInstanceAddRequest.type:type_name -> odal.MsgType
	6,  // 4: odal.AssetInstanceAddRequest.timestamp:type_name -> google.protobuf.Timestamp
	0,  // 5: odal.AssetInstanceAddResponse.type:type_name -> odal.MsgType
	6,  // 6: odal.AssetInstanceAddResponse.timestamp:type_name -> google.protobuf.Timestamp
	0,  // 7: odal.AssetInstanceAddBroadcast.type:type_name -> odal.MsgType
	6,  // 8: odal.AssetInstanceAddBroadcast.timestamp:type_name -> google.protobuf.Timestamp
	6,  // 9: odal.AssetInstanceAddBroadcast.origin_timestamp:type_name -> google.protobuf.Timestamp
	1,  // 10: odal.AssetInstanceAddBroadcast.asset_instance:type_name -> odal.AssetInstance
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_messages_odalpb_odal_proto_init() }
func file_messages_odalpb_odal_proto_init() {
	if File_messages_odalpb_odal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_odalpb_odal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInstance); i {
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
		file_messages_odalpb_odal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_messages_odalpb_odal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInstanceAddRequest); i {
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
		file_messages_odalpb_odal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInstanceAddResponse); i {
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
		file_messages_odalpb_odal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssetInstanceAddBroadcast); i {
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
			RawDescriptor: file_messages_odalpb_odal_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_odalpb_odal_proto_goTypes,
		DependencyIndexes: file_messages_odalpb_odal_proto_depIdxs,
		EnumInfos:         file_messages_odalpb_odal_proto_enumTypes,
		MessageInfos:      file_messages_odalpb_odal_proto_msgTypes,
	}.Build()
	File_messages_odalpb_odal_proto = out.File
	file_messages_odalpb_odal_proto_rawDesc = nil
	file_messages_odalpb_odal_proto_goTypes = nil
	file_messages_odalpb_odal_proto_depIdxs = nil
}
