package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type App struct {
	Port           int
	TemplateSource string
	renderer       renderer
}

type renderer struct {
	template     *template.Template
	templateData map[string]string
}

func StartApp(app *App) {
	app.renderer = renderer{
		templateData: make(map[string]string),
	}
	app.parseAndLoadTemplate()
	app.serve()
}

func (app *App) UpdateTemplateData(name, status string) {
	app.renderer.templateData[name] = status
}

func (app *App) serve() {
	log.Printf("Trying to start App on port [%d]...\n", app.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("error while trying to server mux server on port [%d]. got %s\n", app.Port, err)
	}
}

func (app *App) parseAndLoadTemplate() {
	log.Printf("Trying to parse HTML template [%s]...\n", app.TemplateSource)

	f, err := os.Open(app.TemplateSource)
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
		log.Fatalf("error occured while parsing template [%s]. got %s\n", app.TemplateSource, err)
	}
	app.renderer.template = t
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	err := app.renderer.template.Execute(w, app.renderer.templateData)
	if err != nil {
		log.Printf("error occured during executing template rendering. %s\n", err)
	}
}
