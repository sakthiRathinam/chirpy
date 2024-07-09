package main

import (
	"fmt"
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