package dict

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var items map[string]string

func Load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &items)
}

func LoadFromStr(str string) error {
	return json.Unmarshal([]byte(str), &items)
}

func LoadDefault() error {
	return Load(filepath.Dir(os.Args[0]) + "/config/dict.json")
}

func Get(key string) string {
	return items[key]
}

func Set(key, value string) {
	items[key] = value
}
