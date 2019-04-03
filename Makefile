.PHONY: all check run

all: goxzst

check:
	go vet
	go test -v ./...

goxzst: $(wildcard **/*.go)
	go build -o $@

run: goxzst
	-./goxzst -h
	VERSION=v1.0.0 ./goxzst -o example -i LICENSE -t "usecases/testdata/homebrew.rb usecases/testdata/krew.yaml" -- -ldflags "-X main.version=v1.2.3"
	zipinfo dist/example_linux_amd64.zip
	zipinfo dist/example_darwin_amd64.zip
	zipinfo dist/example_windows_amd64.zip
