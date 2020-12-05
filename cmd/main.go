package main

import (
	"flag"
	"fmt"
	"github.com/gophercises/urlshort"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	mux := defaultMux()
	yamlFilename := flag.String("yaml", "paths.yml", "a yaml file with paths to urls")
	jsonFilename := flag.String("json", "paths.json", "a json file with paths to urls")
	flag.Parse()

	pathsToUrls := map[string]string{
		"/1":              "http://localhost:8080/2",
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yml, err := ioutil.ReadFile(*yamlFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	yamlHandler, err := urlshort.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}

	json, err := ioutil.ReadFile(*jsonFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world")
}
