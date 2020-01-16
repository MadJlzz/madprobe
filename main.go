package main

import (
	"flag"
	"github.com/go-yaml/yaml"
	"github.com/madjlzz/madprobe/internal/probe"
	"github.com/madjlzz/madprobe/internal/server"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var (
	yamlConfig     = flag.String("yamlConfig", "configs/sample.yml", "configuration file used to define probes")
	port           = flag.Int("port", 8081, "port webapp application will listen to")
	templateConfig = flag.String("templateConfig", "configs/index.gohtml", "configuration file used to render app")
)

func main() {
	flag.Parse()

	f, err := os.Open(*yamlConfig)
	if err != nil {
		log.Fatalf("error occured while trying to open file: got '%s'", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error occured while trying to read content of file: got '%s'", err)
	}

	var probes probe.Probes
	err = yaml.Unmarshal(data, &probes)
	if err != nil {
		log.Fatalf("error occured while trying to unmarshal yaml file: got '%s'", err)
	}

	app := &server.App{
		Port:           *port,
		TemplateSource: *templateConfig,
	}
	probe.Run(&probes, app)
	server.StartApp(app)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Stopping probes...")
}
