package main

import "net/http"
func healthCheck(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}