export APP=author
export CONTAINER=infuser-author
export CONTAINER_VERSION=0.1

build:
	go build ./main.go

container:
	docker build --tag $(CONTAINER):$(CONTAINER_VERSION) .

run-container:
	docker run --rm --detach --publish 9090:9090 --name $(APP) $(CONTAINER):$(CONTAINER_VERSION)

container-log:
	docker logs --follow $(APP)