package sweb

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func LoadTemplate(root string) *template.Template {
	log.Println("[TEMPLATE]", "loading templates from", root)

	//found := false

	var html *template.Template

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			buf, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			rel, err := filepath.Rel(root, path)

			if err != nil {
				return err
			}

			name := filepath.ToSlash(rel)

			if html == nil {
				html = template.New(name)
			} else {
				html = html.New(name)
			}

			if _, err := html.Parse(string(buf)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatalln(err)
	}

	/*
	if files, err := ioutil.ReadDir(root); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				html, err = html.ParseFiles(filepath.Join(root, file.Name()))

				if err != nil {
					log.Panicln(err)
				}

				found = true
			}
		}
	} else {
		log.Println(err)
	}

	if !found {
		err := fmt.Errorf ("no templates found in folder %v", root)
		log.Fatalln(err)
	}*/

	return html
}
