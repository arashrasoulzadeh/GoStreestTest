package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar"
	"gopkg.in/yaml.v2"
)

type conf struct {
	Hits   int    `yaml:"hits"`
	Route  string `yaml:"route"`
	Code   int    `yaml:"code"`
	Method string `yaml:"method"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func SendHttpRequest(method string, url string) (*http.Response, error) {
	method = strings.ToUpper(method)

	switch method {
	case http.MethodGet:
		fallthrough
	case http.MethodConnect:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodOptions:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodTrace:
		req, err := http.NewRequest(method, url, new(bytes.Buffer))
		client := &http.Client{}
		resp, err := client.Do(req)
		return resp, err
	}

	fmt.Println("No valid method was provided, using default method GET:")
	resp, err := http.Get(url)
	return resp, err
}

func MakeRequest(url string, method string, ch chan<- string, id int, wg sync.WaitGroup, bar *progressbar.ProgressBar) {
	start := time.Now()

	resp, err := SendHttpRequest(method, url)

	duration := time.Since(start).Seconds()
	if err != nil {
		// handle the error, often:
		bar.Add(1)
		return
	}
	writeToLog(id, resp, err, duration)
	bar.Add(1)
	defer wg.Done()
}

func writeToLog(id int, response *http.Response, e error, duration float64) {
	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%d,%d,%f\n", id, response.StatusCode, duration)); err != nil {
		panic(err)
	}
}

func clearLog() {
	message := []byte("id,code,duration\n")
	err := ioutil.WriteFile("log", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var c conf
	clearLog()
	c.getConf()
	_ = time.Now()
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(c.Hits)
	bar := progressbar.New(c.Hits)
	bar.RenderBlank()

	for i := 1; i <= c.Hits; i++ {
		go MakeRequest(c.Route, c.Method, ch, i, wg, bar)
	}
	wg.Wait()
}
