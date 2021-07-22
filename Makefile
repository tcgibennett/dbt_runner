FONT_ESC := $(shell printf '\033')
FONT_BOLD := ${FONT_ESC}[1m
FONT_NC := ${FONT_ESC}[0m # No colour

all:
	@echo "Use a specific goal. To list all goals, type 'make help'"

.PHONY: version # Prints project version
version:
	@cat VERSION

.PHONY: build # Builds the project
build:
	@go build

.PHONY: build_linux # Builds for linux
build_linux:
	@env GOOS=linux GOARCH=amd64 go build

.PHONY: install # Installs the project
install:
	@go install

.PHONY: test # Runs unit tests
test:
	@go test -v ./...

.PHONY: build-docker # Build Docker image
build-docker:
	$(shell $(MAKE) build_linux)
	@docker build -t talend/dbt_runner:$(shell $(MAKE) version) -f Dockerfile .
	@docker tag talend/dbt_runner:$(shell $(MAKE) version) talend/dbt_runner:latest

.PHONY: help # Generate list of goals with descriptions
help:
	@echo "Available goals:\n"
	@grep '^.PHONY: .* #' Makefile | sed "s/\.PHONY: \(.*\) # \(.*\)/${FONT_BOLD}\1:${FONT_NC}\2~~/" | sed $$'s/~~/\\\n/g' | sed $$'s/~/\\\n\\\t/g'
