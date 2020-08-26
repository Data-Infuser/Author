APP=author
CONTAINER=infuser-author

VERSION:=0.1
ENV:=dev #서비스 환경에 따라 dev, stage, prod로 구분

build:
	go build ./main.go

docker-build:
	docker build --tag $(CONTAINER):$(VERSION) --build-arg=AUTHOR_ENV=$(ENV) .

run-docker:
	docker run --rm --detach --publish 9090:9090 --name $(APP) $(CONTAINER):$(VERSION)

docker-log:
	docker logs --follow $(APP)

.PHONY: build docker run-docker docker-log