// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tx.proto

/*
Package bchain is a generated protocol buffer package.

It is generated from these files:
	tx.proto

It has these top-level messages:
	ProtoTransaction
*/
package bchain

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

type ProtoTransaction struct {
	Txid      []byte                       `protobuf:"bytes,1,opt,name=Txid,json=txid,proto3" json:"Txid,omitempty"`
	Hex       []byte                       `protobuf:"bytes,2,opt,name=Hex,json=hex,proto3" json:"Hex,omitempty"`
	Blocktime uint64                       `protobuf:"varint,3,opt,name=Blocktime,json=blocktime" json:"Blocktime,omitempty"`
	Locktime  uint32                       `protobuf:"varint,4,opt,name=Locktime,json=locktime" json:"Locktime,omitempty"`
	Height    uint32                       `protobuf:"varint,5,opt,name=Height,json=height" json:"Height,omitempty"`
	Vin       []*ProtoTransaction_VinType  `protobuf:"bytes,6,rep,name=Vin,json=vin" json:"Vin,omitempty"`
	Vout      []*ProtoTransaction_VoutType `protobuf:"bytes,7,rep,name=Vout,json=vout" json:"Vout,omitempty"`
	Version   int32                        `protobuf:"varint,8,opt,name=Version,json=version" json:"Version,omitempty"`
	TxType    uint32                       `protobuf:"varint,9,opt,name=TxType,json=txType" json:"TxType,omitempty"`
}

func (m *ProtoTransaction) Reset()                    { *m = ProtoTransaction{} }
func (m *ProtoTransaction) String() string            { return proto.CompactTextString(m) }
func (*ProtoTransaction) ProtoMessage()               {}
func (*ProtoTransaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoTransaction) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *ProtoTransaction) GetHex() []byte {
	if m != nil {
		return m.Hex
	}
	return nil
}

func (m *ProtoTransaction) GetBlocktime() uint64 {
	if m != nil {
		return m.Blocktime
	}
	return 0
}

func (m *ProtoTransaction) GetLocktime() uint32 {
	if m != nil {
		return m.Locktime
	}
	return 0
}

func (m *ProtoTransaction) GetHeight() uint32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ProtoTransaction) GetVin() []*ProtoTransaction_VinType {
	if m != nil {
		return m.Vin
	}
	return nil
}

func (m *ProtoTransaction) GetVout() []*ProtoTransaction_VoutType {
	if m != nil {
		return m.Vout
	}
	return nil
}

func (m *ProtoTransaction) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ProtoTransaction) GetTxType() uint32 {
	if m != nil {
		return m.TxType
	}
	return 0
}

type ProtoTransaction_VinType struct {
	Coinbase     string   `protobuf:"bytes,1,opt,name=Coinbase,json=coinbase" json:"Coinbase,omitempty"`
	Txid         []byte   `protobuf:"bytes,2,opt,name=Txid,json=txid,proto3" json:"Txid,omitempty"`
	Vout         uint32   `protobuf:"varint,3,opt,name=Vout,json=vout" json:"Vout,omitempty"`
	ScriptSigHex []byte   `protobuf:"bytes,4,opt,name=ScriptSigHex,json=scriptSigHex,proto3" json:"ScriptSigHex,omitempty"`
	Sequence     uint32   `protobuf:"varint,5,opt,name=Sequence,json=sequence" json:"Sequence,omitempty"`
	Addresses    []string `protobuf:"bytes,6,rep,name=Addresses,json=addresses" json:"Addresses,omitempty"`
}

func (m *ProtoTransaction_VinType) Reset()                    { *m = ProtoTransaction_VinType{} }
func (m *ProtoTransaction_VinType) String() string            { return proto.CompactTextString(m) }
func (*ProtoTransaction_VinType) ProtoMessage()               {}
func (*ProtoTransaction_VinType) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *ProtoTransaction_VinType) GetCoinbase() string {
	if m != nil {
		return m.Coinbase
	}
	return ""
}

