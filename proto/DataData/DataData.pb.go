// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: proto/DataData/DataData.proto

package DataData

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
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

type DataDataMigrateReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyValues []byte `protobuf:"bytes,1,opt,name=KeyValues,proto3" json:"KeyValues,omitempty"`
}

func (x *DataDataMigrateReq) Reset() {
	*x = DataDataMigrateReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataMigrateReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataMigrateReq) ProtoMessage() {}

func (x *DataDataMigrateReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataMigrateReq.ProtoReflect.Descriptor instead.
func (*DataDataMigrateReq) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{0}
}

func (x *DataDataMigrateReq) GetKeyValues() []byte {
	if x != nil {
		return x.KeyValues
	}
	return nil
}

type DataDataMigrateResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DataDataMigrateResp) Reset() {
	*x = DataDataMigrateResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataMigrateResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataMigrateResp) ProtoMessage() {}

func (x *DataDataMigrateResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataMigrateResp.ProtoReflect.Descriptor instead.
func (*DataDataMigrateResp) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{1}
}

func (x *DataDataMigrateResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DataDataSyncAllReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyValues []byte `protobuf:"bytes,1,opt,name=KeyValues,proto3" json:"KeyValues,omitempty"`
}

func (x *DataDataSyncAllReq) Reset() {
	*x = DataDataSyncAllReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncAllReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncAllReq) ProtoMessage() {}

func (x *DataDataSyncAllReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncAllReq.ProtoReflect.Descriptor instead.
func (*DataDataSyncAllReq) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{2}
}

func (x *DataDataSyncAllReq) GetKeyValues() []byte {
	if x != nil {
		return x.KeyValues
	}
	return nil
}

type DataDataSyncAllResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DataDataSyncAllResp) Reset() {
	*x = DataDataSyncAllResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncAllResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncAllResp) ProtoMessage() {}

func (x *DataDataSyncAllResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncAllResp.ProtoReflect.Descriptor instead.
func (*DataDataSyncAllResp) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{3}
}

func (x *DataDataSyncAllResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DataDataSyncPutReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *DataDataSyncPutReq) Reset() {
	*x = DataDataSyncPutReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncPutReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncPutReq) ProtoMessage() {}

func (x *DataDataSyncPutReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncPutReq.ProtoReflect.Descriptor instead.
func (*DataDataSyncPutReq) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{4}
}

func (x *DataDataSyncPutReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DataDataSyncPutReq) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type DataDataSyncPutResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DataDataSyncPutResp) Reset() {
	*x = DataDataSyncPutResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncPutResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncPutResp) ProtoMessage() {}

func (x *DataDataSyncPutResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncPutResp.ProtoReflect.Descriptor instead.
func (*DataDataSyncPutResp) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{5}
}

func (x *DataDataSyncPutResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DataDataSyncDeleteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *DataDataSyncDeleteReq) Reset() {
	*x = DataDataSyncDeleteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncDeleteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncDeleteReq) ProtoMessage() {}

func (x *DataDataSyncDeleteReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncDeleteReq.ProtoReflect.Descriptor instead.
func (*DataDataSyncDeleteReq) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{6}
}

func (x *DataDataSyncDeleteReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DataDataSyncDeleteResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *DataDataSyncDeleteResp) Reset() {
	*x = DataDataSyncDeleteResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_DataData_DataData_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataDataSyncDeleteResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataDataSyncDeleteResp) ProtoMessage() {}

func (x *DataDataSyncDeleteResp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_DataData_DataData_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataDataSyncDeleteResp.ProtoReflect.Descriptor instead.
func (*DataDataSyncDeleteResp) Descriptor() ([]byte, []int) {
	return file_proto_DataData_DataData_proto_rawDescGZIP(), []int{7}
}

func (x *DataDataSyncDeleteResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_DataData_DataData_proto protoreflect.FileDescriptor

var file_proto_DataData_DataData_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x2f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x22, 0x32, 0x0a, 0x12, 0x44, 0x61, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x1c, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x2f, 0x0a,
	0x13, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x32,
	0x0a, 0x12, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c,
	0x6c, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x22, 0x2f, 0x0a, 0x13, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79,
	0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x3c, 0x0a, 0x12, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x79, 0x6e, 0x63, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x2f, 0x0a, 0x13, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e,
	0x63, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x29, 0x0a, 0x15, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79,
	0x6e, 0x63, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x4b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x32, 0x0a,
	0x16, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x32, 0xdb, 0x02, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x50,
	0x0a, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74,
	0x65, 0x12, 0x1c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a,
	0x1d, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44,
	0x61, 0x74, 0x61, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00,
	0x12, 0x50, 0x0a, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63,
	0x41, 0x6c, 0x6c, 0x12, 0x1c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x52, 0x65,
	0x71, 0x1a, 0x1d, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70,
	0x22, 0x00, 0x12, 0x50, 0x0a, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79,
	0x6e, 0x63, 0x50, 0x75, 0x74, 0x12, 0x1c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x74,
	0x52, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79, 0x6e, 0x63, 0x50, 0x75, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x12, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x79, 0x6e, 0x63, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53, 0x79,
	0x6e, 0x63, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x79, 0x6e, 0x63, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_DataData_DataData_proto_rawDescOnce sync.Once
	file_proto_DataData_DataData_proto_rawDescData = file_proto_DataData_DataData_proto_rawDesc
)

func file_proto_DataData_DataData_proto_rawDescGZIP() []byte {
	file_proto_DataData_DataData_proto_rawDescOnce.Do(func() {
		file_proto_DataData_DataData_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_DataData_DataData_proto_rawDescData)
	})
	return file_proto_DataData_DataData_proto_rawDescData
}

