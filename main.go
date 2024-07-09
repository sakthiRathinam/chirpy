package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)


const addr = ":8080"

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

func main(){
	serveMux := http.ServeMux{}
	apiConfig := apiConfig{}
	serveMux.HandleFunc("/home",homePage)
	serveMux.HandleFunc("/healthz",healthCheck)
	serveMux.HandleFunc("/metrics",apiConfig.getMetrics)
	serveMux.HandleFunc("/reset",apiConfig.resetMetrics)
	fileServer := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/*",http.StripPrefix("/app",apiConfig.middlewareIncrementHit(fileServer)))
	httpServer := http.Server{
		Handler:&serveMux,
		Addr: addr,
	}
	fmt.Printf("Server started on port %s\n",addr)
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("Failed while starting the server",err)
	}
}




func healthCheck(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","text/plain; charset=utf-8")
	w.Write([]byte("OK"))
	w.WriteHeader(200)
}

func homePage(w http.ResponseWriter,r *http.Request){
	htmFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println("error while parsing the request")
		http.Error(w,"Something went wrong",500)
		return 
	}
	defer htmFile.Close()
	w.Header().Set("Content-Type","text/html")

	_,err = io.Copy(w,htmFile)

	if err != nil {
		fmt.Println("error while copying the request")
		http.Error(w,"Something went wrong",500)
		return
	}
}