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
	docker build .

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

GLIDE := $(GOPATH)/bin/glide
$(GLIDE):
	go get -u github.com/Masterminds/glide

vendor: glide.lock | $(GLIDE)
	$(GLIDE) install
	@touch vendor

glide.lock: glide.yaml | $(GLIDE)
	$(GLIDE) update
	@touch vendor

$(BUILD_PATH)/%: vendor $(SOURCES)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval PKG   := ./cmd/$(word 4,$(PARTS)))
	$(eval FLAGS := $(if $(filter debug,$(BUILD)),$(BUILD_DEBUG_FLAGS),$(BUILD_RELEASE_FLAGS)))

	GOARCH=$(ARCH) GOOS=$(OS) go build $(FLAGS) -o "$@" "$(PKG)"

GOCOVMERGE := $(GOPATH)/bin/gocovmerge
$(GOCOVMERGE):
	go get -u github.com/wadey/gocovmerge

$(COVERAGE_PATH)/index.html: $(COVERAGE_PATH)/coverage.cov
	go tool cover -html="$<" -o "$@"

$(COVERAGE_PATH)/coverage.cov: $(foreach P,$(PACKAGES),$(COVERAGE_PATH)/$(P)coverage.partial) | $(GOCOVMERGE)
	@mkdir -p $(@D)
	$(GOCOVMERGE) $^ > "$@"

.SECONDEXPANSION:
%/coverage.partial: vendor $$(subst $(COVERAGE_PATH)/,,$$(@D))/*.go
	$(eval PKG := $(subst $(COVERAGE_PATH)/,,$*))
	@mkdir -p $(@D)
	@touch "$@"
	go test "$(PKG)" -covermode=count -coverprofile="$@"
