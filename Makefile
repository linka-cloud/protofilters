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

SHELL := $(CURDIR)/bin/mise
.SHELLFLAGS := exec -- /bin/sh -c

.PHONY: bin
bin:
	@bin/mise install --locked

clean:
	@find $(PWD) -name '*.pb*.go' -type f -exec rm {} \;

.PHONY: proto
proto: gen-proto lint

.PHONY: gen-proto
gen-proto: bin
	@buf generate --disable-symlinks

.PHONY: lint
lint:
	@goimports -w -local $(MODULE) $$(find . -type d -name '.mise' -prune -o -type f -name '*.go' -print)
	@go fmt $(PWD)

.PHONY: tests
tests: proto
	@go test -v ./...
