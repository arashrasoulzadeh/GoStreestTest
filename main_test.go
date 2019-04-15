package main

import (
	"net/http"
	"strings"
	"testing"
)

const (
	conf_file_path = "./testutils/config.test.yaml"
)

func TestGetConf(t *testing.T) {
	var c conf
	c.getConf(conf_file_path)
	if c.Code != 200 {
		t.Error("Wrong value for - Code -")
	}

	if c.Hits != 50 {
		t.Error("Wrong value for - Hits -")
	}

	if strings.ToUpper(c.Method) != http.MethodGet {
		t.Error("Wrong value for - Method -")
	}

	if c.Route != "test_route" {
		t.Error("Wrong value for - Route -")
	}
}
