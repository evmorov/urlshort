package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

type PathToUrls struct {
	Path string
	Url  string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for shortUrl, redirectUrl := range pathsToUrls {
			if r.URL.Path == shortUrl {
				http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
			}
		}
		fallback.ServeHTTP(w, r)
	})
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]PathToUrls, error) {
	var out []PathToUrls
	err := yaml.UnmarshalStrict([]byte(yml), &out)
	return out, err
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	pathsMap := buildMap(parsedJSON)
	return MapHandler(pathsMap, fallback), nil
}

func parseJSON(jsonData []byte) ([]PathToUrls, error) {
	var out []PathToUrls
	err := json.Unmarshal(jsonData, &out)
	return out, err
}

func buildMap(pathToUrls []PathToUrls) map[string]string {
	pathsToUrlsMap := make(map[string]string)
	for _, pathToUrl := range pathToUrls {
		pathsToUrlsMap[pathToUrl.Path] = pathToUrl.Url
	}
	return pathsToUrlsMap
}
