package language

import (
	"ApiMessenger/models"
	"encoding/json"
	"fmt"
	"os"
)

func Language(key string) string {

	lang := models.Language()

	var path string

	var data map[string]interface{}

	if lang == "en" {
		path = "language/lang_en.json"
	}
	if lang == "ru" {
		path = "language/lang_ru.json"
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("JSON decoding error:", err)
		return ""
	}
	return data[key].(string)
}
