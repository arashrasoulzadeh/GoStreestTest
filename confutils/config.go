package confutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Conf config file strcuture
type Conf struct {
	Hits    int                    `yaml:"hits" json:"hits"`
	Route   string                 `yaml:"route" json:"route"`
	Code    int                    `yaml:"code" json:"code"`
	Method  string                 `yaml:"method" json:"method"`
	Body    map[string]interface{} `yaml:"body" json:"body"`
	Headers map[string]string      `yaml:"headers" json:"headers"`
}

// GetConf method parses json or yaml config file
func (c *Conf) GetConf(dest string) *Conf {
	filetype, err := findType(dest)

	if err != nil {
		panic("Wrong congig file provided.")
	}

	file, err := ioutil.ReadFile(dest)
	if err != nil {
		panic("failed getting config file ")
	}

	if filetype == "yaml" {
		err = yaml.Unmarshal(file, c)
		if err != nil {
			panic("failed to parse config file ")
		}

		return c
	}

	json.Unmarshal(file, c)
	return c
}

// findType determines config file file type, can be `json` or `yaml`
func findType(filename string) (string, error) {
	if isJSON, _ := regexp.Match(`\.json$`, []byte(filename)); isJSON {
		return "json", nil
	} else if isYAML, _ := regexp.Match(`\.yaml$`, []byte(filename)); isYAML {
		return "yaml", nil
	} else {
		return "", fmt.Errorf("invalid file type")
	}
}
