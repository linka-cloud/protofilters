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
)

func main() {
	// gen.Gen()
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
