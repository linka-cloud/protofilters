# Copyright 2021 Linka Cloud  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODULE = go.linka.cloud/protofilters

$(shell mkdir -p .bin)

export GOBIN=$(PWD)/.bin

export PATH := $(GOBIN):$(PATH)

PROTO_BASE_PATH = $(PWD)

INCLUDE_PROTO_PATH = -I$(PROTO_BASE_PATH) \
	-I $(shell go list -m -f {{.Dir}} google.golang.org/protobuf)

bin:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto
	@go install golang.org/x/tools/cmd/goimports

clean:
	@rm -rf .bin
	@find $(PROTO_BASE_PATH) -name '*.pb*.go' -type f -exec rm {} \;


.PHONY: proto
proto: gen-proto lint

.PHONY: gen-proto
gen-proto: bin
	@find $(PROTO_BASE_PATH) -name '*.proto' -type f -exec \
    	protoc $(INCLUDE_PROTO_PATH) \
    		--go-vtproto_out=features=marshal+unmarshal+size+clone,paths=source_relative:. \
    		--go_out=paths=source_relative:. {} \;

.PHONY: lint
lint:
	@goimports -w -local $(MODULE) $(PWD)
	@go fmt $(PWD)

.PHONY: tests
tests: proto
	@go test -v ./...
