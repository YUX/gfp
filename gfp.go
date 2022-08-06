package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var buffer = make([]byte, 128)

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func proxy(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.FormValue("targetUrl")
	if targetUrl == "" {
		http.Error(w, "targetUrl parameter is missing", http.StatusBadRequest)
		return // stop processing
	}
	if !isValidUrl(targetUrl) {
		http.Error(w, "targetUrl parameter is unsupported protocol scheme", http.StatusBadRequest)
		return // stop processing
	}
	log.Println("targetUrl", targetUrl)

	resp, err := http.Get(targetUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = io.CopyBuffer(w, resp.Body, buffer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/", proxy)
	log.Fatalln(http.ListenAndServe(":3333", nil))
}
