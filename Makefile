
TARGETS_NOVENDOR := $(shell glide novendor)

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test -v ./

coverage:
	sh test.sh atomic

html: coverage
	go tool cover -html=coverage.txt && unlink coverage.txt
