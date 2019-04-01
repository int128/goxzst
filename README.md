# goxzst

This is a command to crossbuild, zip, shasum for each GOOS/GOARCH and render templates.
It automates typical release process.

```
Usage:
  goxzst [-d DIR] [-o NAME] [-osarch "GOOS_GOARCH ..."] [-t "FILE ..."] [-tvar "KEY=VALUE ..."]

Options:
  -d string
    	Destination dir (default "dist")
  -o string
    	Target name (default "goxzst")
  -osarch string
    	List of GOOS_GOARCH separated by space (default "linux_amd64 darwin_amd64 windows_amd64")
  -t string
    	List of template files separated by space
  -tvar string
    	List of template variables as KEY=VALUE separated by space
```
