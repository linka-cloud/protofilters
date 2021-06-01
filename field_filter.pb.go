// Copyright 2021 Linka Cloud  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Linka Cloud nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.1
// source: field_filter.proto

package protofilters

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FieldsFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The set of field mask paths.
	Filters map[string]*Filter `protobuf:"bytes,1,rep,name=filters,proto3" json:"filters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FieldsFilter) Reset() {
	*x = FieldsFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldsFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldsFilter) ProtoMessage() {}

func (x *FieldsFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldsFilter.ProtoReflect.Descriptor instead.
func (*FieldsFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{0}
}

func (x *FieldsFilter) GetFilters() map[string]*Filter {
	if x != nil {
		return x.Filters
	}
	return nil
}

type FieldFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field  string  `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Filter *Filter `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *FieldFilter) Reset() {
	*x = FieldFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldFilter) ProtoMessage() {}

func (x *FieldFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldFilter.ProtoReflect.Descriptor instead.
func (*FieldFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{1}
}

func (x *FieldFilter) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *FieldFilter) GetFilter() *Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Match:
	//	*Filter_String_
	//	*Filter_Number
	//	*Filter_Bool
	//	*Filter_Null
	Match isFilter_Match `protobuf_oneof:"match"`
}

func (x *Filter) Reset() {
	*x = Filter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Filter) ProtoMessage() {}

func (x *Filter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Filter.ProtoReflect.Descriptor instead.
func (*Filter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{2}
}

func (m *Filter) GetMatch() isFilter_Match {
	if m != nil {
		return m.Match
	}
	return nil
}

func (x *Filter) GetString_() *StringFilter {
	if x, ok := x.GetMatch().(*Filter_String_); ok {
		return x.String_
	}
	return nil
}

func (x *Filter) GetNumber() *NumberFilter {
	if x, ok := x.GetMatch().(*Filter_Number); ok {
		return x.Number
	}
	return nil
}

func (x *Filter) GetBool() *BoolFilter {
	if x, ok := x.GetMatch().(*Filter_Bool); ok {
		return x.Bool
	}
	return nil
}

func (x *Filter) GetNull() *NullFilter {
	if x, ok := x.GetMatch().(*Filter_Null); ok {
		return x.Null
	}
	return nil
}

type isFilter_Match interface {
	isFilter_Match()
}

type Filter_String_ struct {
	String_ *StringFilter `protobuf:"bytes,1,opt,name=string,proto3,oneof"`
}

type Filter_Number struct {
	Number *NumberFilter `protobuf:"bytes,2,opt,name=number,proto3,oneof"`
}

type Filter_Bool struct {
	Bool *BoolFilter `protobuf:"bytes,3,opt,name=bool,proto3,oneof"`
}

type Filter_Null struct {
	Null *NullFilter `protobuf:"bytes,4,opt,name=null,proto3,oneof"`
}

func (*Filter_String_) isFilter_Match() {}

func (*Filter_Number) isFilter_Match() {}

func (*Filter_Bool) isFilter_Match() {}

func (*Filter_Null) isFilter_Match() {}

type StringFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Condition:
	//	*StringFilter_Equals
	//	*StringFilter_Regex
	//	*StringFilter_In_
	Condition       isStringFilter_Condition `protobuf_oneof:"condition"`
	Not             bool                     `protobuf:"varint,4,opt,name=not,proto3" json:"not,omitempty"`
	CaseInsensitive bool                     `protobuf:"varint,5,opt,name=case_insensitive,json=caseInsensitive,proto3" json:"case_insensitive,omitempty"`
}

func (x *StringFilter) Reset() {
	*x = StringFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringFilter) ProtoMessage() {}

func (x *StringFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringFilter.ProtoReflect.Descriptor instead.
func (*StringFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{3}
}

func (m *StringFilter) GetCondition() isStringFilter_Condition {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (x *StringFilter) GetEquals() string {
	if x, ok := x.GetCondition().(*StringFilter_Equals); ok {
		return x.Equals
	}
	return ""
}

func (x *StringFilter) GetRegex() string {
	if x, ok := x.GetCondition().(*StringFilter_Regex); ok {
		return x.Regex
	}
	return ""
}

func (x *StringFilter) GetIn() *StringFilter_In {
	if x, ok := x.GetCondition().(*StringFilter_In_); ok {
		return x.In
	}
	return nil
}

func (x *StringFilter) GetNot() bool {
	if x != nil {
		return x.Not
	}
	return false
}

func (x *StringFilter) GetCaseInsensitive() bool {
	if x != nil {
		return x.CaseInsensitive
	}
	return false
}

type isStringFilter_Condition interface {
	isStringFilter_Condition()
}

type StringFilter_Equals struct {
	Equals string `protobuf:"bytes,1,opt,name=equals,proto3,oneof"`
}

type StringFilter_Regex struct {
	Regex string `protobuf:"bytes,2,opt,name=regex,proto3,oneof"`
}

type StringFilter_In_ struct {
	In *StringFilter_In `protobuf:"bytes,3,opt,name=in,proto3,oneof"`
}

func (*StringFilter_Equals) isStringFilter_Condition() {}

func (*StringFilter_Regex) isStringFilter_Condition() {}

func (*StringFilter_In_) isStringFilter_Condition() {}

type NumberFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Condition:
	//	*NumberFilter_Equals
	//	*NumberFilter_In_
	Condition isNumberFilter_Condition `protobuf_oneof:"condition"`
	Not       bool                     `protobuf:"varint,4,opt,name=not,proto3" json:"not,omitempty"`
}

func (x *NumberFilter) Reset() {
	*x = NumberFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NumberFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NumberFilter) ProtoMessage() {}

func (x *NumberFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NumberFilter.ProtoReflect.Descriptor instead.
func (*NumberFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{4}
}

func (m *NumberFilter) GetCondition() isNumberFilter_Condition {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (x *NumberFilter) GetEquals() float64 {
	if x, ok := x.GetCondition().(*NumberFilter_Equals); ok {
		return x.Equals
	}
	return 0
}

func (x *NumberFilter) GetIn() *NumberFilter_In {
	if x, ok := x.GetCondition().(*NumberFilter_In_); ok {
		return x.In
	}
	return nil
}

func (x *NumberFilter) GetNot() bool {
	if x != nil {
		return x.Not
	}
	return false
}

type isNumberFilter_Condition interface {
	isNumberFilter_Condition()
}

type NumberFilter_Equals struct {
	Equals float64 `protobuf:"fixed64,1,opt,name=equals,proto3,oneof"`
}

type NumberFilter_In_ struct {
	In *NumberFilter_In `protobuf:"bytes,2,opt,name=in,proto3,oneof"`
}

func (*NumberFilter_Equals) isNumberFilter_Condition() {}

func (*NumberFilter_In_) isNumberFilter_Condition() {}

type NullFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Not bool `protobuf:"varint,1,opt,name=not,proto3" json:"not,omitempty"`
}

func (x *NullFilter) Reset() {
	*x = NullFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NullFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NullFilter) ProtoMessage() {}

func (x *NullFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NullFilter.ProtoReflect.Descriptor instead.
func (*NullFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{5}
}

func (x *NullFilter) GetNot() bool {
	if x != nil {
		return x.Not
	}
	return false
}

type BoolFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equals bool `protobuf:"varint,1,opt,name=equals,proto3" json:"equals,omitempty"`
}

func (x *BoolFilter) Reset() {
	*x = BoolFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolFilter) ProtoMessage() {}

func (x *BoolFilter) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolFilter.ProtoReflect.Descriptor instead.
func (*BoolFilter) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{6}
}

func (x *BoolFilter) GetEquals() bool {
	if x != nil {
		return x.Equals
	}
	return false
}

type StringFilter_In struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *StringFilter_In) Reset() {
	*x = StringFilter_In{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringFilter_In) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringFilter_In) ProtoMessage() {}

func (x *StringFilter_In) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringFilter_In.ProtoReflect.Descriptor instead.
func (*StringFilter_In) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{3, 0}
}

func (x *StringFilter_In) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type NumberFilter_In struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []float64 `protobuf:"fixed64,1,rep,packed,name=values,proto3" json:"values,omitempty"`
}

func (x *NumberFilter_In) Reset() {
	*x = NumberFilter_In{}
	if protoimpl.UnsafeEnabled {
		mi := &file_field_filter_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NumberFilter_In) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NumberFilter_In) ProtoMessage() {}

func (x *NumberFilter_In) ProtoReflect() protoreflect.Message {
	mi := &file_field_filter_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NumberFilter_In.ProtoReflect.Descriptor instead.
func (*NumberFilter_In) Descriptor() ([]byte, []int) {
	return file_field_filter_proto_rawDescGZIP(), []int{4, 0}
}

func (x *NumberFilter_In) GetValues() []float64 {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_field_filter_proto protoreflect.FileDescriptor

var file_field_filter_proto_rawDesc = []byte{
	0x0a, 0x12, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x22, 0xbb,
	0x01, 0x0a, 0x0c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12,
	0x4d, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x33, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x73, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x5c,
	0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x5d, 0x0a, 0x0b,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x38, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x8d, 0x02, 0x0a, 0x06,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x40, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x73, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x48, 0x00,
	0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x40, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x73, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x48, 0x00, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x04, 0x62, 0x6f,
	0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61,
	0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x73, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x48, 0x00,
	0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x12, 0x3a, 0x0a, 0x04, 0x6e, 0x75, 0x6c, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e,
	0x4e, 0x75, 0x6c, 0x6c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x75,
	0x6c, 0x6c, 0x42, 0x07, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x22, 0xe5, 0x01, 0x0a, 0x0c,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x06,
	0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06,
	0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78, 0x12, 0x3b,
	0x0a, 0x02, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6c, 0x69, 0x6e,
	0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x2e, 0x49, 0x6e, 0x48, 0x00, 0x52, 0x02, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6e,
	0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x6e, 0x6f, 0x74, 0x12, 0x29, 0x0a,
	0x10, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x6e, 0x73, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x63, 0x61, 0x73, 0x65, 0x49, 0x6e, 0x73,
	0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x1a, 0x1c, 0x0a, 0x02, 0x49, 0x6e, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0xa2, 0x01, 0x0a, 0x0c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x12, 0x3b,
	0x0a, 0x02, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6c, 0x69, 0x6e,
	0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x2e, 0x49, 0x6e, 0x48, 0x00, 0x52, 0x02, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6e,
	0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x6e, 0x6f, 0x74, 0x1a, 0x1c, 0x0a,
	0x02, 0x49, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x01, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x0b, 0x0a, 0x09, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x1e, 0x0a, 0x0a, 0x4e, 0x75, 0x6c, 0x6c,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x6f, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x03, 0x6e, 0x6f, 0x74, 0x22, 0x24, 0x0a, 0x0a, 0x42, 0x6f, 0x6f, 0x6c,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x71, 0x75, 0x61, 0x6c, 0x73, 0x42, 0x2f,
	0x50, 0x01, 0x5a, 0x28, 0x67, 0x6f, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x61, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x3b,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0xf8, 0x01, 0x01, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_field_filter_proto_rawDescOnce sync.Once
	file_field_filter_proto_rawDescData = file_field_filter_proto_rawDesc
)

func file_field_filter_proto_rawDescGZIP() []byte {
	file_field_filter_proto_rawDescOnce.Do(func() {
		file_field_filter_proto_rawDescData = protoimpl.X.CompressGZIP(file_field_filter_proto_rawDescData)
	})
	return file_field_filter_proto_rawDescData
}

var file_field_filter_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_field_filter_proto_goTypes = []interface{}{
	(*FieldsFilter)(nil),    // 0: linka.cloud.protofilters.FieldsFilter
	(*FieldFilter)(nil),     // 1: linka.cloud.protofilters.FieldFilter
	(*Filter)(nil),          // 2: linka.cloud.protofilters.Filter
	(*StringFilter)(nil),    // 3: linka.cloud.protofilters.StringFilter
	(*NumberFilter)(nil),    // 4: linka.cloud.protofilters.NumberFilter
	(*NullFilter)(nil),      // 5: linka.cloud.protofilters.NullFilter
	(*BoolFilter)(nil),      // 6: linka.cloud.protofilters.BoolFilter
	nil,                     // 7: linka.cloud.protofilters.FieldsFilter.FiltersEntry
	(*StringFilter_In)(nil), // 8: linka.cloud.protofilters.StringFilter.In
	(*NumberFilter_In)(nil), // 9: linka.cloud.protofilters.NumberFilter.In
}
var file_field_filter_proto_depIdxs = []int32{
	7, // 0: linka.cloud.protofilters.FieldsFilter.filters:type_name -> linka.cloud.protofilters.FieldsFilter.FiltersEntry
	2, // 1: linka.cloud.protofilters.FieldFilter.filter:type_name -> linka.cloud.protofilters.Filter
	3, // 2: linka.cloud.protofilters.Filter.string:type_name -> linka.cloud.protofilters.StringFilter
	4, // 3: linka.cloud.protofilters.Filter.number:type_name -> linka.cloud.protofilters.NumberFilter
	6, // 4: linka.cloud.protofilters.Filter.bool:type_name -> linka.cloud.protofilters.BoolFilter
	5, // 5: linka.cloud.protofilters.Filter.null:type_name -> linka.cloud.protofilters.NullFilter
	8, // 6: linka.cloud.protofilters.StringFilter.in:type_name -> linka.cloud.protofilters.StringFilter.In
	9, // 7: linka.cloud.protofilters.NumberFilter.in:type_name -> linka.cloud.protofilters.NumberFilter.In
	2, // 8: linka.cloud.protofilters.FieldsFilter.FiltersEntry.value:type_name -> linka.cloud.protofilters.Filter
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_field_filter_proto_init() }
func file_field_filter_proto_init() {
	if File_field_filter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_field_filter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldsFilter); i {
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
		file_field_filter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldFilter); i {
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
		file_field_filter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Filter); i {
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
		file_field_filter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringFilter); i {
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
		file_field_filter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NumberFilter); i {
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
		file_field_filter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NullFilter); i {
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
		file_field_filter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolFilter); i {
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
		file_field_filter_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringFilter_In); i {
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
		file_field_filter_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NumberFilter_In); i {
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
	file_field_filter_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Filter_String_)(nil),
		(*Filter_Number)(nil),
		(*Filter_Bool)(nil),
		(*Filter_Null)(nil),
	}
	file_field_filter_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*StringFilter_Equals)(nil),
		(*StringFilter_Regex)(nil),
		(*StringFilter_In_)(nil),
	}
	file_field_filter_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*NumberFilter_Equals)(nil),
		(*NumberFilter_In_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_field_filter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_field_filter_proto_goTypes,
		DependencyIndexes: file_field_filter_proto_depIdxs,
		MessageInfos:      file_field_filter_proto_msgTypes,
	}.Build()
	File_field_filter_proto = out.File
	file_field_filter_proto_rawDesc = nil
	file_field_filter_proto_goTypes = nil
	file_field_filter_proto_depIdxs = nil
}