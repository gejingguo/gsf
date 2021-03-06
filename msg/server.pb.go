// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

/*
Package msg is a generated protocol buffer package.

It is generated from these files:
	server.proto

It has these top-level messages:
	ServerInfo
	CmdRegSvrRegReq
	CmdRegSvrRegNtf
*/
package msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Cmd int32

const (
	Cmd_Cmd_ID_None   Cmd = 0
	Cmd_RegSvr_RegReq Cmd = 1001
	Cmd_RegSvr_RegNtf Cmd = 1002
)

var Cmd_name = map[int32]string{
	0:    "Cmd_ID_None",
	1001: "RegSvr_RegReq",
	1002: "RegSvr_RegNtf",
}
var Cmd_value = map[string]int32{
	"Cmd_ID_None":   0,
	"RegSvr_RegReq": 1001,
	"RegSvr_RegNtf": 1002,
}

func (x Cmd) String() string {
	return proto.EnumName(Cmd_name, int32(x))
}
func (Cmd) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ServerType int32

const (
	ServerType_ServerType_None  ServerType = 0
	ServerType_ServerType_Reg   ServerType = 1
	ServerType_ServerType_World ServerType = 2
	ServerType_ServerType_Scene ServerType = 3
	ServerType_ServerType_Gate  ServerType = 4
)

var ServerType_name = map[int32]string{
	0: "ServerType_None",
	1: "ServerType_Reg",
	2: "ServerType_World",
	3: "ServerType_Scene",
	4: "ServerType_Gate",
}
var ServerType_value = map[string]int32{
	"ServerType_None":  0,
	"ServerType_Reg":   1,
	"ServerType_World": 2,
	"ServerType_Scene": 3,
	"ServerType_Gate":  4,
}

func (x ServerType) String() string {
	return proto.EnumName(ServerType_name, int32(x))
}
func (ServerType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ServerInfo struct {
	Id         int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Type       int32  `protobuf:"varint,2,opt,name=type" json:"type,omitempty"`
	Group      int32  `protobuf:"varint,3,opt,name=group" json:"group,omitempty"`
	ListenAddr string `protobuf:"bytes,4,opt,name=listenAddr" json:"listenAddr,omitempty"`
}

func (m *ServerInfo) Reset()                    { *m = ServerInfo{} }
func (m *ServerInfo) String() string            { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()               {}
func (*ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServerInfo) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ServerInfo) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *ServerInfo) GetGroup() int32 {
	if m != nil {
		return m.Group
	}
	return 0
}

func (m *ServerInfo) GetListenAddr() string {
	if m != nil {
		return m.ListenAddr
	}
	return ""
}

type CmdRegSvrRegReq struct {
	Server     *ServerInfo `protobuf:"bytes,1,opt,name=server" json:"server,omitempty"`
	ServerType []int32     `protobuf:"varint,2,rep,packed,name=serverType" json:"serverType,omitempty"`
}

func (m *CmdRegSvrRegReq) Reset()                    { *m = CmdRegSvrRegReq{} }
func (m *CmdRegSvrRegReq) String() string            { return proto.CompactTextString(m) }
func (*CmdRegSvrRegReq) ProtoMessage()               {}
func (*CmdRegSvrRegReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CmdRegSvrRegReq) GetServer() *ServerInfo {
	if m != nil {
		return m.Server
	}
	return nil
}

func (m *CmdRegSvrRegReq) GetServerType() []int32 {
	if m != nil {
		return m.ServerType
	}
	return nil
}

type CmdRegSvrRegNtf struct {
	Server []*ServerInfo `protobuf:"bytes,2,rep,name=server" json:"server,omitempty"`
}

func (m *CmdRegSvrRegNtf) Reset()                    { *m = CmdRegSvrRegNtf{} }
func (m *CmdRegSvrRegNtf) String() string            { return proto.CompactTextString(m) }
func (*CmdRegSvrRegNtf) ProtoMessage()               {}
func (*CmdRegSvrRegNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CmdRegSvrRegNtf) GetServer() []*ServerInfo {
	if m != nil {
		return m.Server
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerInfo)(nil), "msg.ServerInfo")
	proto.RegisterType((*CmdRegSvrRegReq)(nil), "msg.CmdRegSvrRegReq")
	proto.RegisterType((*CmdRegSvrRegNtf)(nil), "msg.CmdRegSvrRegNtf")
	proto.RegisterEnum("msg.Cmd", Cmd_name, Cmd_value)
	proto.RegisterEnum("msg.ServerType", ServerType_name, ServerType_value)
}

