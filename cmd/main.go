package main

import (
	"flag"
	"fmt"
	"github.com/gophercises/urlshort"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
	"os"
)

var yamlFilename *string
var jsonFilename *string

func main() {
	parseFlags()
	handler := redisHandler(jsonHandler(yamlHandler(mapHandler(defaultMux()))))
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world")
}

func parseFlags() {
	yamlFilename = flag.String("yaml", "paths.yml", "a yaml file with paths to urls")
	jsonFilename = flag.String("json", "paths.json", "a json file with paths to urls")
	flag.Parse()
}

func mapHandler(handler http.Handler) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/1":              "http://localhost:8080/2",
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	return urlshort.MapHandler(pathsToUrls, handler)
}

func yamlHandler(handler http.Handler) http.HandlerFunc {
	yml, err := ioutil.ReadFile(*yamlFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	yamlHandler, err := urlshort.YAMLHandler(yml, handler)

	if err != nil {
		panic(err)
	}
	return yamlHandler
}

func jsonHandler(handler http.Handler) http.HandlerFunc {
	json, err := ioutil.ReadFile(*jsonFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonHandler, err := urlshort.JSONHandler(json, handler)
	if err != nil {
		panic(err)
	}
	return jsonHandler
}

func redisHandler(handler http.Handler) http.HandlerFunc {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})

	pathsToUrls := make(map[string]string)

	keys, err := client.Keys("*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		val, err := client.Get(key).Result()
		if err != nil {
			panic(err)
		}
		pathsToUrls[key] = val
	}

	return urlshort.MapHandler(pathsToUrls, handler)
}
