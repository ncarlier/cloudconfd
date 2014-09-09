package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	laddr     = flag.String("l", ":8080", "HTTP service address (e.g.address, ':8080')")
	templates = template.New("cloudconfd")
)

type Config struct {
	Hostname            string
	Ip                  string
	Gateway             string
	Ssh_authorized_keys []string
}

func parseTemplateFiles() {
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filename := filepath.Base(path)
			log.Println(fmt.Sprintf("Loading template: %s ...", filename))
			filetext, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			text := string(filetext)
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			template.Must(templates.New(name).Parse(text))
			log.Println(fmt.Sprintf("Template %s loaded.", name))
		}
		return nil
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	mac := params["mac"]

	// open configuration file
	filename := fmt.Sprintf("./conf/%s/%s.yaml", name, mac)
	filename = strings.Replace(filename, ":", "_", -1)
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Load YAML configuration
	var config Config
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, name, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Load templates
	parseTemplateFiles()

	flag.Parse()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{name:[a-z-]+}/{mac:([0-9a-f]{2}([:-]|$)){6}}", handler).Methods("GET")

	http.Handle("/", rtr)

	log.Println("cloudconfd server listening...")
	log.Fatal(http.ListenAndServe(*laddr, nil))
}
