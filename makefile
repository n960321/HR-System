.PHONY: run build docker-build docker-run db-run db-remove docker-remove test
cur := $(shell pwd)

db-container-id := $(shell docker ps -a| grep mysql-hr-system | awk '{print $$1}')
hr-system-container-id := $(shell docker ps -a | grep hr-system-server | awk '{print $$1}')

test: 
	@go clean -testcache & go test -timeout 30s -v ./test/...

run:
	@go run main.go -l true -c configs/config.yaml

build:
	@go build -v -o bin/hr-system ./main.go

docker-build:
	@docker build --tag n960321/hr-system:latest --file build/dockerfile .

docker-remove:
	@docker rm -f $(hr-system-container-id)

docker-run:
	@docker run --name hr-system-server \
	-p 8080:8080 \
	--link mysql-hr-system:mysql \
	--volume $(cur)/configs:/app/configs \
	n960321/hr-system:latest

db-remove:
	docker rm -f $(db-container-id)

db-run:
	docker run -d -p 3306:3306 --name mysql-hr-system -e MYSQL_ROOT_PASSWORD=123456 -v $(cur)/deploy/mysql:/docker-entrypoint-initdb.d mysql:8.3.0

	