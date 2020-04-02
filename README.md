# goxzst [![CircleCI](https://circleci.com/gh/int128/goxzst.svg?style=shield)](https://circleci.com/gh/int128/goxzst) [![Go Report Card](https://goreportcard.com/badge/github.com/int128/goxzst)](https://goreportcard.com/report/github.com/int128/goxzst)

This is a command to make cross builds, ZIP archives, SHA digests.
As well as it can render templates, for example, Homebrew formula and Kubernetes Krew yaml.

goxzst is an abbreviation of the following letters:

- X (cross builds)
- Z (ZIP archives)
- S (SHA digests)
- T (templates)


## Getting Started

Install [the latest release](https://github.com/int128/goxzst/releases).

```sh
# go
go get github.com/int128/goxzst

# GitHub Releases
curl -sfL -o /tmp/goxzst.zip https://github.com/int128/goxzst/releases/download/v1.0.1/goxzst_linux_amd64.zip
unzip /tmp/goxzst.zip -d ~/bin
```

To make `.zip` and `.sha256` files for the default target platforms:

```sh
goxzst -o hello
```

goxzst will make the following files:

```
dist/hello_darwin_amd64.zip
dist/hello_darwin_amd64.zip.sha256
dist/hello_linux_amd64.zip
dist/hello_linux_amd64.zip.sha256
dist/hello_windows_amd64.zip
dist/hello_windows_amd64.zip.sha256
```

Each zip file contains the executable file:

```
% zipinfo dist/hello_linux_amd64.zip
Archive:  dist/hello_linux_amd64.zip
Zip file size: 2040916 bytes, number of entries: 1
-rwxr-xr-x  2.0 unx  4100026 bl defN 19-Apr-04 14:44 hello
1 file, 4100026 bytes uncompressed, 2040792 bytes compressed:  50.2%
```

To make a Homebrew Formula, create [homebrew.rb](usecases/rendertemplate/testdata/homebrew.rb) and run:

```sh
goxzst -o hello -t homebrew.rb
```

goxzst will make the following files:

```
dist/hello_darwin_amd64.zip
dist/hello_darwin_amd64.zip.sha256
dist/hello_linux_amd64.zip
dist/hello_linux_amd64.zip.sha256
dist/hello_windows_amd64.zip
dist/hello_windows_amd64.zip.sha256
dist/homebrew.rb
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
  -parallelism int
    	Number of parallel build. Default to the current CPU cores
  -t string
    	List of template files separated by space
```

goxzst does the following steps for each platform:

1. Build an executable file for the platform.
1. Pack the executable file into an archive file.
1. Generate the digest of the archive file.

and optionally renders the templates.
Finally it removes the executable files.

### Cross build

You can set the target platforms by `-osarch` option:

```sh
goxzst -o hello -osarch "linux_amd64 linux_arm"
```

You can pass extra arguments to go build after double dash:

```sh
goxzst -o hello -- -ldflags "-X main.version=$VERSION"
```

### Archive

You can add extra files to the archive file by `-i` option:

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

For example, you can render SHA-256 digest of the archive file as follows:

```gotemplate
{{ sha256 .linux_amd64_archive }}
```

See also the examples: [homebrew.rb](usecases/rendertemplate/testdata/homebrew.rb) and [krew.yaml](usecases/rendertemplate/testdata/krew.yaml).


## Related works

This is inspired by [Songmu/goxz](https://github.com/Songmu/goxz).
Thank you for the great work.

You can upload ZIP files to GitHub Releases using [int128/ghcp](https://github.com/int128/ghcp).


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
