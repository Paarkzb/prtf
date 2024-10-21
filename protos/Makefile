generate: 
	protoc --proto_path proto/sso --proto_path ./ --go_out=../auth/sso/protos/gen/go/sso --go_opt=paths=source_relative --go-grpc_out=../auth/sso/protos/gen/go/sso --go-grpc_opt=paths=source_relative --grpc-gateway_out=../auth/sso/protos/gen/go/sso --grpc-gateway_opt=paths=source_relative proto/sso/sso.proto
