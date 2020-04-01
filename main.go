package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

func walkDir(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var dir string
	var found bool
	for _, f := range files {
		dir = filepath.Join(path, f.Name())

		if f.IsDir() {
			if strings.HasPrefix(f.Name(), ".") {
				// skip hidden dirs
				continue
			}

			if strings.Contains(f.Name(), "modules") {
				// skip module dirs
				continue
			}
			walkDir(dir)
		}

		if found {
			continue
		}

		if f.Name() == ".terraform-version" {
			found = true
			// if tfenv file found for now we continue
			continue
		}

		if filepath.Ext(f.Name()) != ".tf" {
			// skip non-terraform files
			continue
		}

		// best attempt to identifying correct folder
		matched, err := regexp.MatchString(`main|provider|remote|vars`, f.Name())
		if err != nil {
			log.Fatal(err)
		}

		if !matched {
			// skip files not matching pattern
			continue
		}

		// best attempt to identifying hcl syntax
		version := "0.12.24"
		if f.Name() == "vars.tf" {
			content, err := ioutil.ReadFile(dir)
			if err != nil {
				log.Fatal(err)
			}
			hcl1, err := isLegacy(content)
			if err != nil {
				log.Fatal(err)
			}
			if hcl1 {
				version = "0.11.14"
			}
		}

		tfv := filepath.Join(filepath.Dir(dir), ".terraform-version")
		err = ioutil.WriteFile(tfv, []byte(version), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(filepath.Dir(dir))
	}
}

func isLegacy(b []byte) (bool, error) {
	return regexp.MatchString(`(type.+=.+"string")`, string(b))
}

func main() {
	walkDir(".")
}
