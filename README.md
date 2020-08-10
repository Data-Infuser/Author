# Data Infuser 인증 서비스

## 개발환경
* Golang 1.14.4
  * grpc-go (https://github.com/grpc/grpc-go)
  * GORM (go orm library, https://github.com/jinzhu/gorm)
* MySQL 5.7

## Configuration

* config/database-sample.yaml 참고하여 DB Config 생성
  * config/dev/database.yaml 또는 config/prod/database.yaml

## Proto Buffer 공통 모듈 다운로드

```sh
$ git clone git@gitlab.com:promptech1/data-infuser/infuser-protobuf.git
```

## 개발 환경 실행
* gRPC Server
```sh
go run main.go -logtostderr=true
```