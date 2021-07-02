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
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/anypb"

	ff "go.linka.cloud/protofilters"
	"go.linka.cloud/protofilters/matcher"
)

type registry struct {
	m  map[pref.FullName]pref.MessageDescriptor
	mu sync.RWMutex
}

func (r *registry) Import(b []byte, allowUnResolved ...bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.m == nil {
		r.m = make(map[pref.FullName]pref.MessageDescriptor)
	}
	var fdp descriptorpb.FileDescriptorProto
	if err := proto.Unmarshal(b, &fdp); err != nil {
		logrus.Fatal(err)
	}
	opts := &protodesc.FileOptions{
		AllowUnresolvable: len(allowUnResolved) > 0 && allowUnResolved[0],
	}
	fd, err := opts.New(&fdp, protoregistry.GlobalFiles)
	if err != nil {
		return err
	}
	msgs := fd.Messages()
	for i := 0; i < msgs.Len(); i++ {
		m := msgs.Get(i)
		r.m[m.FullName()] = m
	}
	return nil
}

func (r *registry) UnmarshalAny(a *anypb.Any) (*dynamicpb.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.m == nil {
		return nil, protoregistry.NotFound
	}
	d, ok := r.m[a.MessageName()]
	if !ok {
		return nil, protoregistry.NotFound
	}
	m := dynamicpb.NewMessage(d)
	if err := anypb.UnmarshalTo(a, m, proto.UnmarshalOptions{}); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *registry) Unmarshal(b []byte) (*dynamicpb.Message, error) {
	var a anypb.Any
	if err := proto.Unmarshal(b, &a); err != nil {
		return nil, err
	}
	return r.UnmarshalAny(&a)
}

func main() {
	b, err := ioutil.ReadFile("test.file-descriptor.bin")
	if err != nil {
		logrus.Fatal(err)
	}
	reg := &registry{}
	if err := reg.Import(b, true); err != nil {
		logrus.Fatal(err)
	}
	b, err = ioutil.ReadFile("test.bin")
	if err != nil {
		logrus.Fatal(err)
	}
	d, err := reg.Unmarshal(b)
	if err != nil {
		logrus.Fatal(err)
	}
	d.Range(func(d pref.FieldDescriptor, v pref.Value) bool {
		i := v.Interface()
		_ = i

		return true
	})
	filter := &ff.FieldFilter{
		Field:  "string_value_field",
		Filter: ff.StringIN("", "..."),
	}
	ok, err := matcher.MatchFilters(d, filter)
	if err != nil {
		logrus.Fatal(err)
	}
	if !ok {
		logrus.Fatal("not found")
	}
	logrus.Infof("%v match %v", d, filter)
}
