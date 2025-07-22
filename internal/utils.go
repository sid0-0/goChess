package internal

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/davecgh/go-spew/spew"
)

func LoadAllTemplates(location string) (*template.Template, error) {
	var allTemplates *template.Template
	spew.Println("Loading templates from:", location)
	allTemplates = template.New("").Funcs(sprig.FuncMap())
	err := filepath.WalkDir(location, func(path string, _ fs.DirEntry, _ error) error {
		if strings.HasSuffix(path, ".gohtml") {
			_, err := allTemplates.ParseFiles(path)
			if err == nil {
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}
