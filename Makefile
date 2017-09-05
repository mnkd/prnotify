NAME     := prnotify
VERSION  := 0.8
REVISION := $(shell git rev-parse --short HEAD)
SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""

bin/$(NAME): $(SRCS) format
	go build $(LDFLAGS) -o bin/$(NAME)

linux: $(SRCS) format
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(NAME)

format:
	go fmt $(SRCS)

clean:
	rm -rf bin/*

install:
	go install $(LDFLAGS)

.PHONY: format, clean, install
