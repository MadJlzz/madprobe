package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	TemplateData interface{}
	tpl          *template.Template
)

type WebApp struct {
	Port        int
	TemplateSrc string
}

func StartApp(webApp WebApp) {
	parseAndLoadTemplate(webApp)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Printf("Trying to start WebApp on port [%d]...\n", webApp.Port)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("error while trying to server mux server on port [%d]. got %s\n", webApp.Port, err)
	}
}

func parseAndLoadTemplate(webApp WebApp) {
	log.Printf("Trying to parse HTML template [%s]...\n", webApp.TemplateSrc)

	f, err := os.Open(webApp.TemplateSrc)
	if err != nil {
		log.Fatalf("error occured while trying to open file: got '%s'", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error occured while trying to read content of file: got '%s'", err)
	}

	t, err := template.New("index.html").Parse(string(data))
	if err != nil {
		log.Fatalf("error occured while parsing template [%s]. got %s\n", webApp.TemplateSrc, err)
	}
	tpl = t
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
