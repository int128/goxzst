package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var o struct {
		dir        string
		target     string
		osarchList string
	}
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&o.dir, "d", "dist", "Destination dir")
	f.StringVar(&o.target, "o", filepath.Base(os.Args[0]), "Target name")
	f.StringVar(&o.osarchList, "osarch", "linux_amd64 darwin_amd64 windows_amd64", "List of GOOS_GOARCH separated by space")
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

	templateData["version"] = "v1.0.0" //TODO: as options

	log.Printf("Rendering templates with data %+v", templateData)
	if err := (&RenderTemplate{}).Do(RenderTemplateIn{
		Entries: []TemplateEntry{
			{
				//TODO: as options
				InputFilename:  "testdata/goxzst.rb",
				OutputFilename: o.dir + "/goxzst.rb",
			},
		},
		Data: templateData,
	}); err != nil {
		log.Fatal(err)
	}
}
