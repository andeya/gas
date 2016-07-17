
TARGETS_NOVENDOR := $(shell glide novendor)

install:
	glide install

lint:
	golint .

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test -v -race .

coverage:
	sh test.sh atomic

html: coverage
	go tool cover -html=coverage.txt
