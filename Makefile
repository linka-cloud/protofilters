


.PHONY: proto
proto:
	@protoc -I. --go_out=paths=source_relative:. ./*.proto


.PHONY: proto-test
proto-test:
	@protoc -I. --go_out=paths=source_relative:. tests/pb/*.proto
