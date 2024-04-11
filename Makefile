

all: generate build

generate:
	./scripts/buf generate

build:
	cd cmd/server && go build -o ../../bin/market $(GOLDFLAGS)
	cd cmd/publisher && go build -o ../../bin/publisher $(GOLDFLAGS)

