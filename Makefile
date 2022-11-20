PKGS ?= $(shell go list ./...)

build-local:
	go build -o app ./main.go

lint:
	$$GOPATH/bin/golint ${PKGS}

test:
	go test -race -covermode=atomic -coverprofile coverage.out -v ${PKGS}

run:
	go run main.go

build:
	$$GOPATH/bin/gox -osarch="linux/amd64" -output="app" ./

ensure-dependencies:
	go mod tidy

services:
	docker-compose up -d

generate-mock:
	cd repository && mockery --name=ITaskRepository --filename=task.go --outpkg=mock --output=../mock

generate-docs:
	swag init --parseDependency