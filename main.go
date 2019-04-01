package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const usage = `Crossbuild, zip, shasum for each GOOS/GOARCH and render templates.

Usage:
  %[1]s [-d DIR] [-o NAME] [-osarch "GOOS_GOARCH ..."] [-t "FILE ..."] [-tvar "KEY=VALUE ..."]

Options:
`

func main() {
	var o struct {
		dir              string
		target           string
		osarchList       string
		templateFileList string
		templateVarList  string
	}
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.Usage = func() {
		_, _ = fmt.Fprintf(f.Output(), usage, f.Name())
		f.PrintDefaults()
	}
	f.StringVar(&o.dir, "d", "dist", "Destination dir")
	f.StringVar(&o.target, "o", filepath.Base(os.Args[0]), "Target name")
	f.StringVar(&o.osarchList, "osarch", "linux_amd64 darwin_amd64 windows_amd64", "List of GOOS_GOARCH separated by space")
	f.StringVar(&o.templateFileList, "t", "", "List of template files separated by space")
	f.StringVar(&o.templateVarList, "tvar", "", "List of template variables as KEY=VALUE separated by space")
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	buildArgs := f.Args()

	if o.dir != "" {
		if err := os.MkdirAll(o.dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	templateData := make(map[string]string)
	for _, osarch := range strings.Split(o.osarchList, " ") {
		s := strings.SplitN(osarch, "_", 2)
		if len(s) != 2 {
			log.Fatalf("Invalid osarch %s", osarch)
		}
		target := filepath.Join(o.dir, o.target+"_"+osarch)

		if err := (&CrossBuild{}).Do(CrossBuildIn{
			OutputFilename: target,
			Args:           buildArgs,
			GOOS:           s[0],
			GOARCH:         s[1],
		}); err != nil {
			log.Fatal(err)
		}

		if err := (&CreateZip{}).Do(CreateZipIn{
			OutputFilename: target + ".zip",
			Entries: []ZipEntry{
				{
					Name:          o.target,
					InputFilename: target,
				},
			},
		}); err != nil {
			log.Fatal(err)
		}

		shaOut, err := (&CreateSHA{}).Do(CreateSHAIn{
			InputFilename:  target + ".zip",
			OutputFilename: target + ".zip.sha256",
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("sha256=%s", shaOut.SHA256)
		templateData[osarch+"_zip_sha256"] = shaOut.SHA256
	}

	var e []TemplateEntry
	for _, t := range strings.Split(o.templateFileList, " ") {
		e = append(e, TemplateEntry{
			InputFilename:  t,
			OutputFilename: filepath.Join(o.dir, filepath.Base(t)),
		})
	}

	for _, v := range strings.Split(o.templateVarList, " ") {
		kv := strings.SplitN(v, "=", 2)
		templateData[kv[0]] = kv[1]
	}

	log.Printf("Rendering templates with data %+v", templateData)
	if err := (&RenderTemplate{}).Do(RenderTemplateIn{
		Entries: e,
		Data:    templateData,
	}); err != nil {
		log.Fatal(err)
	}
}
