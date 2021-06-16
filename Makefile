.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go
GOCGO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-scenario

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_scenario_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/rddigital/device-scenario.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-scenario:
	go mod tidy
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd

test:
	go mod tidy
	$(GOCGO) test ./... -coverprofile=coverage.out
	$(GOCGO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_scenario_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t rddigital/device-scenario:$(GIT_SHA) \
		-t rddigital/device-scenario:$(VERSION)-dev \
		.
