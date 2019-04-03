# goxzst [![CircleCI](https://circleci.com/gh/int128/goxzst.svg?style=shield)](https://circleci.com/gh/int128/goxzst)

This is a command for cross-build, zip, shasum for each GOOS/GOARCH and rendering templates.
It automates typical release process.


## Getting Started

Install:

```
go get github.com/int128/goxzst
```

To make cross-build, zip and sha256 for the default platforms:

```
% goxzst
2019/04/03 12:55:05 GOOS=linux GOARCH=amd64 go build -o dist/goxzst_linux_amd64
2019/04/03 12:55:05 Creating dist/goxzst_linux_amd64.zip
2019/04/03 12:55:06 Creating dist/goxzst_linux_amd64.zip.sha256
2019/04/03 12:55:06 GOOS=darwin GOARCH=amd64 go build -o dist/goxzst_darwin_amd64
2019/04/03 12:55:06 Creating dist/goxzst_darwin_amd64.zip
2019/04/03 12:55:06 Creating dist/goxzst_darwin_amd64.zip.sha256
2019/04/03 12:55:06 GOOS=windows GOARCH=amd64 go build -o dist/goxzst_windows_amd64
2019/04/03 12:55:06 Creating dist/goxzst_windows_amd64.zip
2019/04/03 12:55:07 Creating dist/goxzst_windows_amd64.zip.sha256
```

You can set the target platforms by `-osarch` option:

```
% goxzst -osarch "linux_amd64 linux_arm"
```

You can pass extra arguments to go build after double dash:

```
% goxzst -- -ldflags "-X version=$VERSION"
```

You can render the templates by `-t` option:

```
% goxzst -t "homebrew.rb krew.yaml"
```


## Usage

goxzst performs the following operations for each platform:

1. Run `go build` with `GOOS` and `GOARCH` environment variables.
1. Archive the executable file into a zip file.
1. Generate SHA-256 digest of the zip file.

and optionally generates files from the templates.

Finally, the following files will be created:

- `DIR/NAME_GOOS_GOARCH`: executable
- `DIR/NAME_GOOS_GOARCH.zip`: archive
- `DIR/NAME_GOOS_GOARCH.zip.sha256`: digest

You can set the following options:

```
Usage:
  goxzst [-d DIR] [-o NAME] [-osarch "GOOS_GOARCH ..."] [-t "FILE ..."] [--] [build args]

Options:
  -d string
    	Output directory (default "dist")
  -o string
    	Output name (default "goxzst")
  -osarch string
    	List of GOOS_GOARCH separated by space (default "linux_amd64 darwin_amd64 windows_amd64")
  -t string
    	List of template files separated by space
```

### Template

goxzst creates a file which has same filename (not including directory path) of the template.

You can use the following functions and variables in a template.

Name | Description | Example
-----|-------------|--------
`env(string) string`        | Value of the environment variable. | `env "VERSION"`
`.GOOS_GOARCH_zip_sha256`   | SHA-256 digest of the zip file.    | `.linux_amd64_zip_sha256`


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
