package probe

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Http struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Delay int    `yaml:"delay"`
}

func (h *Http) Probe() {
	for {
		resp, err := http.Get(h.Url)
		if err != nil {
			// handle error from call --> mark service as down...
			log.Printf("<<HTTP PROBE [%s]>> an error occured while doing HTTP call to [%s]. got '%s'\n", h.Name, h.Url, err)
			time.Sleep(time.Duration(h.Delay) * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("<<HTTP PROBE [%s]>> an error occured while reading response body. got '%s'\n", h.Name, err)
			time.Sleep(time.Duration(h.Delay) * time.Second)
			continue
		}

		log.Printf("<<HTTP PROBE [%s]>> Service returned [%s]\n", h.Name, string(body))

		resp.Body.Close()
		time.Sleep(time.Duration(h.Delay) * time.Second)
	}
}
