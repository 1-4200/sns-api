GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=sns-api

all: build

.PHONY:build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v -a

.PHONY:clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY:run
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

.PHONY:build-linux
build-linux:
	$(GOCLEAN)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v -a

.PHONY:dep
dep:
	$(GOGET) -v
