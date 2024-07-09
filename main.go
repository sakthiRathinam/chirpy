package main

import (
	"fmt"
	"net/http"
)


const addr = ":8080"

func main(){
	serveMux := http.ServeMux{}
	apiConfig := apiConfig{}
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