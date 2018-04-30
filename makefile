THIS_FILE := $(lastword $(MAKEFILE_LIST))

deps:
	@echo $@
	@go get -u github.com/go-sql-driver/mysql
	@go get -u github.com/satori/go.uuid


test:
	@echo $@
	@go test -v ./test/...

.PHONY: deps test
