.PHONY: build clean docker docker-push docker-arm64

GO = CGO_ENABLED=0 GO111MODULE=on go
GOCGO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-scenario

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_scenario_go
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 2.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/rddigital/device-scenario.Version=$(VERSION)"

build: $(MICROSERVICES)

cmd/device-scenario:
	go mod tidy
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_scenario_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t rddigital/device-scenario:$(VERSION) \
		.

docker-push:
	docker push \
	rddigital/device-scenario:$(VERSION)

docker-arm64:
	docker buildx build --platform linux/arm64 -t rddigital/device-scenario-arm64:$(VERSION) --push .