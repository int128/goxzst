.PHONY: all check run

all: goxzst

check:
	go vet
	go test -v ./...

goxzst: $(wildcard **/*.go)
	go build -o $@

run: goxzst
	VERSION=v1.0.0 ./goxzst -i LICENSE -t "usecases/testdata/homebrew.rb usecases/testdata/krew.yaml"
	zipinfo dist/goxzst_linux_amd64.zip
	zipinfo dist/goxzst_darwin_amd64.zip
	zipinfo dist/goxzst_windows_amd64.zip
