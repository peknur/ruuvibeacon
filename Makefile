GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BIN_DIR=./bin
BIN_NAME=$(BIN_DIR)/ruuvibeacon
GOFILES=./cmd/ruuvibeacon.go
all: clean build
all-arm: clean build-arm
build: 
	$(GOBUILD) -o $(BIN_NAME) -v $(GOFILES)
clean: 
	$(GOCLEAN)
	rm -f $(BIN_NAME)*
tracer: 
	$(GOBUILD) -race -o $(BIN_NAME).trace -v $(GOFILES)
build-arm: 
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BIN_NAME).arm -v $(GOFILES)
