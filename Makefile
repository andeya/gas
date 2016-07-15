
TARGETS_NOVENDOR := $(shell glide novendor)

lint:
	golint .

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test -v -race ./

coverage:
	sh test.sh count

html: coverage
	go tool cover -html=coverage.txt
