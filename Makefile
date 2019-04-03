.PHONY: all check run

all: goxzst

check:
	go vet
	go test -v ./...

goxzst:
	go build -o $@

run: goxzst
	VERSION=v1.0.0 ./goxzst -t "usecases/testdata/homebrew.rb usecases/testdata/krew.yaml"
