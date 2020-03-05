package sweb

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
)

func LoadTemplate(root string) *template.Template {
	log.Println("[TEMPLATE]", "loading templates from", root)

	html := template.New("")

	found := false

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
	}

	return html
}
