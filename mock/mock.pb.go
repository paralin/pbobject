// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/aperturerobotics/pbobject/mock/mock.proto

/*
Package mock is a generated protocol buffer package.

It is generated from these files:
	github.com/aperturerobotics/pbobject/mock/mock.proto

It has these top-level messages:
	MockObject
*/
package mock

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

// MockObject is a mock object type.
type MockObject struct {
	// Message is a mock data payload.
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *MockObject) Reset()                    { *m = MockObject{} }
func (m *MockObject) String() string            { return proto.CompactTextString(m) }
func (*MockObject) ProtoMessage()               {}
func (*MockObject) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MockObject) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*MockObject)(nil), "mock.MockObject")
}

func init() {
	proto.RegisterFile("github.com/aperturerobotics/pbobject/mock/mock.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x49, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0x2c, 0x48, 0x2d, 0x2a, 0x29, 0x2d, 0x4a, 0x2d,
	0xca, 0x4f, 0xca, 0x2f, 0xc9, 0x4c, 0x2e, 0xd6, 0x2f, 0x48, 0xca, 0x4f, 0xca, 0x4a, 0x4d, 0x2e,
	0xd1, 0xcf, 0xcd, 0x4f, 0xce, 0x06, 0x13, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x2c, 0x20,
	0xb6, 0x92, 0x1a, 0x17, 0x97, 0x6f, 0x7e, 0x72, 0xb6, 0x3f, 0x58, 0x8d, 0x90, 0x04, 0x17, 0x7b,
	0x6e, 0x6a, 0x71, 0x71, 0x62, 0x7a, 0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c, 0x9b,
	0xc4, 0x06, 0xd6, 0x64, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xd3, 0xbf, 0xd0, 0xb7, 0x6c, 0x00,
	0x00, 0x00,
}
