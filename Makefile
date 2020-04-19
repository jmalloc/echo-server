BUILD_OS            := windows linux darwin
BUILD_ARCH          := amd64
BUILD_DEBUG_FLAGS   :=
BUILD_RELEASE_FLAGS := -ldflags "-s -w"

BUILD_PATH          := artifacts/build
COVERAGE_PATH       := artifacts/tests/coverage

CURRENT_OS          := $(shell uname | tr A-Z a-z)
CURRENT_ARCH        := amd64

SOURCES             := $(shell find . -name '*.go' -not -path './vendor/*')
PACKAGES            := $(sort $(dir $(SOURCES)))
BINARIES            := $(notdir $(shell find ./cmd -mindepth 1 -maxdepth 1 -type d))
TARGETS             := $(foreach OS,$(BUILD_OS),$(foreach ARCH,$(BUILD_ARCH),$(foreach BIN,$(BINARIES),$(OS)/$(ARCH)/$(BIN))))

test: vendor
	go test $(PACKAGES)

build: $(addprefix $(BUILD_PATH)/debug/$(CURRENT_OS)/$(CURRENT_ARCH)/,$(BINARIES))

debug: $(addprefix $(BUILD_PATH)/debug/,$(TARGETS))

release: $(addprefix $(BUILD_PATH)/release/,$(TARGETS))

docker: Dockerfile $(BUILD_PATH)/release/linux/amd64/echo-server
	@mkdir -p bin
	cp $(BUILD_PATH)/release/linux/amd64/echo-server bin/echo-server
	docker build -t jmalloc/echo-server:dev .

clean:
	@git check-ignore ./* | grep -v ^./vendor | xargs -t -n1 rm -rf
	docker rmi $(docker images -qf dangling=true)

clean-all:
	@git check-ignore ./* | xargs -t -n1 rm -rf
	docker rmi $(docker images -qf dangling=true)

coverage: $(COVERAGE_PATH)/index.html

open-coverage: $(COVERAGE_PATH)/index.html
	open $(COVERAGE_PATH)/index.html

lint: vendor
	go vet $(PACKAGES)
	! go fmt $(PACKAGES) | grep ^

prepare: lint coverage

ci: $(COVERAGE_PATH)/coverage.cov

.PHONY: build test debug release docker clean clean-all coverage open-coverage lint prepare ci

vendor: go.mod
	go mod vendor

$(BUILD_PATH)/%: vendor $(SOURCES)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval PKG   := ./cmd/$(word 4,$(PARTS)))
	$(eval FLAGS := $(if $(filter debug,$(BUILD)),$(BUILD_DEBUG_FLAGS),$(BUILD_RELEASE_FLAGS)))

	GOARCH=$(ARCH) GOOS=$(OS) go build $(FLAGS) -o "$@" "$(PKG)"

$(COVERAGE_PATH)/index.html: $(COVERAGE_PATH)/coverage.cov
	go tool cover -html="$<" -o "$@"

$(COVERAGE_PATH)/coverage.cov: vendor
	@mkdir -p "$(COVERAGE_PATH)"
	go test $(PACKAGES) -covermode=count -coverprofile="$(COVERAGE_PATH)/coverage.cov"
