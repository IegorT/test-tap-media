
# Env parameters
include .env
EXPORTENV=export PORT=${PORT} LOCATION_DB_PATH=${LOCATION_DB_PATH} UA_PARSER_REGEXP_PATH=${UA_PARSER_REGEXP_PATH}

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./bin/app

all: test build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./src/app
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v
		${EXPORTENV} && $(BINARY_NAME)

deps:
		$(GOGET) github.com/oschwald/maxminddb-golang
		$(GOGET) github.com/ua-parser/uap-go/uaparser
		$(GOGET) github.com/julienschmidt/httprouter
