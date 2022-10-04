gen_proto: 
	protoc -I ./protobuf --go_out=require_unimplemented_servers=false:. --go-grpc_out=. protobuf/auctions_service.proto

