package dict

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	items     map[string]string
	EnableEnv bool
)

func init() {
	// 是否从环境变量中取参数
	EnableEnv = false
}

func Load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = parse(b)
	if err == nil {
		if _, ok := items["configDir"]; !ok {
			items["configDir"] = filepath.Dir(filename) + "/"
		}
	}

	return err
}

func LoadFromStr(str string) error {
	return parse([]byte(str))
}

func LoadDefault() error {
	return Load(filepath.Dir(os.Args[0]) + "/config/dict.json")
}

func parse(b []byte) error {
	err := json.Unmarshal(b, &items)
	if err == nil {
		replaceFromEnv()
	}
	return err
}

func Get(key string) string {
	return items[key]
}

func Set(key, value string) {
	items[key] = value
}

func replaceFromEnv() {
	if !EnableEnv {
		return
	}

	str := `\$\{ENV\.([0-9a-zA-Z._-]+)[:]?(.*)\}`
	reg := regexp.MustCompile(str)

	for key, value := range items {
		rs := reg.FindStringSubmatch(value)
		if len(rs) > 0 {
			envSegment := rs[0]
			envName := rs[1]
			newValue := rs[2]

			if val, ok := os.LookupEnv(envName); ok {
				newValue = val
			}

			value := strings.Replace(value, envSegment, newValue, -1)
			Set(key, value)
		}
	}
}
