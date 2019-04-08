package main

import (
	"fmt"
	"github.com/schollz/progressbar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type conf struct {
	Hits  int    `yaml:"hits"`
	Route string `yaml:"route"`
	Code  int    `yaml:"code"`
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

func MakeRequest(url string, ch chan<- string, id int, wg sync.WaitGroup, bar *progressbar.ProgressBar) {
	start := time.Now()
	resp, err := http.Get(url)
	_ = time.Since(start).Seconds()
	writeToLog(id, resp, err)
	bar.Add(1)
	defer wg.Done()
}

func writeToLog(id int, response *http.Response, e error) {

	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%d,%d\n", id, response.StatusCode)); err != nil {
		panic(err)
	}

}

func clearLog(){
	message := []byte("id,code\n")
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
		go MakeRequest(c.Route, ch, i, wg, bar)
	}
	wg.Wait()
}
