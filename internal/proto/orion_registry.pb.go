// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.2
// source: proto/orion_registry.proto

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

type NewMemberEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FriendlyName string `protobuf:"bytes,1,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	PeerId       uint32 `protobuf:"varint,2,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
}

func (x *NewMemberEvent) Reset() {
	*x = NewMemberEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewMemberEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewMemberEvent) ProtoMessage() {}

func (x *NewMemberEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewMemberEvent.ProtoReflect.Descriptor instead.
func (*NewMemberEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{0}
}

func (x *NewMemberEvent) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *NewMemberEvent) GetPeerId() uint32 {
	if x != nil {
		return x.PeerId
	}
	return 0
}

type MemberDisconnectedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FriendlyName string `protobuf:"bytes,1,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	PeerId       uint32 `protobuf:"varint,2,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
}

func (x *MemberDisconnectedEvent) Reset() {
	*x = MemberDisconnectedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberDisconnectedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberDisconnectedEvent) ProtoMessage() {}

func (x *MemberDisconnectedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberDisconnectedEvent.ProtoReflect.Descriptor instead.
func (*MemberDisconnectedEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{1}
}

func (x *MemberDisconnectedEvent) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *MemberDisconnectedEvent) GetPeerId() uint32 {
	if x != nil {
		return x.PeerId
	}
	return 0
}

type MemberConnectEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EndpointAddr      string `protobuf:"bytes,1,opt,name=endpoint_addr,json=endpointAddr,proto3" json:"endpoint_addr,omitempty"`
	EndpointPort      uint32 `protobuf:"varint,2,opt,name=endpoint_port,json=endpointPort,proto3" json:"endpoint_port,omitempty"`
	PublicKey         []byte `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	FriendlyName      string `protobuf:"bytes,5,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	DestinationPeerId uint32 `protobuf:"varint,6,opt,name=destination_peer_id,json=destinationPeerId,proto3" json:"destination_peer_id,omitempty"`
	SourcePeerId      uint32 `protobuf:"varint,7,opt,name=source_peer_id,json=sourcePeerId,proto3" json:"source_peer_id,omitempty"`
}

func (x *MemberConnectEvent) Reset() {
	*x = MemberConnectEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberConnectEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberConnectEvent) ProtoMessage() {}

func (x *MemberConnectEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberConnectEvent.ProtoReflect.Descriptor instead.
func (*MemberConnectEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{2}
}

func (x *MemberConnectEvent) GetEndpointAddr() string {
	if x != nil {
		return x.EndpointAddr
	}
	return ""
}

func (x *MemberConnectEvent) GetEndpointPort() uint32 {
	if x != nil {
		return x.EndpointPort
	}
	return 0
}

func (x *MemberConnectEvent) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *MemberConnectEvent) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *MemberConnectEvent) GetDestinationPeerId() uint32 {
	if x != nil {
		return x.DestinationPeerId
	}
	return 0
}

func (x *MemberConnectEvent) GetSourcePeerId() uint32 {
	if x != nil {
		return x.SourcePeerId
	}
	return 0
}

type MemberConnectResponseEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EndpointAddr      string `protobuf:"bytes,1,opt,name=endpoint_addr,json=endpointAddr,proto3" json:"endpoint_addr,omitempty"`
	EndpointPort      uint32 `protobuf:"varint,2,opt,name=endpoint_port,json=endpointPort,proto3" json:"endpoint_port,omitempty"`
	PublicKey         []byte `protobuf:"bytes,3,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	FriendlyName      string `protobuf:"bytes,4,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	DestinationPeerId uint32 `protobuf:"varint,5,opt,name=destination_peer_id,json=destinationPeerId,proto3" json:"destination_peer_id,omitempty"`
	PresharedKey      []byte `protobuf:"bytes,7,opt,name=preshared_key,json=presharedKey,proto3" json:"preshared_key,omitempty"`
	SourcePeerId      uint32 `protobuf:"varint,6,opt,name=source_peer_id,json=sourcePeerId,proto3" json:"source_peer_id,omitempty"`
}

func (x *MemberConnectResponseEvent) Reset() {
	*x = MemberConnectResponseEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberConnectResponseEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberConnectResponseEvent) ProtoMessage() {}

func (x *MemberConnectResponseEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberConnectResponseEvent.ProtoReflect.Descriptor instead.
func (*MemberConnectResponseEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{3}
}

func (x *MemberConnectResponseEvent) GetEndpointAddr() string {
	if x != nil {
		return x.EndpointAddr
	}
	return ""
}

func (x *MemberConnectResponseEvent) GetEndpointPort() uint32 {
	if x != nil {
		return x.EndpointPort
	}
	return 0
}

func (x *MemberConnectResponseEvent) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *MemberConnectResponseEvent) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *MemberConnectResponseEvent) GetDestinationPeerId() uint32 {
	if x != nil {
		return x.DestinationPeerId
	}
	return 0
}