func (m *ProtoTransaction_VinType) GetTxid() []byte {
	if m != nil {
		return m.Txid
	}
	return nil
}

func (m *ProtoTransaction_VinType) GetVout() uint32 {
	if m != nil {
		return m.Vout
	}
	return 0
}

func (m *ProtoTransaction_VinType) GetScriptSigHex() []byte {
	if m != nil {
		return m.ScriptSigHex
	}
	return nil
}

func (m *ProtoTransaction_VinType) GetSequence() uint32 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *ProtoTransaction_VinType) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type ProtoTransaction_VoutType struct {
	ValueSat        []byte   `protobuf:"bytes,1,opt,name=ValueSat,json=valueSat,proto3" json:"ValueSat,omitempty"`
	N               uint32   `protobuf:"varint,2,opt,name=N,json=n" json:"N,omitempty"`
	ScriptPubKeyHex []byte   `protobuf:"bytes,3,opt,name=ScriptPubKeyHex,json=scriptPubKeyHex,proto3" json:"ScriptPubKeyHex,omitempty"`
	ScriptType      string   `protobuf:"bytes,4,opt,name=ScriptType,json=scriptType" json:"ScriptType,omitempty"`
	Addresses       []string `protobuf:"bytes,5,rep,name=Addresses,json=addresses" json:"Addresses,omitempty"`
}

func (m *ProtoTransaction_VoutType) Reset()                    { *m = ProtoTransaction_VoutType{} }
func (m *ProtoTransaction_VoutType) String() string            { return proto.CompactTextString(m) }
func (*ProtoTransaction_VoutType) ProtoMessage()               {}
func (*ProtoTransaction_VoutType) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *ProtoTransaction_VoutType) GetValueSat() []byte {
	if m != nil {
		return m.ValueSat
	}
	return nil
}

func (m *ProtoTransaction_VoutType) GetN() uint32 {
	if m != nil {
		return m.N
	}
	return 0
}

func (m *ProtoTransaction_VoutType) GetScriptPubKeyHex() []byte {
	if m != nil {
		return m.ScriptPubKeyHex
	}
	return nil
}

func (m *ProtoTransaction_VoutType) GetScriptType() string {
	if m != nil {
		return m.ScriptType
	}
	return ""
}

func (m *ProtoTransaction_VoutType) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

