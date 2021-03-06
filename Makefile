.PHONY: all check run release clean

all: goxzst

check:
	go vet
	go test -v -race ./...

goxzst: $(wildcard **/*.go)
	go build -o $@

run: goxzst
	-./goxzst -h
	VERSION=v1.0.0 ./goxzst -o example -i LICENSE -t "usecases/rendertemplate/testdata/homebrew.rb usecases/rendertemplate/testdata/krew.yaml" -- -ldflags "-X main.version=v1.2.3"
	zipinfo dist/example_linux_amd64.zip
	zipinfo dist/example_darwin_amd64.zip
	zipinfo dist/example_windows_amd64.zip

release: goxzst
	-rm -r dist/
	./goxzst -o goxzst -- -ldflags "-X main.version=$(CIRCLE_TAG)"
	ghcp release -u "$(CIRCLE_PROJECT_USERNAME)" -r "$(CIRCLE_PROJECT_REPONAME)" -t "$(CIRCLE_TAG)" dist/

clean:
	-rm -r dist/
