# Data Infuser 인증 서비스

## 개발환경
* Golang 1.14.4
  * grpc-go (https://github.com/grpc/grpc-go)
  * GORM (go orm library, https://github.com/jinzhu/gorm)
* MySQL 5.7

## Configuration

* config/database-sample.yaml 참고하여 config/database.yaml 생성

## protoc 수행

```sh
$ make
```

## 개발 환경 실행
* gRPC Server
```sh
go run main.go -logtostderr=true
```

* gRPC Client(server 통신 확인)
```sh
go run grpc_client/main.go -logtostderr=true
```