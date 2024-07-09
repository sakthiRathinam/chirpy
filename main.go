package main

import (
	"fmt"
	"net/http"
)


const addr = ":8080"

func main(){
	serveMux := http.ServeMux{}
	apiConfig := apiConfig{}
	registerRoutes(&serveMux,&apiConfig)
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


func registerRoutes(mux *http.ServeMux,apiConfig *apiConfig){
	mux.HandleFunc("/healthz",healthCheck)
	mux.HandleFunc("/metrics",apiConfig.getMetrics)
	mux.HandleFunc("/reset",apiConfig.resetMetrics)
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/*",http.StripPrefix("/app",apiConfig.middlewareIncrementHit(fileServer)))
}