func (x *MemberConnectResponseEvent) GetPresharedKey() []byte {
	if x != nil {
		return x.PresharedKey
	}
	return nil
}

func (x *MemberConnectResponseEvent) GetSourcePeerId() uint32 {
	if x != nil {
		return x.SourcePeerId
	}
	return 0
}

type SessionIDIssued struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

func (x *SessionIDIssued) Reset() {
	*x = SessionIDIssued{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionIDIssued) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionIDIssued) ProtoMessage() {}

func (x *SessionIDIssued) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionIDIssued.ProtoReflect.Descriptor instead.
func (*SessionIDIssued) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{4}
}

func (x *SessionIDIssued) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type RPCServerEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//
	//	*RPCServerEvent_NewMember
	//	*RPCServerEvent_DisconnectedMember
	//	*RPCServerEvent_MemberConnect
	//	*RPCServerEvent_MemberConnectResponse
	//	*RPCServerEvent_SessionId
	Event isRPCServerEvent_Event `protobuf_oneof:"event"`
}

func (x *RPCServerEvent) Reset() {
	*x = RPCServerEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCServerEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCServerEvent) ProtoMessage() {}

func (x *RPCServerEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCServerEvent.ProtoReflect.Descriptor instead.
func (*RPCServerEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{5}
}

func (m *RPCServerEvent) GetEvent() isRPCServerEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *RPCServerEvent) GetNewMember() *NewMemberEvent {
	if x, ok := x.GetEvent().(*RPCServerEvent_NewMember); ok {
		return x.NewMember
	}
	return nil
}

func (x *RPCServerEvent) GetDisconnectedMember() *MemberDisconnectedEvent {
	if x, ok := x.GetEvent().(*RPCServerEvent_DisconnectedMember); ok {
		return x.DisconnectedMember
	}
	return nil
}

func (x *RPCServerEvent) GetMemberConnect() *MemberConnectEvent {
	if x, ok := x.GetEvent().(*RPCServerEvent_MemberConnect); ok {
		return x.MemberConnect
	}
	return nil
}

func (x *RPCServerEvent) GetMemberConnectResponse() *MemberConnectResponseEvent {
	if x, ok := x.GetEvent().(*RPCServerEvent_MemberConnectResponse); ok {
		return x.MemberConnectResponse
	}
	return nil
}

func (x *RPCServerEvent) GetSessionId() *SessionIDIssued {
	if x, ok := x.GetEvent().(*RPCServerEvent_SessionId); ok {
		return x.SessionId
	}
	return nil
}

type isRPCServerEvent_Event interface {
	isRPCServerEvent_Event()
}

type RPCServerEvent_NewMember struct {
	NewMember *NewMemberEvent `protobuf:"bytes,1,opt,name=new_member,json=newMember,proto3,oneof"`
}

type RPCServerEvent_DisconnectedMember struct {
	DisconnectedMember *MemberDisconnectedEvent `protobuf:"bytes,2,opt,name=disconnected_member,json=disconnectedMember,proto3,oneof"`
}

type RPCServerEvent_MemberConnect struct {
	MemberConnect *MemberConnectEvent `protobuf:"bytes,3,opt,name=member_connect,json=memberConnect,proto3,oneof"`
}

type RPCServerEvent_MemberConnectResponse struct {
	MemberConnectResponse *MemberConnectResponseEvent `protobuf:"bytes,4,opt,name=member_connect_response,json=memberConnectResponse,proto3,oneof"`
}

type RPCServerEvent_SessionId struct {
	SessionId *SessionIDIssued `protobuf:"bytes,5,opt,name=session_id,json=sessionId,proto3,oneof"`
}

func (*RPCServerEvent_NewMember) isRPCServerEvent_Event() {}

func (*RPCServerEvent_DisconnectedMember) isRPCServerEvent_Event() {}

func (*RPCServerEvent_MemberConnect) isRPCServerEvent_Event() {}

func (*RPCServerEvent_MemberConnectResponse) isRPCServerEvent_Event() {}

func (*RPCServerEvent_SessionId) isRPCServerEvent_Event() {}

type InitializeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FriendlyName    string `protobuf:"bytes,1,opt,name=friendly_name,json=friendlyName,proto3" json:"friendly_name,omitempty"`
	TimestampSigned int64  `protobuf:"varint,2,opt,name=timestamp_signed,json=timestampSigned,proto3" json:"timestamp_signed,omitempty"`
	Signed          []byte `protobuf:"bytes,3,opt,name=signed,proto3" json:"signed,omitempty"`
	MemberId        uint32 `protobuf:"varint,4,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	Certificate     []byte `protobuf:"bytes,5,opt,name=certificate,proto3" json:"certificate,omitempty"`
	SessionId       string `protobuf:"bytes,6,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Reconnect       *bool  `protobuf:"varint,7,opt,name=reconnect,proto3,oneof" json:"reconnect,omitempty"`
}

