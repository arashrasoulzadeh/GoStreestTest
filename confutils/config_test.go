package confutils

import (
	"testing"
)

const (
	jsonExample      = "conf.json"
	yamlExample      = "conf.yaml"
	jsonTestFileDest = "../testutils/config.test.json"
	yamlTestFileDest = "../testutils/config.test.yaml"
)

func TestFindType(t *testing.T) {
	var confType string
	confType, err := findType(jsonExample)

	if err != nil {
		t.Error(err)
	}

	if confType != "json" {
		t.Error("wrong type returned -> should return 'json'")
	}

	confType, err = findType(yamlExample)

	if err != nil {
		t.Error(err)
	}

	if confType != "yaml" {
		t.Error("wrong type returned -> should return 'yaml'")
	}
}

func TestGetConf(t *testing.T) {
	var cJSON Conf
	cJSON.GetConf(jsonTestFileDest)

	if cJSON.Hits != 50 {
		t.Error("config file and config struct does not match")
	}

	var cYAML Conf
	cYAML.GetConf(yamlTestFileDest)

	if cYAML.Hits != 50 {
		t.Error("config file and config struct does not match")
	}
}
