package main

import (
	"github.com/schollz/progressbar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type conf struct {
	Hits  int    `yaml:"hits"`
	Route string `yaml:"route"`
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
func MakeRequest(url string, ch chan<- string, id int, wg sync.WaitGroup ,bar *progressbar.ProgressBar) {
 	start := time.Now()
	_, _ = http.Get(url)
	_ = time.Since(start).Seconds()
 	bar.Add(1)
 	defer wg.Done()
}

func main() {
	var c conf
	c.getConf()
	_ = time.Now()
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(c.Hits)
	bar := progressbar.New(c.Hits)
	bar.RenderBlank()

	for i := 1; i <= c.Hits; i++ {
		go MakeRequest(c.Route, ch, i, wg,bar)
	}
 	wg.Wait()
}