func (x *InitializeRequest) Reset() {
	*x = InitializeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitializeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeRequest) ProtoMessage() {}

func (x *InitializeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeRequest.ProtoReflect.Descriptor instead.
func (*InitializeRequest) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{6}
}

func (x *InitializeRequest) GetFriendlyName() string {
	if x != nil {
		return x.FriendlyName
	}
	return ""
}

func (x *InitializeRequest) GetTimestampSigned() int64 {
	if x != nil {
		return x.TimestampSigned
	}
	return 0
}

func (x *InitializeRequest) GetSigned() []byte {
	if x != nil {
		return x.Signed
	}
	return nil
}

func (x *InitializeRequest) GetMemberId() uint32 {
	if x != nil {
		return x.MemberId
	}
	return 0
}

func (x *InitializeRequest) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

func (x *InitializeRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *InitializeRequest) GetReconnect() bool {
	if x != nil && x.Reconnect != nil {
		return *x.Reconnect
	}
	return false
}

type RPCClientEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//
	//	*RPCClientEvent_Initialize
	//	*RPCClientEvent_Connect
	//	*RPCClientEvent_ConnectResponse
	Event isRPCClientEvent_Event `protobuf_oneof:"event"`
}

func (x *RPCClientEvent) Reset() {
	*x = RPCClientEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_orion_registry_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPCClientEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPCClientEvent) ProtoMessage() {}

func (x *RPCClientEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_orion_registry_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPCClientEvent.ProtoReflect.Descriptor instead.
func (*RPCClientEvent) Descriptor() ([]byte, []int) {
	return file_proto_orion_registry_proto_rawDescGZIP(), []int{7}
}

func (m *RPCClientEvent) GetEvent() isRPCClientEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *RPCClientEvent) GetInitialize() *InitializeRequest {
	if x, ok := x.GetEvent().(*RPCClientEvent_Initialize); ok {
		return x.Initialize
	}
	return nil
}

func (x *RPCClientEvent) GetConnect() *MemberConnectEvent {
	if x, ok := x.GetEvent().(*RPCClientEvent_Connect); ok {
		return x.Connect
	}
	return nil
}

func (x *RPCClientEvent) GetConnectResponse() *MemberConnectResponseEvent {
	if x, ok := x.GetEvent().(*RPCClientEvent_ConnectResponse); ok {
		return x.ConnectResponse
	}
	return nil
}

type isRPCClientEvent_Event interface {
	isRPCClientEvent_Event()
}

type RPCClientEvent_Initialize struct {
	Initialize *InitializeRequest `protobuf:"bytes,1,opt,name=initialize,proto3,oneof"`
}

type RPCClientEvent_Connect struct {
	Connect *MemberConnectEvent `protobuf:"bytes,2,opt,name=connect,proto3,oneof"`
}

type RPCClientEvent_ConnectResponse struct {
	ConnectResponse *MemberConnectResponseEvent `protobuf:"bytes,3,opt,name=connect_response,json=connectResponse,proto3,oneof"`
}

func (*RPCClientEvent_Initialize) isRPCClientEvent_Event() {}

func (*RPCClientEvent_Connect) isRPCClientEvent_Event() {}

func (*RPCClientEvent_ConnectResponse) isRPCClientEvent_Event() {}

var File_proto_orion_registry_proto protoreflect.FileDescriptor

var file_proto_orion_registry_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x72, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x0e,
	0x4e, 0x65, 0x77, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x22, 0x57, 0x0a, 0x17,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x69, 0x65, 0x6e,
	0x64, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70,
	0x65, 0x65, 0x72, 0x49, 0x64, 0x22, 0xf8, 0x01, 0x0a, 0x12, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c,
	0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72,
	0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x64, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64,
	0x22, 0xa5, 0x02, 0x0a, 0x1a, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x23, 0x0a, 0x0d, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x69, 0x65,
	0x6e, 0x64, 0x6c, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a,
	0x13, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x65, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x70, 0x72, 0x65, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4b,
	0x65, 0x79, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64, 0x22, 0x30, 0x0a, 0x0f, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xe0, 0x02, 0x0a, 0x0e, 0x52,
	0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x0a,
	0x0a, 0x6e, 0x65, 0x77, 0x5f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x4e, 0x65, 0x77, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x65, 0x77, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x4b, 0x0a, 0x13, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x12, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3c, 0x0a, 0x0e,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0d, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x55, 0x0a, 0x17, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x5f, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x4d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x15, 0x6d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x31, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x44, 0x49, 0x73, 0x73, 0x75, 0x65, 0x64, 0x48, 0x00, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x8a, 0x02,
	0x0a, 0x11, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x69, 0x65,
	0x6e, 0x64, 0x6c, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x53, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x09, 0x72, 0x65, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09,
	0x72, 0x65, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a,
	0x5f, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x22, 0xca, 0x01, 0x0a, 0x0e, 0x52,
	0x50, 0x43, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a,
	0x0a, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c,
	0x69, 0x7a, 0x65, 0x12, 0x2f, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x12, 0x48, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x5f,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0f, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x07,
	0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x32, 0x45, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x12, 0x39, 0x0a, 0x11, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65,
	0x54, 0x6f, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0f, 0x2e, 0x52, 0x50, 0x43, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x0f, 0x2e, 0x52, 0x50, 0x43, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x28, 0x01, 0x30, 0x01, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_orion_registry_proto_rawDescOnce sync.Once
	file_proto_orion_registry_proto_rawDescData = file_proto_orion_registry_proto_rawDesc
)

