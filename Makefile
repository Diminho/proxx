ifndef GO_OS
GO_OS=linux
endif

ifndef GO_ARCH
GO_ARCH=amd64
endif

## enabling go module inside the GOPATH
export GO111MODULE=on
export GOSUMDB=off

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor

packages = \
	./service/launch \
	./service/ticket \
	./web/launch \
	./web/ticket \

.PHONY: test
test:
	@$(foreach package,$(packages), \
		set -e; \
		go test -tags musl -coverprofile $(package)/cover.out.tmp -covermode=count $(package); \
		cat $(package)/cover.out.tmp | grep -v "_mock.go" > $(package)/cover.out; \
		rm $(package)/cover.out.tmp;)


.PHONY: code-quality-print ## Run golang-cilint with printing to stdout
code-quality-print: bin/golangci-lint
	./bin/golangci-lint --build-tags=musl --exclude-use-default=false --tests=false --out-format tab run ./...

.PHONY: code-quality-ci
code-quality-ci: bin/golangci-lint ## Run golang-ci linter for all packages with printing to stdout
	bin/golangci-lint run --out-format colored-line-number

.PHONY: golangci-install
golangci-install:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.49.0


