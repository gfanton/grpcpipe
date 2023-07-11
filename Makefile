generate: rpcmanager/rpcmanager.pb.go
rpcmanager/rpcmanager.pb.go: rpcmanager/rpcmanager.proto
	protoc -I rpcmanager/ --go_opt=paths=source_relative --go_out=plugins=grpc:rpcmanager/ rpcmanager/rpcmanager.proto