var file_proto_DataData_DataData_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_DataData_DataData_proto_goTypes = []interface{}{
	(*DataDataMigrateReq)(nil),     // 0: DataData.DataDataMigrateReq
	(*DataDataMigrateResp)(nil),    // 1: DataData.DataDataMigrateResp
	(*DataDataSyncAllReq)(nil),     // 2: DataData.DataDataSyncAllReq
	(*DataDataSyncAllResp)(nil),    // 3: DataData.DataDataSyncAllResp
	(*DataDataSyncPutReq)(nil),     // 4: DataData.DataDataSyncPutReq
	(*DataDataSyncPutResp)(nil),    // 5: DataData.DataDataSyncPutResp
	(*DataDataSyncDeleteReq)(nil),  // 6: DataData.DataDataSyncDeleteReq
	(*DataDataSyncDeleteResp)(nil), // 7: DataData.DataDataSyncDeleteResp
}
var file_proto_DataData_DataData_proto_depIdxs = []int32{
	0, // 0: DataData.DataData.DataDataMigrate:input_type -> DataData.DataDataMigrateReq
	2, // 1: DataData.DataData.DataDataSyncAll:input_type -> DataData.DataDataSyncAllReq
	4, // 2: DataData.DataData.DataDataSyncPut:input_type -> DataData.DataDataSyncPutReq
	6, // 3: DataData.DataData.DataDataSyncDelete:input_type -> DataData.DataDataSyncDeleteReq
	1, // 4: DataData.DataData.DataDataMigrate:output_type -> DataData.DataDataMigrateResp
	3, // 5: DataData.DataData.DataDataSyncAll:output_type -> DataData.DataDataSyncAllResp
	5, // 6: DataData.DataData.DataDataSyncPut:output_type -> DataData.DataDataSyncPutResp
	7, // 7: DataData.DataData.DataDataSyncDelete:output_type -> DataData.DataDataSyncDeleteResp
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_DataData_DataData_proto_init() }
func file_proto_DataData_DataData_proto_init() {
	if File_proto_DataData_DataData_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_DataData_DataData_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataMigrateReq); i {
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
		file_proto_DataData_DataData_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataMigrateResp); i {
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
		file_proto_DataData_DataData_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncAllReq); i {
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
		file_proto_DataData_DataData_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncAllResp); i {
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
		file_proto_DataData_DataData_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncPutReq); i {
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
		file_proto_DataData_DataData_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncPutResp); i {
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
		file_proto_DataData_DataData_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncDeleteReq); i {
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
		file_proto_DataData_DataData_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataDataSyncDeleteResp); i {
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
			RawDescriptor: file_proto_DataData_DataData_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_DataData_DataData_proto_goTypes,
		DependencyIndexes: file_proto_DataData_DataData_proto_depIdxs,
		MessageInfos:      file_proto_DataData_DataData_proto_msgTypes,
	}.Build()
	File_proto_DataData_DataData_proto = out.File
	file_proto_DataData_DataData_proto_rawDesc = nil
	file_proto_DataData_DataData_proto_goTypes = nil
	file_proto_DataData_DataData_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DataDataClient is the client API for DataData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataDataClient interface {
	DataDataMigrate(ctx context.Context, in *DataDataMigrateReq, opts ...grpc.CallOption) (*DataDataMigrateResp, error)
	DataDataSyncAll(ctx context.Context, in *DataDataSyncAllReq, opts ...grpc.CallOption) (*DataDataSyncAllResp, error)
	DataDataSyncPut(ctx context.Context, in *DataDataSyncPutReq, opts ...grpc.CallOption) (*DataDataSyncPutResp, error)
	DataDataSyncDelete(ctx context.Context, in *DataDataSyncDeleteReq, opts ...grpc.CallOption) (*DataDataSyncDeleteResp, error)
}

