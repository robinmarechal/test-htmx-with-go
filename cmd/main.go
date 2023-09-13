package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	db "robinmarechal/mod/pkg/database"
	"robinmarechal/mod/pkg/web"
	"strings"
	"text/template"
)

type IndexData struct {
	Name string
}

func loadTemplates() (*template.Template, error) {
	tmpls := template.New("")
	err := filepath.WalkDir("public/views", func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, ".html") {
			log.Printf("parsing file %s \n", path)
			_, err = tmpls.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}
		return err
	})

	return tmpls, err
}

func initDatabase() error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "file:///tmp/text-htmx.db"
	}

	return db.InitDatabase(dbUrl)
}

func main() {
	err := initDatabase()
	if err != nil {
		log.Fatalf("couldn't initialize db: %v", err)
	}

	tmpls, err := loadTemplates()
	if err != nil {
		log.Fatalf("couldn't load template files: %v", err)
	}

	e := web.SetupWebServer(tmpls)
	e.Logger.Fatal(e.Start(":8888"))
}
