# goxzst

This is a command to crossbuild, zip, shasum for each GOOS/GOARCH and render templates.
It automates typical release process.

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
