proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative services/companyInfoGetter/companyInfoGetter.proto
	protoc -I . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true services/companyInfoGetter/companyInfoGetter.proto
.PHONY: proto