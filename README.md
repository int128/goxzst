# goxzst [![CircleCI](https://circleci.com/gh/int128/goxzst.svg?style=shield)](https://circleci.com/gh/int128/goxzst)

This is a command for cross-build, zip archive, sha digest and template rendering.


## Getting Started

Install the package:

```sh
go get github.com/int128/goxzst
```

To make cross-build, zip and sha256:

```sh
goxzst -o hello
```

It will create the following files:

- `dist/hello_darwin_amd64.zip`
- `dist/hello_darwin_amd64.zip.sha256`
- `dist/hello_linux_amd64.zip`
- `dist/hello_linux_amd64.zip.sha256`
- `dist/hello_windows_amd64.zip`
- `dist/hello_windows_amd64.zip.sha256`

Each zip file contains an executable file as follows:

```
% zipinfo dist/hello_linux_amd64.zip
Archive:  dist/hello_linux_amd64.zip
Zip file size: 2040916 bytes, number of entries: 1
-rwxr-xr-x  2.0 unx  4100026 bl defN 19-Apr-04 14:44 hello
1 file, 4100026 bytes uncompressed, 2040792 bytes compressed:  50.2%
```


## Usage

You can set the following options:

```
Usage:
  goxzst -o NAME [-d DIR] [-osarch "GOOS_GOARCH ..."] [-i "FILE ..."] [-a ALGORITHM] [-t "FILE ..."] [--] [build args]

Options:
  -a string
    	Digest algorithm. One of (sha256|sha512) (default "sha256")
  -d string
    	Output directory (default "dist")
  -i string
    	List of extra files to add to the zip, separated by space
  -o string
    	Output name (mandatory)
  -osarch string
    	List of GOOS_GOARCH separated by space (default "linux_amd64 darwin_amd64 windows_amd64")
  -t string
    	List of template files separated by space
```

goxzst performs the following operations for each platform:

1. Run `go build` with `GOOS` and `GOARCH` environment variables.
1. Archive the executable file into a zip file.
1. Generate the digest of the zip file.

and optionally renders the templates.
Finally it removes the executable files.

### Cross-build

You can set the target platforms by `-osarch` option:

```sh
goxzst -o hello -osarch "linux_amd64 linux_arm"
```

You can pass extra arguments to go build after double dash:

```sh
goxzst -o hello -- -ldflags "-X main.version=$VERSION"
```

### Archive

You can add extra files to the zip by `-i` option:

```sh
goxzst -o hello -i "LICENSE README.md"
```

### Digest

You can set the digest algorithm by `-a` option:

```sh
goxzst -o hello -a sha512
```

### Template

You can pass template files by `-t` option:

```sh
goxzst -o hello -t homebrew.rb
```

goxzst will render the template as a [Go template](https://golang.org/pkg/text/template/)
and write it to a file which has the same filename (not including directory path) of the template.

You can use the following functions and variables in a template.

Name | Description | Example
-----|-------------|--------
`env(string) string`        | Value of the environment variable. | `env "VERSION"`
`sha256(string) string`     | SHA-256 digest of the file.        | `sha256 .linux_amd64_archive`
`sha512(string) string`     | SHA-512 digest of the file.        | `sha512 .linux_amd64_archive`
`.GOOS_GOARCH_executable`   | Path to the executable file.       | `.linux_amd64_executable`
`.GOOS_GOARCH_archive`      | Path to the archive file.          | `.linux_amd64_archive`
`.GOOS_GOARCH_digest`       | Path to the digest file.           | `.linux_amd64_digest`

See also the examples: [homebrew.rb](usecases/testdata/homebrew.rb) and [krew.yaml](usecases/testdata/krew.yaml).


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
