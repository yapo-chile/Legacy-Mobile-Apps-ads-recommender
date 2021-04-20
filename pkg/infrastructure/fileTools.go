package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FileTools struct to interact with a directory files
type FileTools struct {
	FilesPath string
	Extension string
}

// NewFileTools will create a new instance of a custom File Tool handler
func NewFileTools(filePath string, extension string) *FileTools {
	return &FileTools{
		FilesPath: filePath,
		Extension: extension,
	}
}

// FileExist checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func (t *FileTools) FileExist(fileName string) bool {
	info, err := os.Stat(t.FilesPath + fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// LoadJSONFromFile Loads a Json from file
func (t *FileTools) LoadJSONFromFile(fileName string, dest interface{}) error {
	content, err := os.ReadFile(t.FilesPath + fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, dest)
}

// ListFilesFromPath gets a list of files in a directory
func (t *FileTools) ListFilesFromPath() (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.Walk(t.FilesPath, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("\n file: %+v", info.Name())

		if !info.IsDir() && filepath.Ext(path) == t.Extension {
			files[strings.TrimSuffix(info.Name(), t.Extension)] = info.Name()
		}
		return nil
	})
	return files, err
}

// LoadTemplatesFromFolder loads templates from
// an specific folder
func (t *FileTools) LoadTemplatesFromFolder() map[string]*template.Template {
	fmt.Printf("\nt : %+v", t)
	fmt.Printf("\npath : %+v", t.FilesPath)
	templatesToLoad, _ := t.ListFilesFromPath()
	templates := make(map[string]*template.Template)
	for key, file := range templatesToLoad {
		rawTemplate, err := template.ParseFiles(t.FilesPath + file)
		if err == nil {
			templates[getKey(key)] = rawTemplate
		}
	}
	return templates
}

// LoadTemplatesFromFiles loads templates from
// multiple files of a folder
func (t *FileTools) LoadTemplatesFromFiles(templatesToLoad map[string]string) map[string]*template.Template {
	templates := make(map[string]*template.Template)
	for key, file := range templatesToLoad {
		rawTemplate, err := template.ParseFiles(t.FilesPath + file)
		if err == nil {
			templates[getKey(key)] = rawTemplate
		}
	}
	return templates
}

// LoadTemplateFromFile loads a single template
func (t *FileTools) LoadTemplateFromFile(templateName string) (*template.Template, error) {
	return template.ParseFiles(t.FilesPath + templateName)
}

func getKey(key string) string {
	return strings.Split(key, ".")[0]
}
