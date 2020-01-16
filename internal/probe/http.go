package probe

import (
	"log"
	"net/http"
	"time"
)

type Http struct {
	Probe `yaml:",inline"`
	Url   string `yaml:"url"`
}

func (h *Http) Check() {
	for {
		resp, err := http.Get(h.Url)
		if err != nil {
			log.Printf("<<HTTP PROBE [%s]>> an error occured while doing HTTP call to [%s]. got '%s'\n", h.Name, h.Url, err)
			h.app.UpdateTemplateData(h.Name, StatusDown)
		} else {
			log.Printf("<<HTTP PROBE [%s]>> Service is alive.\n", h.Name)
			h.app.UpdateTemplateData(h.Name, StatusUp)
			resp.Body.Close()
		}
		time.Sleep(time.Duration(h.Delay) * time.Second)
	}
}
