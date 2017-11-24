ifeq ($(origin VERSION), undefined)
	VERSION != 1
endif

HOST_GOOS=$(shell go env GOOS)
HOST_GOARCH=$(shell go env GOARCH)
REPOPATH = github.com/rafael/cryptoticker

VERBOSE_1 := -v
VERBOSE_2 := -v -x
WHAT :=  googlecreds ticketupdater

build: vendor
	for target in $(WHAT); do \
		$(BUILD_ENV_FLAGS) go build $(VERBOSE_$(V)) -o bin/$$target -ldflags "-X $(REPOPATH).Version=$(VERSION)" ./cmd/$$target; \
	done

