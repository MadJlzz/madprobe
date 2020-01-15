package main

import (
	"flag"
	"github.com/go-yaml/yaml"
	"github.com/madjlzz/madprobe/internal"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var (
	config = flag.String("config", "test/sample.yml", "configuration file used to define probes")
)

func main() {
	flag.Parse()

	f, err := os.Open(*config)
	if err != nil {
		log.Fatalf("error occured while trying to open file: got '%s'", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("error occured while trying to read content of file: got '%s'", err)
	}

	var doc internal.Doc
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		log.Fatalf("error occured while trying to unmarshal yaml file: got '%s'", err)
	}

	var probers []internal.Prober
	registerPIDProbes(doc.PidProbes, &probers)
	registerHttpProbes(doc.HttpProbes, &probers)

	log.Printf("Running probes defined in [%s] yaml file.\n", *config)
	runProbers(probers)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Stopping probes...")
}

func runProbers(probers []internal.Prober) {
	for _, prober := range probers {
		go prober.Probe()
	}
}

func registerHttpProbes(httpProbes []internal.HttpProbe, probers *[]internal.Prober) {
	for _, probe := range httpProbes {
		*probers = append(*probers, &probe)
	}
}

func registerPIDProbes(pidProbes []internal.PidProbe, probers *[]internal.Prober) {
	for _, probe := range pidProbes {
		*probers = append(*probers, &probe)
	}
}
