.PHONY: format
format:
	@find . -type f -name "*.go" -print0 | xargs -0 gofmt -s -w

.PHONY: bench
bench:
	GOPATH=$(GOPATH) go test -bench=.

# Clean junk
.PHONY: clean
clean:
	GOPATH=$(GOPATH) go clean ./...

.PHONY: clean-mac
clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm

.PHONY: test
test:
	-rm coverage.txt
	@for package in $$(go list ./... | grep -v example) ; do \
		GOPATH=$(GOPATH) go test -race -coverprofile=profile.out -covermode=atomic $$package ; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt ; \
			rm profile.out ; \
		fi \
	done

update_deps:
	GOPATH=$(GOPATH) GO111MODULE=on go mod verify
	GOPATH=$(GOPATH) GO111MODULE=on go mod tidy
	GOPATH=$(GOPATH) rm -rf vendor
	GOPATH=$(GOPATH) GO111MODULE=on go mod vendor