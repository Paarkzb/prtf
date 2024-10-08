generate: 
	protoc -I proto proto/sso/sso.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

generate-auth: 
	protoc -I proto proto/sso/sso.proto --go_out=../auth/sso/protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=../auth/sso/protos/gen/go/ --go-grpc_opt=paths=source_relative

generate-gateway: 
	protoc -I proto proto/sso/sso.proto --go_out=../gateway/protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=../gateway/protos/gen/go/ --go-grpc_opt=paths=source_relative

generate-services:
	protoc -I proto proto/sso/sso.proto --go_out=../auth/sso/protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=../auth/sso/protos/gen/go/ --go-grpc_opt=paths=source_relative
	protoc -I proto proto/sso/sso.proto --go_out=../gateway/protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=../gateway/protos/gen/go/ --go-grpc_opt=paths=source_relative
