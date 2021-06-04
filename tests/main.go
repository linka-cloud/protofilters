/*
 Copyright 2021 Linka Cloud  All rights reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"

	ff "go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/matcher"
	"go.linka.cloud/protofilters/tests/gen"
)

func main() {
	gen.Gen()
	b, err := ioutil.ReadFile("test.file-descriptor.bin")
	if err != nil {
		logrus.Fatal(err)
	}
	var fdp descriptorpb.FileDescriptorProto
	if err := proto.Unmarshal(b, &fdp); err != nil {
		logrus.Fatal(err)
	}
	fd, err := protodesc.NewFile(&fdp, protoregistry.GlobalFiles)
	if err != nil {
		logrus.Fatal(err)
	}
	// protodesc.NewFile(d.)
	d := dynamicpb.NewMessage(fd.Messages().Get(0))
	b, err = ioutil.ReadFile("test.bin")
	if err != nil {
		logrus.Fatal(err)
	}
	a := anypb.Any{}
	if err := proto.Unmarshal(b, &a); err != nil {
		logrus.Fatal(err)
	}
	if err := anypb.UnmarshalTo(&a, d, proto.UnmarshalOptions{}); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("%+v", d)
	logrus.Info(d.ProtoReflect().Descriptor().FullName())
	filter := &ff.FieldFilter{
		Field: "string_field",
		// Filter: ff.StringRegex(".*....*"),
		Filter: ff.StringIN("", "whatever..."),
	}
	ok, err := matcher.MatchFilters(d, filter)
	if err != nil {
		logrus.Fatal(err)
	}
	if !ok {
		logrus.Fatal("reg ex not found")
	}
	logrus.Infof("%v match %v", d, filter)
}
