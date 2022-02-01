.PHONY: all clean clean-coverage generate install install-dependencies install-tools lint test test-verbose test-verbose-with-coverage

export ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
export PKG := github.com/smgladkovskiy/go-mutesting
export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

export TEST_TIMEOUT_IN_SECONDS := 240

$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

ifdef ARGS
	PKG_TEST := $(ARGS)
else
	PKG_TEST := $(PKG)/...
endif

all: install-dependencies install-tools install lint test
.PHONY: all

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
.PHONY: clean

clean-coverage:
	find $(ROOT_DIR) | grep .coverprofile | xargs rm
.PHONY: clean-coverage

generate: clean
	go generate $(PKG)/...
.PHONY: generate

install:
	go install -v $(PKG)/...
.PHONY: install

install-dependencies:
	go mod vendor
	go test -i -v $(PKG)/...
.PHONY: install-dependencies

install-tools:
	# generation
	go get golang.org/x/tools/cmd/stringer

	# code coverage
	go get golang.org/x/tools/cmd/cover
	go get github.com/onsi/ginkgo/ginkgo/...
	go get github.com/modocache/gover/...
	go get github.com/mattn/goveralls/...
.PHONY: install-tools

#test:
#	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" $(PKG_TEST)
#.PHONY: test
#
#test-verbose:
#	go test -race -test.timeout "$(TEST_TIMEOUT_IN_SECONDS)s" -v $(PKG_TEST)
#.PHONY: test-verbose
#
#test-verbose-with-coverage:
#	ginkgo -r -v -cover -race -skipPackage="testdata"
#.PHONY: test-verbose-with-coverage


# ----
## LINTER stuff start

linter_include_check:
	@[ -f linter.mk ] && echo "linter.mk include exists" || (echo "getting linter.mk from github.com" && curl -sO https://raw.githubusercontent.com/spacetab-io/makefiles/master/golang/linter.mk)

.PHONY: lint
lint: linter_include_check
	@make -f linter.mk go_lint

## LINTER stuff end
# ----

# ----
## TESTS stuff start

tests_include_check:
	@[ -f tests.mk ] && echo "tests.mk include exists" || (echo "getting tests.mk from github.com" && curl -sO https://raw.githubusercontent.com/spacetab-io/makefiles/master/golang/tests.mk)

tests: tests_include_check
	@make -f tests.mk go_tests
.PHONY: tests

tests_html: tests_include_check
	@make -f tests.mk go_tests_html
	@open coverage.html
.PHONY: tests

## TESTS stuff end
# ----

build:
	go build -mod=vendor -o bin/go-mutesting .