func init() {
	proto.RegisterType((*ProtoTransaction)(nil), "bchain.ProtoTransaction")
	proto.RegisterType((*ProtoTransaction_VinType)(nil), "bchain.ProtoTransaction.VinType")
	proto.RegisterType((*ProtoTransaction_VoutType)(nil), "bchain.ProtoTransaction.VoutType")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x6a, 0xdc, 0x30,
	0x18, 0xc5, 0x51, 0xa4, 0xf1, 0xc8, 0x5f, 0x1d, 0x12, 0xb4, 0x28, 0x62, 0x28, 0xc5, 0xcd, 0xca,
	0x2b, 0x2f, 0x52, 0x7a, 0x80, 0xb6, 0x9b, 0x40, 0x4b, 0x09, 0xb2, 0xf1, 0xde, 0x7f, 0x44, 0x2c,
	0x9a, 0x4a, 0x53, 0x4b, 0x36, 0xce, 0x5d, 0x7a, 0x83, 0x5e, 0xaf, 0x07, 0x28, 0x92, 0xed, 0x21,
	0x19, 0xc8, 0xf2, 0xbd, 0xef, 0x7b, 0xd2, 0x4f, 0xcf, 0x06, 0xea, 0xe6, 0xfc, 0x38, 0x18, 0x67,
	0x58, 0xd4, 0xb4, 0x7d, 0xad, 0xf4, 0xcd, 0x3f, 0x02, 0xd7, 0xf7, 0xde, 0x29, 0x87, 0x5a, 0xdb,
	0xba, 0x75, 0xca, 0x68, 0xc6, 0x80, 0x94, 0xb3, 0xea, 0x38, 0x4a, 0x51, 0x96, 0x08, 0xe2, 0x66,
	0xd5, 0xb1, 0x6b, 0xc0, 0x77, 0x72, 0xe6, 0x17, 0xc1, 0xc2, 0xbd, 0x9c, 0xd9, 0x3b, 0x88, 0xbf,
	0x3c, 0x9a, 0xf6, 0xa7, 0x53, 0xbf, 0x24, 0xc7, 0x29, 0xca, 0x88, 0x88, 0x9b, 0xcd, 0x60, 0x07,
	0xa0, 0xdf, 0xb7, 0x21, 0x49, 0x51, 0x76, 0x29, 0xe8, 0x69, 0xf6, 0x16, 0xa2, 0x3b, 0xa9, 0x1e,
	0x7a, 0xc7, 0x77, 0x61, 0x12, 0xf5, 0x41, 0xb1, 0x5b, 0xc0, 0x95, 0xd2, 0x3c, 0x4a, 0x71, 0xf6,
	0xe6, 0x36, 0xcd, 0x17, 0xc4, 0xfc, 0x1c, 0x2f, 0xaf, 0x94, 0x2e, 0x9f, 0x8e, 0x52, 0xe0, 0x49,
	0x69, 0xf6, 0x09, 0x48, 0x65, 0x46, 0xc7, 0xf7, 0x21, 0xf4, 0xe1, 0xf5, 0x90, 0x19, 0x5d, 0x48,
	0x91, 0xc9, 0x8c, 0x8e, 0x71, 0xd8, 0x57, 0x72, 0xb0, 0xca, 0x68, 0x4e, 0x53, 0x94, 0xed, 0xc4,
	0x7e, 0x5a, 0xa4, 0x87, 0x2b, 0x67, 0xbf, 0xc9, 0xe3, 0x05, 0xce, 0x05, 0x75, 0xf8, 0x8b, 0x60,
	0xbf, 0xde, 0xec, 0x1f, 0xf7, 0xd5, 0x28, 0xdd, 0xd4, 0x56, 0x86, 0x92, 0x62, 0x41, 0xdb, 0x55,
	0x9f, 0xca, 0xbb, 0x78, 0x56, 0x1e, 0x5b, 0x21, 0x71, 0x38, 0x71, 0x21, 0xb8, 0x81, 0xa4, 0x68,
	0x07, 0x75, 0x74, 0x85, 0x7a, 0xf0, 0xcd, 0x92, 0xb0, 0x9f, 0xd8, 0x67, 0x9e, 0xbf, 0xa7, 0x90,
	0xbf, 0x47, 0xa9, 0x5b, 0xb9, 0x56, 0x45, 0xed, 0xaa, 0x7d, 0xfd, 0x9f, 0xbb, 0x6e, 0x90, 0xd6,
	0x4a, 0x1b, 0x2a, 0x8b, 0x45, 0x5c, 0x6f, 0xc6, 0xe1, 0x0f, 0x02, 0xba, 0x3d, 0xd9, 0x1f, 0x53,
	0xd5, 0x8f, 0xa3, 0x2c, 0x6a, 0xb7, 0x7e, 0x53, 0x3a, 0xad, 0x9a, 0x25, 0x80, 0x7e, 0x04, 0xd6,
	0x4b, 0x81, 0x34, 0xcb, 0xe0, 0x6a, 0x81, 0xba, 0x1f, 0x9b, 0x6f, 0xf2, 0xc9, 0x73, 0xe1, 0x10,
	0xb8, 0xb2, 0x2f, 0x6d, 0xf6, 0x1e, 0x60, 0xd9, 0x0c, 0x55, 0x91, 0x50, 0x02, 0xd8, 0x93, 0xf3,
	0x12, 0x6f, 0x77, 0x86, 0xd7, 0x44, 0xe1, 0x2f, 0xfc, 0xf8, 0x3f, 0x00, 0x00, 0xff, 0xff, 0x99,
	0xbb, 0x72, 0x25, 0x91, 0x02, 0x00, 0x00,
}
