include ./MANIFEST
include ./scripts/release.mk

DATE := $(shell date | sed 's/\ /_/g')

GOPATH ?= $$PWD/../../../..
GOOS ?= linux

SERVER_SWAGGER_FILE=api/server.yml
CLIENT_SWAGGER_FILE=api/client.yml

PKGS := $(shell go list ./... | grep -vF /vendor/)

JUPITER_PACKAGE=github.com/home-IoT/jupiter/internal/jupiter

# --- Repo 

initialize: clean swagger-gen
	dep init
	$(MAKE) go-dep

clean:
	mkdir -p bin
	rm -rf ./bin/*
	rm -rf ./pkg/*

# --- Tools

get-tools:
	go get -u github.com/golang/lint/golint
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# --- Swagger

get-swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger-serve:
	swagger serve $(SERVER_SWAGGER_FILE)

swagger-validate:
	swagger validate $(SERVER_SWAGGER_FILE)

swagger-gen: swagger-gen-server swagger-gen-client

swagger-gen-server: clean
	swagger generate server -f $(SERVER_SWAGGER_FILE) -t server -A jupiter

swagger-gen-client: clean
	swagger generate client -f $(CLIENT_SWAGGER_FILE) -t client -A jupiter

# --- Common Go

go-dep:
	dep ensure

go-dep-status:
	dep status

go-dep-clean: clean
	mkdir -p ./bin
	rm -rf ./vendor/*
	dep ensure

go-fmt:
	@go fmt $(PKGS)

go-validate:
	@echo go vet
	@go vet $(PKGS)
	@echo golint
	@golint -set_exit_status $(PKGS)

# --- Jupiter Server

go-build-linux:
	@echo "build linux binary"
	$(MAKE) go-build GOOS=linux GOARCH=amd64 TARGET=$(PROJECT)-linux-amd64

go-build-pi:
	@echo "build linux binary for raspberry pi"
	$(MAKE) go-build GOOS=linux GOARCH=arm GOARM=7 TARGET=$(PROJECT)-linux-arm7

go-build-windows:
	@echo "build windows binary"
	$(MAKE) go-build GOOS=windows GOARCH=386 TARGET=$(PROJECT)-windows-386.exe

go-build-mac:
	@echo "build Mac binary"
	$(MAKE) go-build GOOS=darwin GOARCH=amd64 TARGET=$(PROJECT)-darwin-amd64

TARGET ?= $(PROJECT)

go-build: 
	go build -ldflags="-X $(JUPITER_PACKAGE).GitRevision=$(shell git rev-parse HEAD) -X $(JUPITER_PACKAGE).BuildVersion=$(VERSION) -X $(JUPITER_PACKAGE).BuildTime=$(DATE)" -i -o ./bin/$(TARGET) server/cmd/jupiter-server/main.go

go-build-all: go-build-pi go-build-linux go-build-windows go-build-mac

run: go-build
	./bin/$(TARGET) --port 8080 -c configs/test.yml

# --- Release

go-release-all: clean 
	$(MAKE) go-build-all
	mkdir -p ./release
	rm -rf ./release/*
	chmod +x bin/*
	cp ./bin/* ./release
	for bf in ./release/*; do shasum -a 256 "$$bf" > "$$bf".sha256; done


