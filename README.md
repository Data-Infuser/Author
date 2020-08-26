# Data Infuser 인증 서비스

## 개발환경
* Golang 1.14.4
  * grpc-go (https://github.com/grpc/grpc-go)
  * GORM (go orm library, https://github.com/jinzhu/gorm)
* MySQL 5.7
* Redis
* Docker

## Configuration
> config 파일 생성
* config/database-sample.yaml 참고하여 DB Config 생성
  * config/dev/database.yaml 또는 config/(stage | prod)/database.yaml
* config/redis-sample.yaml 참고하여 Redis Config 생성
  * config/dev/redis.yaml 또는 config/(stage | prod)/redis.yaml
> Proto Buffer 공통 모듈 다운로드
```sh
$ git clone git@gitlab.com:promptech1/data-infuser/infuser-protobuf.git
```

## 배포환경 설정(배포 환경에 따라 dev, stage, prod로 구분되며 각 설정 파일 필요)
> Docker Build 
```sh
make docker-build ENV=dev
```

> Docker Run
```sh
make run-docker
```

> Log Tailing
```sh
make docker-log
```