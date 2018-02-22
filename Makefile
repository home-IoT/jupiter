include ./MANIFEST

DATE := $(shell date | sed 's/\ /_/g')

GOPATH ?= $$PWD/../../../..
GOOS ?= linux

SERVER_SWAGGER_FILE=api/server.yml
CLIENT_SWAGGER_FILE=api/client.yml

MOCK=dht22-mock

PKGS := $(shell go list ./... | grep -vF /vendor/)

initialize: clean swagger-generate
	dep init
	$(MAKE) dep

clean:
	rm -rf ./bin/*
	rm -rf ./pkg/*

swagger-serve:
	swagger serve $(SERVER_SWAGGER_FILE)

swagger-validate:
	swagger validate $(SERVER_SWAGGER_FILE)

swagger-gen: swagger-gen-server swagger-gen-client swagger-gen-mock

swagger-gen-server: clean
	swagger generate server -f $(SERVER_SWAGGER_FILE) -t server -A jupiter

swagger-gen-client: clean
	swagger generate client -f $(CLIENT_SWAGGER_FILE) -t client -A jupiter

swagger-gen-mock: clean
	swagger generate server -f $(CLIENT_SWAGGER_FILE) -t dht22-mock -A dht22-mock

dep:
	dep ensure

dep-status:
	dep status

dep-clean: clean
	mkdir -p ./bin
	rm -rf ./vendor/*
	dep ensure

# --- Common Go

go-fmt:
	@go fmt $(PKGS)

go-validate:
	@go vet $(PKGS)
	@golint -set_exit_status $(PKGS)

# --- Jupiter Server

go-build-linux:
	@echo "build linux binary"
	$(MAKE) go-build GOOS=linux GOARCH=amd64 TARGET=$(PROJECT)_linux

go-build-pi:
	@echo "build linux binary for raspberry pi"
	$(MAKE) go-build GOOS=linux GOARCH=arm GOARM=7 TARGET=$(PROJECT)_pi

go-build-windows:
	@echo "build windows binary"
	$(MAKE) go-build GOOS=windows GOARCH=386 TARGET=$(PROJECT).exe

go-build-mac:
	@echo "build Mac binary"
	$(MAKE) go-build GOOS=darwin GOARCH=amd64 TARGET=$(PROJECT)_darwin

go-build-all: go-build-pi go-build-linux go-build-windows go-build-mac

TARGET ?= $(PROJECT)
JUPITER_PACKAGE=github.com/home-IoT/jupiter/internal/jupiter

go-build: 
	go build -ldflags="-X $(JUPITER_PACKAGE).GitRevision=$(shell git rev-parse HEAD) -X $(JUPITER_PACKAGE).BuildVersion=$(VERSION) -X $(JUPITER_PACKAGE).BuildTime=$(DATE)" -i -o ./bin/$(TARGET) server/cmd/jupiter-server/main.go

run: go-build
	./bin/$(TARGET) --port 8080 -c configs/test.yml

# --- Mock

go-build-mock-linux:
	@echo "build linux binary"
	$(MAKE) go-build-mock GOOS=linux GOARCH=amd64 MOCK_TARGET=$(MOCK)_linux

go-build-mock-pi:
	@echo "build linux binary for raspberry pi"
	$(MAKE) go-build-mock GOOS=linux GOARCH=arm GOARM=7 MOCK_TARGET=$(MOCK)_pi

go-build-mock-windows:
	@echo "build windows binary"
	$(MAKE) go-build-mock GOOS=windows GOARCH=amd64 MOCK_TARGET=$(MOCK).exe

go-build-mock-mac:
	@echo "build Mac binary"
	$(MAKE) go-build-mock GOOS=darwin GOARCH=amd64 MOCK_TARGET=$(MOCK)_darwin

go-build-mock-all: go-build-mock-pi go-build-mock-linux go-build-mock-windows go-build-mock-mac

MOCK_TARGET ?= $(MOCK)
DHT22_PACKAGE=github.com/home-IoT/jupiter/internal/dht22

go-build-mock: 
	go build -ldflags="-X $(DHT22_PACKAGE).GitRevision=$(shell git rev-parse HEAD) -X $(DHT22_PACKAGE).BuildVersion=$(VERSION) -X $(DHT22_PACKAGE).BuildTime=$(DATE)" -i -o ./bin/$(MOCK_TARGET) dht22-mock/cmd/dht22-mock-server/main.go

run-mock: 
	./bin/$(MOCK_TARGET) --port 8081

