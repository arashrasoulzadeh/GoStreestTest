package main

import (
	"bytes"
	"encoding/json"
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

//strucure of the yaml file
type conf struct {
	Hits   int                    `yaml:"hits"`
	Route  string                 `yaml:"route"`
	Code   int                    `yaml:"code"`
	Method string                 `yaml:"method"`
	Body   map[string]interface{} `yaml:"body"`
}
//get the configuration
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
//send http call
func SendHttpRequest(method string, url string, body string,wg *sync.WaitGroup) (*http.Response, error) {
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
		reqBody := []byte(body)
		req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		return resp, err
	}

	fmt.Println("No valid method was provided, using default method GET:")
	resp, err := http.Get(url)
	return resp, err
}
//recover if error occured in MakeRequest
func MakeRequestRecovery(wg *sync.WaitGroup){
	defer wg.Done()
	if r := recover(); r != nil {
		fmt.Println("Recovered in f", r)
	}
}
//make request handler
func MakeRequest(url string, method string, body string, ch chan<- string, id int, wg *sync.WaitGroup, bar *progressbar.ProgressBar,f *os.File,ferr error) {
	defer MakeRequestRecovery(wg)
	start := time.Now()
	resp, err := SendHttpRequest(method, url, body,wg)
	duration := time.Since(start).Seconds()
	if err != nil {
		// handle the error, often:
		writeToLog(id, resp, err, duration,f,ferr)
		bar.Add(1)
		return
	}
	writeToLog(id, resp, err, duration,f,ferr)
	bar.Add(1)
}
//write to the log
func writeToLog(id int, response *http.Response, e error, duration float64,f *os.File,ferr error) {
	if e != nil {
	}else{
		if _, ferr = f.WriteString(fmt.Sprintf("%d,%d,%f,%s\n", id, response.StatusCode, duration, "0")); ferr != nil {
			panic(ferr)
		}
	}
}
//clear log for use
func clearLog() {
	message := []byte("id,code,duration,error\n")
	err := ioutil.WriteFile("log.csv", message, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
//main function
func main() {
	f, err := os.OpenFile("log.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
	}()
	var wg sync.WaitGroup
	var c conf
	clearLog()
	c.getConf()
	_ = time.Now()
	ch := make(chan string)
	wg.Add(c.Hits)
	bar := progressbar.New(c.Hits)
	bar.RenderBlank()
	var body string
	if c.Body != nil {
		marshaled, err := json.Marshal(c.Body)
		if err != nil {
			panic(err)
		}
		body = string(marshaled)
	} else {
		marshaled, _ := json.Marshal(map[string]interface{}{})
		body = string(marshaled)
	}
	for i := 1; i <= c.Hits; i++ {
		go MakeRequest(c.Route, c.Method, body, ch, i, &wg, bar,f,err)
	}
	wg.Wait()
}