func init() { proto.RegisterFile("server.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4b, 0xfb, 0x40,
	0x10, 0x40, 0x7f, 0x49, 0xfa, 0x87, 0xdf, 0x54, 0x9b, 0x65, 0xec, 0x21, 0x27, 0x29, 0xbd, 0x58,
	0x7a, 0xe8, 0x41, 0x6f, 0x1e, 0x04, 0x89, 0x20, 0xbd, 0xe4, 0xb0, 0x11, 0x04, 0x2f, 0x41, 0xdd,
	0xc9, 0x12, 0x68, 0xb2, 0x71, 0xb3, 0x06, 0xfc, 0xc8, 0xfa, 0x29, 0x24, 0xbb, 0x4a, 0xd2, 0x82,
	0xb7, 0x99, 0x97, 0xf0, 0xf6, 0xb1, 0x0b, 0x27, 0x0d, 0xe9, 0x96, 0xf4, 0xb6, 0xd6, 0xca, 0x28,
	0x0c, 0xca, 0x46, 0xae, 0x72, 0x80, 0xd4, 0xc2, 0x5d, 0x95, 0x2b, 0x9c, 0x83, 0x5f, 0x88, 0xc8,
	0x5b, 0x7a, 0xeb, 0x31, 0xf7, 0x0b, 0x81, 0x08, 0x23, 0xf3, 0x51, 0x53, 0xe4, 0x5b, 0x62, 0x67,
	0x5c, 0xc0, 0x58, 0x6a, 0xf5, 0x5e, 0x47, 0x81, 0x85, 0x6e, 0xc1, 0x73, 0x80, 0x7d, 0xd1, 0x18,
	0xaa, 0x6e, 0x85, 0xd0, 0xd1, 0x68, 0xe9, 0xad, 0xff, 0xf3, 0x01, 0x59, 0x3d, 0x41, 0x18, 0x97,
	0x82, 0x93, 0x4c, 0x5b, 0xcd, 0x49, 0x72, 0x7a, 0xc3, 0x0b, 0x98, 0xb8, 0x1e, 0x7b, 0xe0, 0xec,
	0x32, 0xdc, 0x96, 0x8d, 0xdc, 0xf6, 0x35, 0xfc, 0xe7, 0x73, 0xe7, 0x76, 0xd3, 0x83, 0x6b, 0x09,
	0xd6, 0x63, 0x3e, 0x20, 0xab, 0xeb, 0x43, 0x77, 0x62, 0xf2, 0x81, 0xbb, 0xfb, 0xfd, 0x6f, 0xf7,
	0xe6, 0x06, 0x82, 0xb8, 0x14, 0x18, 0xc2, 0x2c, 0x2e, 0x45, 0xb6, 0xbb, 0xcb, 0x12, 0x55, 0x11,
	0xfb, 0x87, 0x08, 0xa7, 0x4e, 0x98, 0xb9, 0x5a, 0xf6, 0x39, 0x3d, 0x64, 0x89, 0xc9, 0xd9, 0xd7,
	0x74, 0xd3, 0xfe, 0xde, 0x5f, 0x57, 0x82, 0x67, 0x10, 0xf6, 0x5b, 0xaf, 0x9a, 0x0f, 0x20, 0x27,
	0xc9, 0x3c, 0x5c, 0x00, 0x1b, 0xb0, 0x47, 0xa5, 0xf7, 0x82, 0xf9, 0x47, 0x34, 0x7d, 0xa5, 0x8a,
	0x58, 0x70, 0x24, 0xbd, 0x7f, 0x36, 0xc4, 0x46, 0x2f, 0x13, 0xfb, 0x86, 0x57, 0xdf, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x5e, 0xc4, 0x5d, 0x33, 0xd3, 0x01, 0x00, 0x00,
}
