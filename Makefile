all:
	protoc -I. \
		-I${GOPATH}/src \
		-I=third_party/googleapis \
		--go_out=plugins=grpc:gen \
		proto/*.proto
