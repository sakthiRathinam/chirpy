package main

import (
	"fmt"
	"html/template"
	"net/http"
)
type apiConfig struct {
	fileServerHits int
}

func (ac *apiConfig) middlewareIncrementHit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ac.fileServerHits = ac.fileServerHits + 1
		next.ServeHTTP(w,r)
	})
}

func (ac *apiConfig) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d",ac.fileServerHits)))
	w.WriteHeader(200)
}

func (ac *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	ac.fileServerHits = 0
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d",ac.fileServerHits)))
	w.WriteHeader(200)
}

func (ac *apiConfig) getAdminMetrics(w http.ResponseWriter, r *http.Request) {
	teml, err := template.ParseFiles("admin.html")
	if err != nil {
		panic(err)
	}
	fmt.Println(ac.fileServerHits)
	data := map[string]int{"fileServerHits":ac.fileServerHits}
	err = teml.Execute(w, data)
	if err != nil {
		panic(err)
	}
}