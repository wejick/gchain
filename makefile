# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

EXAMPLES = $(shell find example -type d -not -path '*/\.*')

# Test and integration tag
INTEGRATIONTAGS=integration

all: setup build
setup:
	@echo "Setting up the Go environment"
	$(GOGET) -v -t -d ./...
test:
	@echo "Running the Go unit tests"
	$(GOTEST) -v  ./...
integration-test:
	@echo "Running the Go integration tests"
	$(GOTEST) -v -tags=$(INTEGRATIONTAGS) ./...
build-example:
	@echo "Building examples"
	@for dir in $(shell find ./example -type d); do \
		(cd $$dir && go build .); \
	done
generate-mock:
	mockery --all --keeptree
	-rm -r mocks/internal_mock
	mv mocks/internal mocks/internal_mock