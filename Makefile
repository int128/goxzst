.PHONY: all check run

all: goxzst

check:
	go vet
	go test -v ./...

goxzst:
	go build -o $@

run: goxzst
	./goxzst
