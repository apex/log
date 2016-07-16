
TESTFLAG ?= -cover

test:
	@go test $(TESTFLAG) ./...
.PHONY: test