func file_proto_orion_registry_proto_rawDescGZIP() []byte {
	file_proto_orion_registry_proto_rawDescOnce.Do(func() {
		file_proto_orion_registry_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_orion_registry_proto_rawDescData)
	})
	return file_proto_orion_registry_proto_rawDescData
}

var file_proto_orion_registry_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_orion_registry_proto_goTypes = []interface{}{
	(*NewMemberEvent)(nil),             // 0: NewMemberEvent
	(*MemberDisconnectedEvent)(nil),    // 1: MemberDisconnectedEvent
	(*MemberConnectEvent)(nil),         // 2: MemberConnectEvent
	(*MemberConnectResponseEvent)(nil), // 3: MemberConnectResponseEvent
	(*SessionIDIssued)(nil),            // 4: SessionIDIssued
	(*RPCServerEvent)(nil),             // 5: RPCServerEvent
	(*InitializeRequest)(nil),          // 6: InitializeRequest
	(*RPCClientEvent)(nil),             // 7: RPCClientEvent
}
var file_proto_orion_registry_proto_depIdxs = []int32{
	0, // 0: RPCServerEvent.new_member:type_name -> NewMemberEvent
	1, // 1: RPCServerEvent.disconnected_member:type_name -> MemberDisconnectedEvent
	2, // 2: RPCServerEvent.member_connect:type_name -> MemberConnectEvent
	3, // 3: RPCServerEvent.member_connect_response:type_name -> MemberConnectResponseEvent
	4, // 4: RPCServerEvent.session_id:type_name -> SessionIDIssued
	6, // 5: RPCClientEvent.initialize:type_name -> InitializeRequest
	2, // 6: RPCClientEvent.connect:type_name -> MemberConnectEvent
	3, // 7: RPCClientEvent.connect_response:type_name -> MemberConnectResponseEvent
	7, // 8: Registry.SubscribeToStream:input_type -> RPCClientEvent
	5, // 9: Registry.SubscribeToStream:output_type -> RPCServerEvent
	9, // [9:10] is the sub-list for method output_type
	8, // [8:9] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_proto_orion_registry_proto_init() }
func file_proto_orion_registry_proto_init() {
	if File_proto_orion_registry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_orion_registry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewMemberEvent); i {
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
		file_proto_orion_registry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberDisconnectedEvent); i {
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
		file_proto_orion_registry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberConnectEvent); i {
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
		file_proto_orion_registry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberConnectResponseEvent); i {
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
		file_proto_orion_registry_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionIDIssued); i {
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
		file_proto_orion_registry_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCServerEvent); i {
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
		file_proto_orion_registry_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitializeRequest); i {
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
		file_proto_orion_registry_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPCClientEvent); i {
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
	file_proto_orion_registry_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*RPCServerEvent_NewMember)(nil),
		(*RPCServerEvent_DisconnectedMember)(nil),
		(*RPCServerEvent_MemberConnect)(nil),
		(*RPCServerEvent_MemberConnectResponse)(nil),
		(*RPCServerEvent_SessionId)(nil),
	}
	file_proto_orion_registry_proto_msgTypes[6].OneofWrappers = []interface{}{}
	file_proto_orion_registry_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*RPCClientEvent_Initialize)(nil),
		(*RPCClientEvent_Connect)(nil),
		(*RPCClientEvent_ConnectResponse)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_orion_registry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_orion_registry_proto_goTypes,
		DependencyIndexes: file_proto_orion_registry_proto_depIdxs,
		MessageInfos:      file_proto_orion_registry_proto_msgTypes,
	}.Build()
	File_proto_orion_registry_proto = out.File
	file_proto_orion_registry_proto_rawDesc = nil
	file_proto_orion_registry_proto_goTypes = nil
	file_proto_orion_registry_proto_depIdxs = nil
}