type dataDataClient struct {
	cc grpc.ClientConnInterface
}

func NewDataDataClient(cc grpc.ClientConnInterface) DataDataClient {
	return &dataDataClient{cc}
}

func (c *dataDataClient) DataDataMigrate(ctx context.Context, in *DataDataMigrateReq, opts ...grpc.CallOption) (*DataDataMigrateResp, error) {
	out := new(DataDataMigrateResp)
	err := c.cc.Invoke(ctx, "/DataData.DataData/DataDataMigrate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataDataClient) DataDataSyncAll(ctx context.Context, in *DataDataSyncAllReq, opts ...grpc.CallOption) (*DataDataSyncAllResp, error) {
	out := new(DataDataSyncAllResp)
	err := c.cc.Invoke(ctx, "/DataData.DataData/DataDataSyncAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataDataClient) DataDataSyncPut(ctx context.Context, in *DataDataSyncPutReq, opts ...grpc.CallOption) (*DataDataSyncPutResp, error) {
	out := new(DataDataSyncPutResp)
	err := c.cc.Invoke(ctx, "/DataData.DataData/DataDataSyncPut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataDataClient) DataDataSyncDelete(ctx context.Context, in *DataDataSyncDeleteReq, opts ...grpc.CallOption) (*DataDataSyncDeleteResp, error) {
	out := new(DataDataSyncDeleteResp)
	err := c.cc.Invoke(ctx, "/DataData.DataData/DataDataSyncDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataDataServer is the server API for DataData service.
type DataDataServer interface {
	DataDataMigrate(context.Context, *DataDataMigrateReq) (*DataDataMigrateResp, error)
	DataDataSyncAll(context.Context, *DataDataSyncAllReq) (*DataDataSyncAllResp, error)
	DataDataSyncPut(context.Context, *DataDataSyncPutReq) (*DataDataSyncPutResp, error)
	DataDataSyncDelete(context.Context, *DataDataSyncDeleteReq) (*DataDataSyncDeleteResp, error)
}

// UnimplementedDataDataServer can be embedded to have forward compatible implementations.
type UnimplementedDataDataServer struct {
}

func (*UnimplementedDataDataServer) DataDataMigrate(context.Context, *DataDataMigrateReq) (*DataDataMigrateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataDataMigrate not implemented")
}
func (*UnimplementedDataDataServer) DataDataSyncAll(context.Context, *DataDataSyncAllReq) (*DataDataSyncAllResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataDataSyncAll not implemented")
}
func (*UnimplementedDataDataServer) DataDataSyncPut(context.Context, *DataDataSyncPutReq) (*DataDataSyncPutResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataDataSyncPut not implemented")
}
func (*UnimplementedDataDataServer) DataDataSyncDelete(context.Context, *DataDataSyncDeleteReq) (*DataDataSyncDeleteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DataDataSyncDelete not implemented")
}

func RegisterDataDataServer(s *grpc.Server, srv DataDataServer) {
	s.RegisterService(&_DataData_serviceDesc, srv)
}

func _DataData_DataDataMigrate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataDataMigrateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataDataServer).DataDataMigrate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DataData.DataData/DataDataMigrate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataDataServer).DataDataMigrate(ctx, req.(*DataDataMigrateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataData_DataDataSyncAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataDataSyncAllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataDataServer).DataDataSyncAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DataData.DataData/DataDataSyncAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataDataServer).DataDataSyncAll(ctx, req.(*DataDataSyncAllReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataData_DataDataSyncPut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataDataSyncPutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataDataServer).DataDataSyncPut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DataData.DataData/DataDataSyncPut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataDataServer).DataDataSyncPut(ctx, req.(*DataDataSyncPutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataData_DataDataSyncDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DataDataSyncDeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataDataServer).DataDataSyncDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/DataData.DataData/DataDataSyncDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataDataServer).DataDataSyncDelete(ctx, req.(*DataDataSyncDeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataData_serviceDesc = grpc.ServiceDesc{
	ServiceName: "DataData.DataData",
	HandlerType: (*DataDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DataDataMigrate",
			Handler:    _DataData_DataDataMigrate_Handler,
		},
		{
			MethodName: "DataDataSyncAll",
			Handler:    _DataData_DataDataSyncAll_Handler,
		},
		{
			MethodName: "DataDataSyncPut",
			Handler:    _DataData_DataDataSyncPut_Handler,
		},
		{
			MethodName: "DataDataSyncDelete",
			Handler:    _DataData_DataDataSyncDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/DataData/DataData.proto",
}
