gen_proto: 
	protoc -I ./protobuf --go_out=. --go-grpc_out=. protobuf/auctions_service.proto