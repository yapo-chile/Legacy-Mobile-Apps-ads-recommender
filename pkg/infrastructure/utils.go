package infrastructure

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
)

// LoadJSONFromFile loads the data in the JSON file in the given path
// and loads it in the given interface
func LoadJSONFromFile(filePath string, dest interface{}) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, dest)
}

// LoadTemplatesFromFiles loads the given templates in the given folder and returns a map
// with all the loaded templates
func LoadTemplatesFromFiles(
	templatesFolder string,
	templatesToLoad map[string]string,
) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	for key, file := range templatesToLoad {
		rawTemplate, err := template.ParseFiles(templatesFolder + file)
		if err == nil {
			templates[key] = rawTemplate
		} else {
			return templates, err
		}
	}
	return templates, nil
}
